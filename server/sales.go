package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/config"
	"github.com/raiyin/artserver/dtos"
	"github.com/raiyin/artserver/models"
)

func MaterialsBySaleId(saleId int) []models.Material {
	query := "select m.* from materials m join sales_materials sm on m.id = sm.material_id where sm.sale_id = ?"
	rows, err := db.Query(query, saleId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	materials := []models.Material{}
	for rows.Next() {
		m := models.Material{}
		err := rows.Scan(
			&m.Id, &m.MaterialRu, &m.MaterialEn)

		if err != nil {
			fmt.Println(err)
			continue
		}

		materials = append(materials, m)
	}

	return materials
}

func Map[T, U any](slice []T, fn func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func GetSales(c *gin.Context) {

	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "select s.*, b.base_ru, b.base_en from sales s join bases b on s.base_id = b.id"

	if len(limit) > 0 {
		query = query + " limit " + limit

		if len(offset) > 0 {
			query = query + " offset " + offset
		}
	}

	rows, err := db.Query(query)

	if err != nil {
		panic(err)
	}

	var images sql.NullString
	var base_en sql.NullString
	var base_ru sql.NullString

	defer rows.Close()
	sales := []dtos.GetSaleDto{}
	for rows.Next() {
		p := models.Sale{}
		err := rows.Scan(
			&p.Id, &p.Dir, &p.Width,
			&p.Height, &p.Year, &p.Price,
			&p.NameRu, &p.NameEn, &p.BaseId, &p.StrId,
			&p.Descr, &images, &base_ru, &base_en)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if images.Valid {
			p.Images = strings.Split(images.String, ";")
		} else {
			p.Images = []string{}
		}

		dto := dtos.GetSaleDto{
			Id:          p.Id,
			Dir:         p.Dir,
			Width:       p.Width,
			Height:      p.Height,
			Year:        p.Year,
			Price:       p.Price,
			NameRu:      p.NameRu,
			NameEn:      p.NameEn,
			BaseRu:      base_ru.String,
			BaseEn:      base_en.String,
			StrId:       p.StrId,
			Descr:       p.Descr,
			Images:      p.Images,
			MaterialsEn: []string{},
			MaterialsRu: []string{},
		}

		var materials []models.Material = MaterialsBySaleId(p.Id)
		dto.MaterialsEn = Map(materials, func(m models.Material) string {
			return m.MaterialEn
		})
		dto.MaterialsRu = Map(materials, func(m models.Material) string {
			return m.MaterialRu
		})

		sales = append(sales, dto)
	}

	c.JSON(http.StatusOK, sales)
}

func CreateSale(c *gin.Context) {

	// 1. Get the JSON metadata from form field
	dataJson := c.PostForm("data")
	if dataJson == "" {
		c.JSON(400, gin.H{"error": "data is required"})
		return
	}

	fmt.Printf("json is: %q\n", dataJson)

	// 2. Parse the JSON
	var sale models.Sale
	if err := json.Unmarshal([]byte(dataJson), &sale); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error: %q\n", err)
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var maxID int
	err = tx.QueryRow("select MAX(id) from sales").Scan(&maxID)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error: %q\n", err)
		return
	}

	_, err = tx.Exec(
		"insert into sales (id, dir, width, height, year, price, name_ru, name_en, base_id, str_id, descr) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID+1,
		config.AppConfigInstance.Directories.SaleDbDirPrefix+strings.Replace(sale.NameEn, " ", "_", -1)+"/",
		sale.Width,
		sale.Height,
		sale.Year,
		sale.Price,
		sale.NameRu,
		sale.NameEn,
		sale.BaseId,
		strings.Replace(sale.NameEn, " ", "_", -1),
		sale.Descr)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error: %q\n", err)
		return
	}

	for _, material_id := range sale.MaterialsIds {
		_, err = tx.Exec(
			"insert into sales_materials (sale_id, material_id) "+
				"values (?, ?)",
			maxID+1,
			material_id)

		if err != nil {
			tx.Rollback()
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Error : %q\n", err)
			return
		}
	}

	// 3. Get the uploaded files and save them
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error : %q\n", err)
		return
	}

	// Get all files from the form
	files := form.File

	// Process each file
	dirPath := config.AppConfigInstance.Directories.SaleDirSave + strings.Replace(sale.NameEn, " ", "_", -1) + "/"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory: %v\n"})
			tx.Rollback()
			log.Fatal(err)
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Error: %s\n", err)
			fmt.Printf("Error : %q\n", err)
			return
		}
		fmt.Printf("Directory '%s' created successfully.\n", dirPath)
	}

	for fieldName, fileHeaders := range files {
		for index, fileHeader := range fileHeaders {
			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				log.Fatal(err)
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Error: %s\n", err)
				fmt.Printf("Error : %q\n", err)
				return
			}
			defer file.Close()

			// Create a destination file
			dst, err := os.Create(dirPath + strconv.Itoa(index+1) + filepath.Ext(fileHeader.Filename))
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				log.Fatal(err)
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Error: %s\n", err)
				fmt.Printf("Error : %q\n", err)
				return
			}
			defer dst.Close()

			// Copy the file data
			if _, err := io.Copy(dst, file); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				log.Fatal(err)
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Error: %s\n", err)
				fmt.Printf("Error : %q\n", err)
				return
			}

			log.Printf("Saved file %s from field %s", fileHeader.Filename, fieldName)
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error : %q\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sale saved successfully"})
}

func GetSaleById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Get sale data
	query := `
		SELECT s.*, b.base_ru, b.base_en
		FROM sales s
		JOIN bases b ON s.base_id = b.id
		WHERE s.id = ?
	`

	var images sql.NullString
	var sale dtos.GetSaleDto
	err := db.QueryRow(query, id).Scan(
		&sale.Id, &sale.Dir, &sale.Width, &sale.Height, &sale.Year,
		&sale.Price, &sale.NameRu, &sale.NameEn, &sale.StrId,
		&sale.Descr, &images, &sale.BaseRu, &sale.BaseEn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "sale not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if images.Valid {
		sale.Images = strings.Split(images.String, ";")
	} else {
		sale.Images = []string{}
	}

	// Get materials for this sale
	materialsQuery := `
		SELECT m.id, m.material_ru, m.material_en
		FROM materials m
		JOIN sales_materials sm ON m.id = sm.material_id
		WHERE sm.sale_id = ?
	`

	materialRows, err := db.Query(materialsQuery, sale.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer materialRows.Close()

	var materialsIds []int
	for materialRows.Next() {
		var material models.Material
		err := materialRows.Scan(&material.Id, &material.MaterialRu, &material.MaterialEn)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		materialsIds = append(materialsIds, material.Id)
	}

	// Create a Sale struct with materials_ids
	saleWithMaterials := struct {
		dtos.GetSaleDto
		MaterialsIds []int `json:"materials_ids"`
	}{
		GetSaleDto:   sale,
		MaterialsIds: materialsIds,
	}

	c.JSON(http.StatusOK, saleWithMaterials)
}

// Update sale by id
func UpdateSale(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 1. Get the JSON metadata from form field
	dataJson := c.PostForm("data")
	if dataJson == "" {
		c.JSON(400, gin.H{"error": "data is required"})
		return
	}

	// 2. Parse the JSON
	var sale models.Sale
	if err := json.Unmarshal([]byte(dataJson), &sale); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 3. Either need to delete the old dir
	// Get sale data, if
	query := `
		SELECT dir, name_en
		FROM sales
		WHERE id = ?
	`

	var sale_for_dir models.Sale
	err := db.QueryRow(query, id).Scan(&sale_for_dir.Dir, &sale_for_dir.NameEn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var needToDeleteOldDir = sale_for_dir.Dir == sale.Dir

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Update sale record
	_, err = tx.Exec(`
		UPDATE sales
		SET dir = ?, width = ?, height = ?, year = ?, price = ?, name_ru = ?, name_en = ?,
			base_id = ?, str_id = ?, descr = ?
		WHERE id = ?`,
		config.AppConfigInstance.Directories.SaleDbDirPrefix+strings.Replace(sale.NameEn, " ", "_", -1)+"/",
		sale.Width, sale.Height, sale.Year, sale.Price, sale.NameRu, sale.NameEn,
		sale.BaseId, strings.Replace(sale.NameEn, " ", "_", -1), sale.Descr, id)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete existing material associations
	_, err = tx.Exec("DELETE FROM sales_materials WHERE sale_id IN (SELECT id FROM sales WHERE id = ?)", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert new material associations
	for _, material_id := range sale.MaterialsIds {
		// Get the sale ID first
		var saleId int
		err = tx.QueryRow("SELECT id FROM sales WHERE id = ?", id).Scan(&saleId)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = tx.Exec(
			"INSERT INTO sales_materials (sale_id, material_id) VALUES (?, ?)",
			saleId, material_id)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Handle file uploads if any
	// form, err := c.MultipartForm()
	// if err != nil && err != http.ErrNotMultipart {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	tx.Rollback()
	// 	return
	// }

	// Process directory and images
	dirPath := config.AppConfigInstance.Directories.SaleDirSave + strings.Replace(sale.NameEn, " ", "_", -1) + "/"
	fmt.Printf("dirPath: %s\n", dirPath)

	// Create directory if it doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory: %v\n"})
			tx.Rollback()
			return
		}
		fmt.Printf("Directory '%s' created successfully.\n", dirPath)
	}

	// Handle removed images if any
	// if len(sale.RemovedIndices) > 0 {
	// 	// Remove specified images from directory
	// 	for _, index := range sale.RemovedIndices {
	// 		imagePath := filepath.Join(dirPath, strconv.Itoa(index)+".jpg")
	// 		if _, err := os.Stat(imagePath); err == nil {
	// 			err := os.Remove(imagePath)
	// 			if err != nil {
	// 				fmt.Printf("Error removing file %s: %v\n", imagePath, err)
	// 			} else {
	// 				fmt.Printf("Removed file: %s\n", imagePath)
	// 			}
	// 		}
	// 	}

	// 	// Rename remaining files to fill gaps
	// 	// First, get all existing files and sort them
	// 	files, err := os.ReadDir(dirPath)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read directory"})
	// 		tx.Rollback()
	// 		return
	// 	}

	// 	// Get list of remaining image numbers
	// 	var remainingNumbers []int
	// 	for _, file := range files {
	// 		if !file.IsDir() && filepath.Ext(file.Name()) == ".jpg" {
	// 			if num, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".jpg")); err == nil {
	// 				// Check if this number is not in removed indices
	// 				isRemoved := false
	// 				for _, removedIndex := range sale.RemovedIndices {
	// 					if num == removedIndex {
	// 						isRemoved = true
	// 						break
	// 					}
	// 				}
	// 				if !isRemoved {
	// 					remainingNumbers = append(remainingNumbers, num)
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// Sort remaining numbers
	// 	sort.Ints(remainingNumbers)

	// 	// Rename files to be sequential starting from 1
	// 	for i, num := range remainingNumbers {
	// 		oldPath := filepath.Join(dirPath, strconv.Itoa(num)+".jpg")
	// 		newPath := filepath.Join(dirPath, strconv.Itoa(i+1)+".jpg")
	// 		if oldPath != newPath {
	// 			err := os.Rename(oldPath, newPath)
	// 			if err != nil {
	// 				fmt.Printf("Error renaming file from %s to %s: %v\n", oldPath, newPath, err)
	// 			} else {
	// 				fmt.Printf("Renamed file from %s to %s\n", oldPath, newPath)
	// 			}
	// 		}
	// 	}
	// }

	// Handle file uploads if any
	// if form != nil {
	// 	// Get all files from the form
	// 	files := form.File

	// 	// Determine the next available index for new files
	// 	nextIndex := 1
	// 	if len(sale.RemovedIndices) == 0 {
	// 		// If no files were removed, next index is img_count + 1
	// 		// But we need to check existing files
	// 		if filesInDir, err := os.ReadDir(dirPath); err == nil {
	// 			for _, file := range filesInDir {
	// 				if !file.IsDir() && filepath.Ext(file.Name()) == ".jpg" {
	// 					if num, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".jpg")); err == nil {
	// 						if num >= nextIndex {
	// 							nextIndex = num + 1
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	} else {
	// 		// If files were removed, next index is the count of remaining files + 1
	// 		if filesInDir, err := os.ReadDir(dirPath); err == nil {
	// 			count := 0
	// 			for _, file := range filesInDir {
	// 				if !file.IsDir() && filepath.Ext(file.Name()) == ".jpg" {
	// 					count++
	// 				}
	// 			}
	// 			nextIndex = count + 1
	// 		}
	// 	}

	// 	// Save new images
	// 	for fieldName, fileHeaders := range files {
	// 		for index, fileHeader := range fileHeaders {
	// 			// Open the file
	// 			file, err := fileHeader.Open()
	// 			if err != nil {
	// 				c.JSON(500, gin.H{"error": err.Error()})
	// 				tx.Rollback()
	// 				return
	// 			}
	// 			defer file.Close()

	// 			// Create a destination file
	// 			dst, err := os.Create(dirPath + strconv.Itoa(nextIndex+index) + filepath.Ext(fileHeader.Filename))
	// 			if err != nil {
	// 				c.JSON(500, gin.H{"error": err.Error()})
	// 				tx.Rollback()
	// 				return
	// 			}
	// 			defer dst.Close()

	// 			// Copy the file data
	// 			if _, err := io.Copy(dst, file); err != nil {
	// 				c.JSON(500, gin.H{"error": err.Error()})
	// 				tx.Rollback()
	// 				return
	// 			}

	// 			log.Printf("Saved file %s from field %s", fileHeader.Filename, fieldName)
	// 		}
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	if needToDeleteOldDir {
		dirPath := config.AppConfigInstance.Directories.SaleDirSave + strings.Replace(sale_for_dir.NameEn, " ", "_", -1) + "/"
		err = os.RemoveAll(dirPath)
		fmt.Printf("Delete dir is: %v\n", dirPath)
		if err != nil {
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sale updated successfully"})
}
