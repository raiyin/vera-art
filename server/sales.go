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

func mapSaleModelToSaleResponse(sale models.Sale) dtos.SaleResponse {
	var dir = config.AppConfigInstance.Directories.RelSalesDir +
		sale.StrId + "/"
	dto := dtos.SaleResponse{
		Id:          sale.Id,
		StrId:       sale.StrId,
		Dir:         dir,
		NameRu:      sale.NameRu,
		NameEn:      sale.NameEn,
		Year:        sale.Year,
		Descr:       sale.Descr,
		Width:       sale.Width,
		Height:      sale.Height,
		Price:       sale.Price,
		BaseRu:      "",
		BaseEn:      "",
		Images:      strings.Split(sale.Images, ";"),
		MaterialsEn: []string{},
		MaterialsRu: []string{},
	}
	// Query base_ru and base_en
	db.QueryRow("select base_ru, base_en from bases where id = ?", sale.BaseId).Scan(&dto.BaseRu, &dto.BaseEn)

	// Query materials_ru and materials_en arrays from many-to-many works_materials table and materials tables
	rows, err := db.Query("select material_id from sales_materials where sale_id = ?", sale.Id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	materialIds := []int{}
	for rows.Next() {
		var materialId int
		err := rows.Scan(&materialId)
		if err != nil {
			panic(err)
		}
		materialIds = append(materialIds, materialId)
	}
	for _, materialId := range materialIds {
		var materialRu, materialEn string
		db.QueryRow("select material_ru, material_en from materials where id = ?", materialId).Scan(&materialRu, &materialEn)
		dto.MaterialsRu = append(dto.MaterialsRu, materialRu)
		dto.MaterialsEn = append(dto.MaterialsEn, materialEn)
	}

	return dto
}

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

	query := "select * from sales"

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

	defer rows.Close()
	sales := []dtos.SaleResponse{}
	for rows.Next() {
		s := models.Sale{}
		err := rows.Scan(
			&s.Id, &s.Width,
			&s.Height, &s.Year, &s.Price,
			&s.NameRu, &s.NameEn, &s.BaseId, &s.StrId,
			&s.Descr, &s.Images)

		if err != nil {
			fmt.Println(err)
			continue
		}

		dto := mapSaleModelToSaleResponse(s)
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
	var sale dtos.CreateSaleRequest
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
		config.AppConfigInstance.Directories.RelSalesDir+strings.Replace(sale.NameEn, " ", "_", -1)+"/",
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
	dirPath := config.AppConfigInstance.Directories.AbsSalesDir + strings.Replace(sale.NameEn, " ", "_", -1) + "/"
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
	query := `SELECT * FROM sales WHERE id = ?`

	var sale models.Sale
	err := db.QueryRow(query, id).Scan(
		&sale.Id, &sale.Width, &sale.Height, &sale.Year,
		&sale.Price, &sale.NameRu, &sale.NameEn, &sale.BaseId, &sale.StrId,
		&sale.Descr, &sale.Images,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "sale not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	getSaleDto := mapSaleModelToSaleResponse(sale)

	c.JSON(http.StatusOK, getSaleDto)
}

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
	var sale dtos.UpdateSaleRequest
	if err := json.Unmarshal([]byte(dataJson), &sale); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 3. Either need to delete the old dir
	// Get sale data, if
	query := `
		SELECT name_en
		FROM sales
		WHERE id = ?
	`

	var sale_for_dir models.Sale
	err := db.QueryRow(query, id).Scan(&sale_for_dir.NameEn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
		config.AppConfigInstance.Directories.RelSalesDir+strings.Replace(sale.NameEn, " ", "_", -1)+"/",
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

	// Process directory and images
	dirPath := config.AppConfigInstance.Directories.AbsSalesDir + strings.Replace(sale.NameEn, " ", "_", -1) + "/"
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

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sale updated successfully"})
}
