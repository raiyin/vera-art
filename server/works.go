package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/dtos"
	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/models"
)

// workModelToWorkResponse converts a models.Work to dtos.GetWorkDto using material maps.
func workModelToWorkResponse(work models.Work) dtos.WorkResponse {
	var dir = config.AppConfigInstance.Directories.RelWorksDir +
		work.StrId + "/"
	dto := dtos.WorkResponse{
		Id:          work.Id,
		StrId:       work.StrId,
		Dir:         dir,
		NameRu:      work.NameRu,
		NameEn:      work.NameEn,
		Year:        work.Year,
		DescrRu:     work.DescrRu,
		DescrEn:     work.DescrEn,
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

func workModelToUpdateWorkResponse(work models.Work) dtos.UpdateWorkResponse {
	var dir = config.AppConfigInstance.Directories.RelWorksDir +
		work.StrId + "/"
	dto := dtos.UpdateWorkResponse{
		Id:           work.Id,
		StrId:        work.StrId,
		Dir:          dir,
		NameRu:       work.NameRu,
		NameEn:       work.NameEn,
		BaseId:       work.BaseId,
		Year:         work.Year,
		DescrRu:      work.DescrRu,
		DescrEn:      work.DescrEn,
		Width:        work.Width,
		Height:       work.Height,
		Type:         work.Type,
		Images:       strings.Split(work.Images, ";"),
		MaterialsIds: []int{},
	}

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

	dto.MaterialsIds = materialIds
	return dto
}

func GetWorks(c *gin.Context) {
	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "SELECT * FROM works"
	var args []interface{}

	if limit != "" {
		// Validate limit is a positive integer
		if limitInt, err := strconv.Atoi(limit); err == nil && limitInt > 0 {
			query += " LIMIT ?"
			args = append(args, limitInt)

			if offset != "" {
				// Validate offset is a non-negative integer
				if offsetInt, err := strconv.Atoi(offset); err == nil && offsetInt >= 0 {
					query += " OFFSET ?"
					args = append(args, offsetInt)
				}
			}
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	works := []dtos.WorkResponse{}
	for rows.Next() {
		w := models.Work{}
		err := rows.Scan(
			&w.Id, &w.Width, &w.Height, &w.Year, &w.NameRu,
			&w.NameEn, &w.BaseId, &w.StrId,
			&w.Type, &w.Images, &w.DescrRu, &w.DescrEn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		works = append(works, workModelToWorkResponse(w))
	}

	// If no works, return empty array
	if len(works) == 0 {
		c.JSON(http.StatusOK, []dtos.WorkResponse{})
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
	var work dtos.CreateWorkRequest
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
		log.Printf("Error getting max ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Convert images array to semicolon-separated string
	imagesStr := strings.Join(work.Images, ";")

	_, err = tx.Exec(
		"insert into works (id, width, height, year, name_ru, name_en, base_id, str_id, type, images, descr_ru, descr_en) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID+1,
		work.Width,
		work.Height,
		work.Year,
		work.NameRu,
		work.NameEn,
		work.BaseId,
		strings.Replace(work.NameEn, " ", "_", -1),
		work.Type,
		imagesStr,
		work.DescrRu,
		work.DescrEn,
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting work: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save work to database"})
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
			log.Printf("Error inserting work material: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save work materials"})
			return
		}
	}

	// 3. Get the uploaded files and save them
	form, err := c.MultipartForm()
	if err != nil {
		tx.Rollback()
		log.Printf("Error getting multipart form: %v", err)
		c.JSON(400, gin.H{"error": "Invalid form data"})
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
				tx.Rollback()
				log.Printf("Error creating destination file: %v", err)
				c.JSON(500, gin.H{"error": "Failed to save file"})
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
				tx.Rollback()
				log.Printf("Error copying file: %v", err)
				c.JSON(500, gin.H{"error": "Failed to save file"})
				return
			}

			log.Printf("Saved file %s from field %s", fileHeader.Filename, fieldName)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(500, gin.H{"error": "Failed to save work"})
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
	query := ` SELECT * FROM works WHERE id = ?`

	var work models.Work
	err := db.QueryRow(query, id).Scan(
		&work.Id, &work.Width, &work.Height, &work.Year,
		&work.NameRu, &work.NameEn, &work.BaseId, &work.StrId,
		&work.Type, &work.Images, &work.DescrRu, &work.DescrEn)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	workResponse := workModelToWorkResponse(work)

	c.JSON(http.StatusOK, workResponse)
}

func GetWorkByIdForEdit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// Get work data
	query := `SELECT * FROM works WHERE id = ?`

	var work models.Work
	err := db.QueryRow(query, id).Scan(
		&work.Id, &work.Width, &work.Height, &work.Year,
		&work.NameRu, &work.NameEn, &work.BaseId, &work.StrId,
		&work.Type, &work.Images,
		&work.DescrRu, &work.DescrEn)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var workResponse dtos.UpdateWorkResponse
	workResponse = workModelToUpdateWorkResponse(work)
	c.JSON(http.StatusOK, workResponse)

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
	var work dtos.UpdateWorkRequest
	if err := json.Unmarshal([]byte(dataJson), &work); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 3. Get old work data to compare and clean up
	query := `SELECT name_en, str_id FROM works WHERE id = ?`

	var oldNameEn string
	var oldStrId string
	if err := db.QueryRow(query, id).Scan(&oldNameEn, &oldStrId); err != nil {
		log.Printf("Error getting work data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Work not found"})
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Update work record
	// Build images string from work.Images slice
	imagesStr := strings.Join(work.Images, ";")
	_, err = tx.Exec(`
		UPDATE works
		SET width = ?, height = ?, year = ?, name_ru = ?, name_en = ?,
			base_id = ?, str_id = ?, descr_ru = ?, descr_en = ?, type = ?, images = ?
		WHERE id = ?`,
		work.Width, work.Height, work.Year, work.NameRu, work.NameEn,
		work.BaseId, strings.Replace(work.NameEn, " ", "_", -1),
		work.DescrRu, work.DescrEn, work.Type, imagesStr, id)

	if err != nil {
		tx.Rollback()
		log.Printf("Error updating work: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work"})
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
		_, err = tx.Exec(
			"INSERT INTO works_materials (work_id, material_id) VALUES (?, ?)",
			id, material_id)

		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting work material: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work materials"})
			return
		}
	}

	// rename old dir name by new
	newStrId := strings.Replace(work.NameEn, " ", "_", -1)
	if oldStrId != newStrId {
		err = os.Rename(
			config.AppConfigInstance.Directories.AbsWorksDir+
				config.AppConfigInstance.Directories.RelWorksDir+oldStrId+"/",
			config.AppConfigInstance.Directories.AbsWorksDir+
				config.AppConfigInstance.Directories.RelWorksDir+newStrId+"/")

		if err != nil && !os.IsNotExist(err) {
			tx.Rollback()
			log.Printf("Error renaming directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename work directory"})
			return
		}
	}

	// Delete only files in dir with names that do not exist in work.Images
	files, err := os.ReadDir(
		config.AppConfigInstance.Directories.AbsWorksDir +
			config.AppConfigInstance.Directories.RelWorksDir + strings.Replace(work.NameEn, " ", "_", -1) + "/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, file := range files {
		if !slices.Contains(work.Images, file.Name()) {
			err = os.Remove(
				config.AppConfigInstance.Directories.AbsWorksDir +
					config.AppConfigInstance.Directories.RelWorksDir +
					strings.Replace(work.NameEn, " ", "_", -1) + "/" +
					file.Name())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Handle file uploads if any
	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	if form != nil {
		for _, f := range form.File {
			for _, file := range f {
				// Upload the file to the server
				err = c.SaveUploadedFile(file,
					config.AppConfigInstance.Directories.AbsWorksDir+
						config.AppConfigInstance.Directories.RelWorksDir+
						strings.Replace(work.NameEn, " ", "_", -1)+"/"+
						file.Filename)
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					tx.Rollback()
					return
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update work"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Work updated successfully"})
}
