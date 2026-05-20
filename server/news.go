package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/models"
)

func GetNewsById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var images sql.NullString
	var videos sql.NullString
	query := "SELECT * FROM news WHERE id = ?"
	var news models.News
	err := db.QueryRow(query, id).Scan(
		&news.Id, &news.Datetime, &news.TitleRu, &news.TitleEn,
		&news.SubtitleRu, &news.SubtitleEn, &news.Dir, &news.ImgBack,
		&news.ImgBackfull,
		&news.TextRu, &news.TextEn, &images, &videos)
	// &news.TextRu, &news.TextEn, &news.Images, &news.Videos)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if images.Valid {
		news.Images = strings.Split(images.String, ";")
	} else {
		news.Images = []string{}
	}

	if videos.Valid {
		news.Videos = strings.Split(videos.String, ";")
	} else {
		news.Videos = []string{}
	}

	c.JSON(http.StatusOK, news)
}

func GetNews(c *gin.Context) {
	id := c.Query("id")
	offset := c.Query("offset")
	limit := c.Query("limit")
	id_ne := c.Query("id_ne")

	// Build query with parameterized placeholders
	query := "SELECT * FROM news"
	var args []interface{}
	var whereClauses []string

	if id != "" {
		whereClauses = append(whereClauses, "id = ?")
		args = append(args, id)
	}

	if id_ne != "" {
		whereClauses = append(whereClauses, "id != ?")
		args = append(args, id_ne)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY datetime DESC"

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

	fmt.Println("Executing query:", query, "with args:", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var images sql.NullString
	var videos sql.NullString
	news := []models.News{}
	for rows.Next() {
		p := models.News{}
		err := rows.Scan(
			&p.Id, &p.Datetime, &p.TitleRu, &p.TitleEn, &p.SubtitleRu,
			&p.SubtitleEn, &p.Dir, &p.ImgBack, &p.ImgBackfull, &p.TextRu, &p.TextEn, &images, &videos)

		if images.Valid {
			p.Images = strings.Split(images.String, ";")
		} else {
			p.Images = []string{}
		}

		if videos.Valid {
			p.Videos = strings.Split(videos.String, ";")
		} else {
			p.Videos = []string{}
		}

		if err != nil {
			fmt.Println(err)
			continue
		}
		news = append(news, p)
	}

	c.JSON(http.StatusOK, news)
}

func AddNews(c *gin.Context) {
	// 1. Get the JSON metadata from form field
	dataJson := c.PostForm("data")
	if dataJson == "" {
		c.JSON(400, gin.H{"error": "data is required"})
		return
	}

	// 2. Parse the JSON
	var news models.News
	if err := json.Unmarshal([]byte(dataJson), &news); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	var counter = 0
	var maxID string = strings.ReplaceAll(news.Datetime, "-", "")
	var isIdExists = checkIdIdExists(maxID)
	for isIdExists {
		counter++
		maxID = strings.ReplaceAll(news.Datetime, "-", "") + strconv.Itoa(counter)
		isIdExists = checkIdIdExists(maxID)
	}
	fmt.Println(news.Datetime)

	var dirPostfix string = filepath.Join(
		strings.Split(news.Datetime, "-")[0],
		strings.Split(news.Datetime, "-")[1],
		strings.Split(news.Datetime, "-")[2])

	// Update news.Dir with the actual directory path
	//config.AppConfigInstance.Directories.NewsDirSave
	var temp_db_dir = filepath.Join(config.AppConfigInstance.Directories.NewsDbDirPrefix, dirPostfix)
	if !strings.HasSuffix(temp_db_dir, string(os.PathSeparator)) {
		temp_db_dir += string(os.PathSeparator)
	}
	dir_int_postfix := 0

	info, err := os.Stat(temp_db_dir)
	for err == nil && info.IsDir() {
		dir_int_postfix++
		temp_db_dir = filepath.Join(config.AppConfigInstance.Directories.NewsDbDirPrefix, dirPostfix+"_"+strconv.Itoa(dir_int_postfix))
		if !strings.HasSuffix(temp_db_dir, string(os.PathSeparator)) {
			temp_db_dir += string(os.PathSeparator)
		}
		info, err = os.Stat(temp_db_dir)
	}
	// news.Dir = filepath.Join(config.AppConfigInstance.Directories.NewsDirSave, temp_db_dir)
	if !strings.HasSuffix(news.Dir, string(os.PathSeparator)) {
		news.Dir += string(os.PathSeparator)
	}

	// Convert images and videos slices to semicolon-separated strings
	imagesStr := strings.Join(news.Images, ";")
	videosStr := strings.Join(news.Videos, ";")

	_, err = tx.Exec(
		"insert into news (id, datetime, title_ru, title_en, subtitle_ru, subtitle_en, dir, img_back, img_backfull, text_ru, text_en, images, videos) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID,
		news.Datetime,
		news.TitleRu,
		news.TitleEn,
		news.SubtitleRu,
		news.SubtitleEn,
		temp_db_dir,
		news.ImgBack,
		news.ImgBackfull,
		news.TextRu,
		news.TextEn,
		imagesStr,
		videosStr)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create directory for news item
	var dirPath = filepath.Join(config.AppConfigInstance.Directories.NewsDirSave, temp_db_dir)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0777)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating directory"})
			tx.Rollback()
			return
		}
		fmt.Printf("Directory '%s' created successfully.\n", dirPath)
	}

	// 3. Get the uploaded files and save them
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	// Get all files from the form
	files := form.File

	// Process each file
	for fieldName, fileHeaders := range files {
		// Skip the "data" field as it's not a file
		if fieldName == "data" {
			continue
		}

		for _, fileHeader := range fileHeaders {
			// Validate the uploaded file
			if err := validateUploadedFile(fileHeader); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				tx.Rollback()
				return
			}

			// Sanitize filename
			safeFilename := sanitizeFilename(fileHeader.Filename)

			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				return
			}
			defer file.Close()

			// Determine file name based on field type
			switch fieldName {
			case "images":
				// TODO make a bunch request
				_, err = tx.Exec(
					"insert into images (table_type_link, filename, entity_id) "+
						"values (?, ?, ?)",
					3,
					safeFilename,
					maxID)

				if err != nil {
					tx.Rollback()
					log.Fatal(err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			case "videos":

				// TODO make a bunch request
				_, err = tx.Exec(
					"insert into videos (table_type_link, filename, entity_id) "+
						"values (?, ?, ?)",
					3,
					safeFilename,
					maxID)

				if err != nil {
					tx.Rollback()
					log.Fatal(err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			default:
				fmt.Println("Unknown field type:", safeFilename)
			}

			// Create a destination file
			dst, err := os.Create(filepath.Join(dirPath, safeFilename))
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

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "News saved successfully"})
}

func DeleteNews(c *gin.Context) {
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

	// Get the news directory path before deleting the record
	var dir string
	query := "SELECT dir FROM news WHERE id = ?"
	err = tx.QueryRow(query, id).Scan(&dir)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Delete from news table
	result, err := tx.Exec("DELETE FROM news WHERE id = ?", id)
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
		c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the directory and its contents
	dirPath := filepath.Join(
		config.AppConfigInstance.Directories.NewsDirSave +
			strings.Replace(strings.TrimSuffix(dir, "/"), config.AppConfigInstance.Directories.NewsDbDirPrefix, "", -1))
	if _, err := os.Stat(dirPath); err == nil {
		err := os.RemoveAll(dirPath)
		if err != nil {
			log.Printf("Error deleting directory %s: %v", dirPath, err)
			// We don't return an error here because the database record was successfully deleted
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

// Update news by id
func UpdateNews(c *gin.Context) {
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
	var news models.News
	if err := json.Unmarshal([]byte(dataJson), &news); err != nil {
		c.JSON(400, gin.H{"error": "invalid data format"})
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Convert images and videos slices to semicolon-separated strings
	imagesStr := strings.Join(news.Images, ";")
	videosStr := strings.Join(news.Videos, ";")

	// Update news record
	_, err = tx.Exec(`
		UPDATE news
		SET datetime = ?, title_ru = ?, title_en = ?, subtitle_ru = ?, subtitle_en = ?,
		    dir = ?, img_back = ?, img_backfull = ?,
		    text_ru = ?, text_en = ?, images = ?, videos = ?
		WHERE id = ?`,
		news.Datetime, news.TitleRu, news.TitleEn, news.SubtitleRu, news.SubtitleEn,
		news.Dir, news.ImgBack, news.ImgBackfull,
		news.TextRu, news.TextEn, imagesStr, videosStr, id)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process directory and images
	var dirPostfix string = strings.Split(news.Datetime, "-")[0] +
		"/" +
		strings.Split(news.Datetime, "-")[1] +
		"/" +
		strings.Split(news.Datetime, "-")[2] +
		"/"

	var dirPath = config.AppConfigInstance.Directories.NewsDirSave + dirPostfix
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

	// Handle file uploads if any
	form, err := c.MultipartForm()
	if err != nil && err != http.ErrNotMultipart {
		c.JSON(400, gin.H{"error": err.Error()})
		tx.Rollback()
		return
	}

	// Handle file uploads if any
	if form != nil {
		// Get all files from the form
		files := form.File

		// Process each file
		for fieldName, fileHeaders := range files {
			// Skip the "data" field as it's not a file
			if fieldName == "data" {
				continue
			}

			for index, fileHeader := range fileHeaders {
				// Open the file
				file, err := fileHeader.Open()
				if err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					tx.Rollback()
					return
				}
				defer file.Close()

				// Determine file name based on field type
				var fileName string
				switch fieldName {
				case "img_back":
					fileName = "back.jpg"
				case "img_backfull":
					fileName = "back_full.jpg"
				case "images":
					fileName = fmt.Sprintf("%d.jpg", index+1)
				case "videos":
					fileName = fmt.Sprintf("%d%s", index+1, filepath.Ext(fileHeader.Filename))
				default:
					fileName = fmt.Sprintf("%s_%d%s", fieldName, index+1, filepath.Ext(fileHeader.Filename))
				}

				// Create a destination file
				dst, err := os.Create(dirPath + fileName)
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

				log.Printf("Saved file %s from field %s", fileName, fieldName)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "News updated successfully"})
}

// validateUploadedFile checks if an uploaded file is safe to save
func validateUploadedFile(fileHeader *multipart.FileHeader) error {
	// Check file size (max 10MB)
	const maxFileSize = 10 << 20 // 10 MB
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("file %s is too large: %d bytes (max: %d bytes)",
			fileHeader.Filename, fileHeader.Size, maxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".pdf":  true,
	}
	if !allowedExtensions[ext] {
		return fmt.Errorf("file %s has disallowed extension: %s", fileHeader.Filename, ext)
	}

	// Check for path traversal in filename
	if strings.Contains(fileHeader.Filename, "..") || strings.Contains(fileHeader.Filename, "/") {
		return fmt.Errorf("invalid filename: %s", fileHeader.Filename)
	}

	return nil
}

// sanitizeFilename removes dangerous characters from filename
func sanitizeFilename(filename string) string {
	// Remove path components
	filename = filepath.Base(filename)
	// Replace spaces and special characters
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "..", "")
	// Keep only alphanumeric, dots, underscores, and hyphens
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = reg.ReplaceAllString(filename, "")
	return filename
}

func checkIdIdExists(id string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM news WHERE id = ?", id).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
