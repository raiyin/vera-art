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
func workToDto(work models.Work) dtos.GetWorkDto {
	var dir = config.AppConfigInstance.Directories.RelWorksDir +
		work.StrId + "/"
	dto := dtos.GetWorkDto{
		Id:          work.Id,
		StrId:       work.StrId,
		Dir:         dir,
		NameRu:      work.NameRu,
		NameEn:      work.NameEn,
		Year:        work.Year,
		Descr:       work.Descr,
		Width:       work.Width,
		Height:      work.Height,
		Type:        work.Type,
		BaseRu:      "",
		BaseEn:      "",
		Images:      strings.Split(work.Images, ";"),
		MaterialsEn: []string{},
		MaterialsRu: []string{},
	}
	// Query base_ru and base_en
	db.QueryRow("select base_ru, base_en from bases where id = ?", work.BaseId).Scan(&dto.BaseRu, &dto.BaseEn)

	// Query materials_ru and materials_en arrays from many-to-many works_materials table and materials tables
	rows, err := db.Query("select material_id from works_materials where work_id = ?", work.Id)
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

func GetWorks(c *gin.Context) {

	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "select *  from works"

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

	works := []dtos.GetWorkDto{}
	for rows.Next() {
		w := models.Work{}
		err := rows.Scan(
			&w.Id, &w.Width, &w.Height, &w.Year, &w.NameRu,
			&w.NameEn, &w.BaseId, &w.StrId, &w.Descr,
			&w.Type, &w.Images)
		if err != nil {
			fmt.Println(err)
			continue
		}

		works = append(works, workToDto(w))
	}

	// If no works, return empty array
	if len(works) == 0 {
		c.JSON(http.StatusOK, []dtos.GetWorkDto{})
		return
	}

	c.JSON(http.StatusOK, works)
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
	var work dtos.AddWorkDto
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
		"insert into works (id, width, height, year, name_ru, name_en, base_id, str_id, descr, type, images) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID+1,
		work.Width,
		work.Height,
		work.Year,
		work.NameRu,
		work.NameEn,
		work.BaseId,
		strings.Replace(work.NameEn, " ", "_", -1),
		work.Descr,
		work.Type,
		work.Images,
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
	dirPath := config.AppConfigInstance.Directories.AbsWorksDir +
		config.AppConfigInstance.Directories.RelWorksDir +
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
		for _, fileHeader := range fileHeaders {
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
			dst, err := os.Create(dirPath + fileHeader.Filename)
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
	dirPath := config.AppConfigInstance.Directories.AbsWorksDir + strings.Replace(strings.TrimSuffix(dir, "/"), config.AppConfigInstance.Directories.RelWorksDir, "", -1)
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
	var work dtos.EditWorkDto
	err := db.QueryRow(query, id).Scan(
		&work.Id, &work.Width, &work.Height, &work.Year,
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
	var work dtos.EditWorkDto
	if err := json.Unmarshal([]byte(dataJson), &work); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 3. Get old work data to compare and clean up
	query := `
		SELECT name_en, str_id
		FROM works
		WHERE id = ?
	`

	var oldNameEn string
	var oldStrId string
	err := db.QueryRow(query, id).Scan(&oldNameEn, &oldStrId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Update work record
	// Build images string from work.Images slice
	imagesStr := work.Images
	// newDir := config.AppConfigInstance.Directories.WorksDbDirPrefix + strings.Replace(work.NameEn, " ", "_", -1) + "/"
	_, err = tx.Exec(`
		UPDATE works
		SET width = ?, height = ?, year = ?, name_ru = ?, name_en = ?,
			base_id = ?, str_id = ?, descr = ?, type = ?, images = ?
		WHERE id = ?`,
		work.Width, work.Height, work.Year, work.NameRu, work.NameEn,
		work.BaseId, strings.Replace(work.NameEn, " ", "_", -1), work.Descr, work.Type, imagesStr, id)

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
	// First, determine if we need to move to a new directory (name changed)
	newNameNormalized := strings.Replace(work.NameEn, " ", "_", -1)
	nameChanged := oldStrId != newNameNormalized

	oldDirPath := config.AppConfigInstance.Directories.AbsWorksDir +
		config.AppConfigInstance.Directories.RelWorksDir +
		oldStrId + "/"

	if nameChanged {

		// New directory path
		newDirPath := config.AppConfigInstance.Directories.AbsWorksDir +
			config.AppConfigInstance.Directories.RelWorksDir +
			newNameNormalized + "/"
		fmt.Printf("New dirPath: %s\n", newDirPath)

		// Create new directory if it doesn't exist
		if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
			err := os.MkdirAll(newDirPath, 0777)
			if err != nil {
				fmt.Printf("Error creating directory: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory: %v\n"})
				tx.Rollback()
				return
			}
			fmt.Printf("Directory '%s' created successfully.\n", newDirPath)
		}

		// Clean up old images in the new directory (if directory already existed)
		// Delete all files in the new directory to start fresh
		if _, err := os.Stat(newDirPath); err == nil {
			files, err := filepath.Glob(filepath.Join(newDirPath, "*"))
			if err == nil {
				for _, file := range files {
					os.Remove(file)
				}
			}
		}

		// Handle file uploads if any
		if form != nil {
			// Get all files from the form
			files := form.File

			// Save new images with their original filenames (sanitized)
			for fieldName, fileHeaders := range files {
				for _, fileHeader := range fileHeaders {
					// Open the file
					file, err := fileHeader.Open()
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						tx.Rollback()
						return
					}
					defer file.Close()

					// Use original filename but sanitize it
					originalFilename := fileHeader.Filename
					// Sanitize filename: replace spaces with underscores, remove special characters
					// sanitizedFilename := strings.Replace(originalFilename, " ", "_", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, "(", "", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, ")", "", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, "'", "", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, "\"", "", -1)

					// Create a destination file
					// dst, err := os.Create(newDirPath + sanitizedFilename)
					dst, err := os.Create(newDirPath + originalFilename)
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

					// log.Printf("Saved file %s from field %s as %s", fileHeader.Filename, fieldName, sanitizedFilename)
					log.Printf("Saved file %s from field %s as %s", fileHeader.Filename, fieldName, originalFilename)
				}
			}
		}

		// Clean up old directory if name changed
		if nameChanged {
			if _, err := os.Stat(oldDirPath); err == nil {
				err = os.RemoveAll(oldDirPath)
				fmt.Printf("Deleted old directory: %v\n", oldDirPath)
				if err != nil {
					log.Printf("Error deleting old directory %s: %v", oldDirPath, err)
				}
			}
		}

		if form != nil {
			// Get all files from the form
			files := form.File

			// Save new images with their original filenames (sanitized)
			for fieldName, fileHeaders := range files {
				for _, fileHeader := range fileHeaders {
					// Open the file
					file, err := fileHeader.Open()
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						tx.Rollback()
						return
					}
					defer file.Close()

					// Use original filename but sanitize it
					originalFilename := fileHeader.Filename
					// Sanitize filename: replace spaces with underscores, remove special characters
					// sanitizedFilename := strings.Replace(originalFilename, " ", "_", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, "(", "", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, ")", "", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, "'", "", -1)
					// sanitizedFilename = strings.Replace(sanitizedFilename, "\"", "", -1)

					// Create a destination file
					// dst, err := os.Create(newDirPath + sanitizedFilename)
					dst, err := os.Create(newDirPath + originalFilename)
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

					// log.Printf("Saved file %s from field %s as %s", fileHeader.Filename, fieldName, sanitizedFilename)
					log.Printf("Saved file %s from field %s as %s", fileHeader.Filename, fieldName, originalFilename)
				}
			}
		}

	} else {
		// remove all images from the old directory
		files, err := filepath.Glob(filepath.Join(oldDirPath, "*"))
		if err == nil {
			for _, file := range files {
				os.Remove(file)
			}
		}

	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work updated successfully"})
}
