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
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/config"
	"github.com/raiyin/artserver/dtos"
	"github.com/raiyin/artserver/models"
)

func getFileNamesString(files []string) string {
	names := make([]string, len(files))
	copy(names, files)
	return strings.Join(names, ";")
}

// workToDto converts a models.Work to dtos.GetWorkDto using material maps.
func workToDto(work models.Work, materialsEn map[int]string, materialsRu map[int]string) dtos.GetWorkDto {
	dto := dtos.GetWorkDto{
		Id:          work.Id,
		StrId:       work.StrId,
		Dir:         work.Dir,
		NameRu:      work.NameRu,
		NameEn:      work.NameEn,
		Year:        work.Year,
		Descr:       work.Descr,
		BaseRu:      work.BaseRu,
		BaseEn:      work.BaseEn,
		Width:       work.Width,
		Height:      work.Height,
		Type:        work.Type,
		Images:      work.Images,
		MaterialsEn: []string{},
		MaterialsRu: []string{},
	}
	// Populate material names
	for _, mid := range work.MaterialsIds {
		if en, ok := materialsEn[mid]; ok {
			dto.MaterialsEn = append(dto.MaterialsEn, en)
		}
		if ru, ok := materialsRu[mid]; ok {
			dto.MaterialsRu = append(dto.MaterialsRu, ru)
		}
	}
	return dto
}

func GetWorks(c *gin.Context) {

	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "select w.*, b.base_ru, b.base_en from works w join bases b on w.base_id = b.id"

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

	var images sql.NullString
	works := []models.Work{}
	for rows.Next() {
		w := models.Work{}
		err := rows.Scan(
			&w.Id, &w.Dir, &w.Width, &w.Height, &w.Year, &w.NameRu,
			&w.NameEn, &w.BaseId, &w.StrId, &w.Descr,
			&w.Type, &images, &w.BaseRu, &w.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if images.Valid {
			w.Images = strings.Split(images.String, ";")
		} else {
			w.Images = []string{}
		}

		works = append(works, w)
	}

	// If no works, return empty array
	if len(works) == 0 {
		c.JSON(http.StatusOK, []dtos.GetWorkDto{})
		return
	}

	// Collect work IDs
	workIds := make([]int, len(works))
	for i, w := range works {
		workIds[i] = w.Id
	}

	// Convert to []any for query
	args := make([]any, len(workIds))
	for i, id := range workIds {
		args[i] = id
	}

	// Fetch materials for these works
	queryStr := `
		SELECT wm.work_id, m.id, m.material_en, m.material_ru
		FROM works_materials wm
		JOIN materials m ON wm.material_id = m.id
		WHERE wm.work_id IN (` + strings.Repeat("?,", len(workIds)-1) + `?)`
	materialRows, err := db.Query(queryStr, args...)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer materialRows.Close()

	// Map work ID -> material IDs -> names
	materialsEn := make(map[int]string)
	materialsRu := make(map[int]string)
	workMaterials := make(map[int][]int)
	for materialRows.Next() {
		var workId, materialId int
		var materialEn, materialRu string
		err := materialRows.Scan(&workId, &materialId, &materialEn, &materialRu)
		if err != nil {
			fmt.Println(err)
			continue
		}
		materialsEn[materialId] = materialEn
		materialsRu[materialId] = materialRu
		workMaterials[workId] = append(workMaterials[workId], materialId)
	}

	// Assign material IDs to each work
	for i := range works {
		works[i].MaterialsIds = workMaterials[works[i].Id]
	}

	// Convert to DTOs
	dtos := make([]dtos.GetWorkDto, len(works))
	for i, w := range works {
		dtos[i] = workToDto(w, materialsEn, materialsRu)
	}

	c.JSON(http.StatusOK, dtos)
}

func AddWork(c *gin.Context) {

	// 1. Get the JSON metadata from form field
	dataJson := c.PostForm("data")
	if dataJson == "" {
		c.JSON(400, gin.H{"error": "data is required"})
		return
	}

	// fmt.Printf("json is: %q\n", dataJson)

	// 2. Parse the JSON
	var work models.Work
	if err := json.Unmarshal([]byte(dataJson), &work); err != nil {
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
	err = tx.QueryRow("select MAX(id) from works").Scan(&maxID)
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
		"insert into works (id, dir, width, height, year, name_ru, name_en, base_id, str_id, descr, type, images) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID+1,
		config.AppConfigInstance.Directories.WorksDbDirPrefix+strings.Replace(work.NameEn, " ", "_", -1)+"/",
		work.Width,
		work.Height,
		work.Year,
		work.NameRu,
		work.NameEn,
		work.BaseId,
		strings.Replace(work.NameEn, " ", "_", -1),
		work.Descr,
		work.Type,
		strings.Join(work.Images, ";"),
	)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error: %q\n", err)
		return
	}

	for _, material_id := range work.MaterialsIds {
		_, err = tx.Exec(
			"insert into works_materials (work_id, material_id) "+
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
	dirPath := config.AppConfigInstance.Directories.WorksDirSave +
		config.AppConfigInstance.Directories.WorksDbDirPrefix +
		strings.Replace(work.NameEn, " ", "_", -1) + "/"

	fmt.Printf("dirPath: %s\n", dirPath)
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
			dst, err := os.Create(dirPath + work.Images[index])
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

	c.JSON(http.StatusOK, gin.H{"message": "Painting saved successfully"})
}

func DeleteWork(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Get the work directory path before deleting the record
	var dir string
	query := "SELECT dir FROM works WHERE id = ?"
	err = tx.QueryRow(query, id).Scan(&dir)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Delete from works_materials table
	_, err = tx.Exec("DELETE FROM works_materials WHERE work_id IN (SELECT id FROM works WHERE id = ?)", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete from works table
	result, err := tx.Exec("DELETE FROM works WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the directory and its contents
	dirPath := config.AppConfigInstance.Directories.WorksDirSave + strings.Replace(strings.TrimSuffix(dir, "/"), config.AppConfigInstance.Directories.WorksDbDirPrefix, "", -1)
	if _, err := os.Stat(dirPath); err == nil {
		err := os.RemoveAll(dirPath)
		if err != nil {
			log.Printf("Error deleting directory %s: %v", dirPath, err)
			// We don't return an error here because the database record was successfully deleted
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work deleted successfully"})
}

func GetWorkById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Get work data
	query := `
		SELECT w.*, b.base_ru, b.base_en
		FROM works w
		JOIN bases b ON w.base_id = b.id
		WHERE w.id = ?
	`

	var images sql.NullString
	var work models.Work
	err := db.QueryRow(query, id).Scan(
		&work.Id, &work.Dir, &work.Width, &work.Height, &work.Year,
		&work.NameRu, &work.NameEn, &work.BaseId, &work.StrId,
		&work.Descr, &work.Type, &images,
		&work.BaseRu, &work.BaseEn)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if images.Valid {
		work.Images = strings.Split(images.String, ";")
	} else {
		work.Images = []string{}
	}

	// Get materials for this work
	materialsQuery := `
		SELECT m.id, m.material_ru, m.material_en
		FROM materials m
		JOIN works_materials wm ON m.id = wm.material_id
		WHERE wm.work_id = ?
	`

	materialRows, err := db.Query(materialsQuery, work.Id)
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

	work.MaterialsIds = materialsIds

	c.JSON(http.StatusOK, work)
}

// Update work by id
func UpdateWork(c *gin.Context) {
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
	var work models.Work
	if err := json.Unmarshal([]byte(dataJson), &work); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 3. Either need to delete the old dir
	// Get work data, if
	query := `
		SELECT dir, name_en
		FROM works
		WHERE id = ?
	`

	var work_for_dir models.Work
	err := db.QueryRow(query, id).Scan(&work_for_dir.Dir, &work_for_dir.NameEn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var needToDeleteOldDir = work_for_dir.Dir == work.Dir

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Update work record
	_, err = tx.Exec(`
		UPDATE works
		SET dir = ?, width = ?, height = ?, year = ?, name_ru = ?, name_en = ?,
			base_id = ?, str_id = ?, descr = ?, type = ?
		WHERE id = ?`,
		config.AppConfigInstance.Directories.WorksDbDirPrefix+strings.Replace(work.NameEn, " ", "_", -1)+"/",
		work.Width, work.Height, work.Year, work.NameRu, work.NameEn,
		work.BaseId, strings.Replace(work.NameEn, " ", "_", -1), work.Descr, work.Type, id)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete existing material associations
	_, err = tx.Exec("DELETE FROM works_materials WHERE work_id IN (SELECT id FROM works WHERE id = ?)", id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Insert new material associations
	for _, material_id := range work.MaterialsIds {
		// Get the work ID first
		var workId int
		err = tx.QueryRow("SELECT id FROM works WHERE id = ?", id).Scan(&workId)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = tx.Exec(
			"INSERT INTO works_materials (work_id, material_id) VALUES (?, ?)",
			workId, material_id)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Handle file uploads if any
	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	// Process directory and images
	dirPath := config.AppConfigInstance.Directories.WorksDirSave + strings.Replace(work.NameEn, " ", "_", -1) + "/"
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
	if len(work.RemovedIndices) > 0 {
		// Remove specified images from directory
		for _, index := range work.RemovedIndices {
			imagePath := filepath.Join(dirPath, strconv.Itoa(index)+".jpg")
			if _, err := os.Stat(imagePath); err == nil {
				err := os.Remove(imagePath)
				if err != nil {
					fmt.Printf("Error removing file %s: %v\n", imagePath, err)
				} else {
					fmt.Printf("Removed file: %s\n", imagePath)
				}
			}
		}

		// Rename remaining files to fill gaps
		// First, get all existing files and sort them
		files, err := os.ReadDir(dirPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read directory"})
			tx.Rollback()
			return
		}

		// Get list of remaining image numbers
		var remainingNumbers []int
		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".jpg" {
				if num, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".jpg")); err == nil {
					// Check if this number is not in removed indices
					isRemoved := false
					for _, removedIndex := range work.RemovedIndices {
						if num == removedIndex {
							isRemoved = true
							break
						}
					}
					if !isRemoved {
						remainingNumbers = append(remainingNumbers, num)
					}
				}
			}
		}

		// Sort remaining numbers
		sort.Ints(remainingNumbers)

		// Rename files to be sequential starting from 1
		for i, num := range remainingNumbers {
			oldPath := filepath.Join(dirPath, strconv.Itoa(num)+".jpg")
			newPath := filepath.Join(dirPath, strconv.Itoa(i+1)+".jpg")
			if oldPath != newPath {
				err := os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Printf("Error renaming file from %s to %s: %v\n", oldPath, newPath, err)
				} else {
					fmt.Printf("Renamed file from %s to %s\n", oldPath, newPath)
				}
			}
		}
	}

	// Handle file uploads if any
	if form != nil {
		// Get all files from the form
		files := form.File

		// Determine the next available index for new files
		nextIndex := 1
		if len(work.RemovedIndices) == 0 {
			// If no files were removed, next index is img_count + 1
			// But we need to check existing files
			if filesInDir, err := os.ReadDir(dirPath); err == nil {
				for _, file := range filesInDir {
					if !file.IsDir() && filepath.Ext(file.Name()) == ".jpg" {
						if num, err := strconv.Atoi(strings.TrimSuffix(file.Name(), ".jpg")); err == nil {
							if num >= nextIndex {
								nextIndex = num + 1
							}
						}
					}
				}
			}
		} else {
			// If files were removed, next index is the count of remaining files + 1
			if filesInDir, err := os.ReadDir(dirPath); err == nil {
				count := 0
				for _, file := range filesInDir {
					if !file.IsDir() && filepath.Ext(file.Name()) == ".jpg" {
						count++
					}
				}
				nextIndex = count + 1
			}
		}

		// Save new images
		for fieldName, fileHeaders := range files {
			for index, fileHeader := range fileHeaders {
				// Open the file
				file, err := fileHeader.Open()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					tx.Rollback()
					return
				}
				defer file.Close()

				// Create a destination file
				dst, err := os.Create(dirPath + strconv.Itoa(nextIndex+index) + filepath.Ext(fileHeader.Filename))
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					tx.Rollback()
					return
				}
				defer dst.Close()

				// Copy the file data
				if _, err := io.Copy(dst, file); err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					tx.Rollback()
					return
				}

				log.Printf("Saved file %s from field %s", fileHeader.Filename, fieldName)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	if needToDeleteOldDir {
		dirPath := config.AppConfigInstance.Directories.WorksDirSave + strings.Replace(work_for_dir.NameEn, " ", "_", -1) + "/"
		err = os.RemoveAll(dirPath)
		fmt.Printf("Delete dir is: %v\n", dirPath)
		if err != nil {
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work updated successfully"})
}
