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
		&news.ImgBackfull, &news.ImagesCount, &news.VideosCount,
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

	query := "select * from news"

	if len(id) > 0 {
		query = query + " where id = " + id
	} else if len(id_ne) > 0 && len(limit) > 0 {
		query = query +
			" where id != " +
			id_ne +
			" order by datetime desc" +
			" limit " + limit
	} else if len(limit) > 0 {
		query = query + " order by datetime desc" + " limit " + limit

		if len(offset) > 0 {
			query = query + " offset " + offset
		}
	}
	// query = query + " order by datetime desc"
	fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var images sql.NullString
	var videos sql.NullString
	news := []models.News{}
	for rows.Next() {
		p := models.News{}
		err := rows.Scan(
			&p.Id, &p.Datetime, &p.TitleRu, &p.TitleEn, &p.SubtitleRu,
			&p.SubtitleEn, &p.Dir, &p.ImgBack, &p.ImgBackfull, &p.ImagesCount,
			&p.VideosCount, &p.TextRu, &p.TextEn, &images, &videos)

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

	_, err = tx.Exec(
		"insert into news (id, datetime, title_ru, title_en, subtitle_ru, subtitle_en, dir, img_back, img_backfull, imagescount, videoscount, text_ru, text_en, images, videos) "+
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
		news.ImagesCount,
		news.VideosCount,
		news.TextRu,
		news.TextEn,
		news.Images,
		news.Videos)

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
			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				return
			}
			defer file.Close()

			// Determine file name based on field type
			if fieldName == "images" {
				// TODO make a bunch request
				_, err = tx.Exec(
					"insert into images (table_type_link, filename, entity_id) "+
						"values (?, ?, ?)",
					3,
					fileHeader.Filename,
					maxID)

				if err != nil {
					tx.Rollback()
					log.Fatal(err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else if fieldName == "videos" {

				// TODO make a bunch request
				_, err = tx.Exec(
					"insert into videos (table_type_link, filename, entity_id) "+
						"values (?, ?, ?)",
					3,
					fileHeader.Filename,
					maxID)

				if err != nil {
					tx.Rollback()
					log.Fatal(err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
			} else {
				fmt.Println("Unknown field type:", fileHeader.Filename)
			}

			// Create a destination file
			dst, err := os.Create(filepath.Join(dirPath, fileHeader.Filename))
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

	// Update news record
	_, err = tx.Exec(`
		UPDATE news
		SET datetime = ?, title_ru = ?, title_en = ?, subtitle_ru = ?, subtitle_en = ?,
		    dir = ?, img_back = ?, img_backfull = ?, imagescount = ?, videoscount = ?,
		    text_ru = ?, text_en = ?
		WHERE id = ?`,
		news.Datetime, news.TitleRu, news.TitleEn, news.SubtitleRu, news.SubtitleEn,
		news.Dir, "back.jpg", "backfull.jpg", news.ImagesCount, news.VideosCount,
		news.TextRu, news.TextEn, id)

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
				if fieldName == "img_back" {
					fileName = "back.jpg"
				} else if fieldName == "img_backfull" {
					fileName = "back_full.jpg"
				} else if fieldName == "images" {
					fileName = fmt.Sprintf("%d.jpg", index+1)
				} else if fieldName == "videos" {
					fileName = fmt.Sprintf("%d%s", index+1, filepath.Ext(fileHeader.Filename))
				} else {
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

func checkIdIdExists(id string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM news WHERE id = ?", id).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}
