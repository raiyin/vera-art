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
	"time"

	"github.com/joho/godotenv"
	"github.com/raiyin/artserver/config"
	"github.com/raiyin/artserver/dtos"
	"github.com/raiyin/artserver/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

// var jwtKey = "mydevsecret"
var users = make(map[string]string)

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		panic(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
}

func getSales(c *gin.Context) {

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

	defer rows.Close()
	sales := []models.SaleWithBase{}
	for rows.Next() {
		p := models.SaleWithBase{}
		err := rows.Scan(
			&p.Id, &p.Dir, &p.Width,
			&p.Height, &p.Year, &p.Price,
			&p.NameRu, &p.NameEn,
			&p.BaseId, &p.StrId,
			&p.ImgCount, &p.Descr,
			&p.BaseRu, &p.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sales = append(sales, p)
	}

	c.JSON(http.StatusOK, sales)
}

func createSale(c *gin.Context) {

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
		"insert into sales (id, dir, width, height, year, price, name_ru, name_en, base_id, str_id, img_count, descr) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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
		sale.ImgCount,
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
			"insert into sale_materials (sale_id, material_id) "+
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

func getWorks(c *gin.Context) {

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
	works := []models.Work{}
	for rows.Next() {
		w := models.Work{}
		err := rows.Scan(&w.Id, &w.Dir, &w.Width, &w.Height, &w.Year, &w.NameRu, &w.NameEn, &w.BaseId, &w.StrId, &w.ImgCount, &w.Descr, &w.Type, &w.BaseRu, &w.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		works = append(works, w)
	}

	c.JSON(http.StatusOK, works)
}

func addWork(c *gin.Context) {

	// 1. Get the JSON metadata from form field
	dataJson := c.PostForm("data")
	if dataJson == "" {
		c.JSON(400, gin.H{"error": "data is required"})
		return
	}

	fmt.Printf("json is: %q\n", dataJson)

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
		"insert into works (id, dir, width, height, year, name_ru, name_en, base_id, str_id, img_count, descr, type) "+
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
		work.ImgCount,
		work.Descr,
		work.Type)

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
	dirPath := config.AppConfigInstance.Directories.WorksDirSave + strings.Replace(work.NameEn, " ", "_", -1) + "/"
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

	c.JSON(http.StatusOK, gin.H{"message": "Painting saved successfully"})
}

func deleteWork(c *gin.Context) {
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

func getNews(c *gin.Context) {

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
	news := []models.News{}
	for rows.Next() {
		p := models.News{}
		err := rows.Scan(&p.Id, &p.Datetime, &p.TitleRu, &p.TitleEn, &p.SubtitleRu, &p.SubtitleEn, &p.Dir, &p.ImgBack, &p.ImgBackfull, &p.ImagesCount, &p.VideosCount, &p.TextRu, &p.TextEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		news = append(news, p)
	}

	c.JSON(http.StatusOK, news)
}

func addNews(c *gin.Context) {
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

	var maxID string = strings.ReplaceAll(news.Datetime, "-", "")
	fmt.Println(news.Datetime)

	var dirPostfix string = strings.Split(news.Datetime, "-")[0] +
		"/" +
		strings.Split(news.Datetime, "-")[1] +
		"/" +
		strings.Split(news.Datetime, "-")[2] +
		"/"

	// Create directory for news item
	var dirPath = config.AppConfigInstance.Directories.NewsDirSave + dirPostfix
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

	// Update news.Dir with the actual directory path
	news.Dir = config.AppConfigInstance.Directories.NewsDbDirPrefix + dirPostfix

	// Update image file names
	news.ImgBack = "back.jpg"
	news.ImgBackfull = "back_full.jpg"

	_, err = tx.Exec(
		"insert into news (id, datetime, title_ru, title_en, subtitle_ru, subtitle_en, dir, img_back, img_backfull, imagescount, videoscount, text_ru, text_en) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID,
		news.Datetime,
		news.TitleRu,
		news.TitleEn,
		news.SubtitleRu,
		news.SubtitleEn,
		news.Dir,
		news.ImgBack,
		news.ImgBackfull,
		news.ImagesCount,
		news.VideosCount,
		news.TextRu,
		news.TextEn)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

	err = tx.Commit()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		tx.Rollback()
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "News saved successfully"})
}

func deleteNews(c *gin.Context) {
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
	dirPath := config.AppConfigInstance.Directories.NewsDirSave + strings.Replace(strings.TrimSuffix(dir, "/"), config.AppConfigInstance.Directories.NewsDbDirPrefix, "", -1)
	if _, err := os.Stat(dirPath); err == nil {
		err := os.RemoveAll(dirPath)
		if err != nil {
			log.Printf("Error deleting directory %s: %v", dirPath, err)
			// We don't return an error here because the database record was successfully deleted
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

func getMaterials(c *gin.Context) {

	query := "select * from materials"

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	materials := []models.Material{}
	for rows.Next() {
		p := models.Material{}
		err := rows.Scan(&p.Id, &p.MaterialRu, &p.MaterialEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		materials = append(materials, p)
	}

	c.JSON(http.StatusOK, materials)
}

func getBases(c *gin.Context) {

	query := "select * from bases"

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	bases := []models.Base{}
	for rows.Next() {
		p := models.Base{}
		err := rows.Scan(&p.Id, &p.BaseRu, &p.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		bases = append(bases, p)
	}

	c.JSON(http.StatusOK, bases)
}

func register(c *gin.Context) {
	var registerUserRequest dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&registerUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var count int64
	query := "select count(*) from users where username = ?"
	err := db.QueryRow(query, registerUserRequest.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check if user exists"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUserRequest.Password), bcrypt.DefaultCost)
	var stringHashedPassword = string(hashedPassword)

	// fmt.Println("%X ", hashedPassword)
	// for _, num := range hashedPassword {
	// 	fmt.Printf("%X ", num)
	// }
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	query = "insert into users (username, pass_hash) " +
		"values (?, ?)"
	_, err = db.Exec(query, registerUserRequest.Username, stringHashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func login(c *gin.Context) {
	var loginUserRequest dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	var username string
	var passhash string

	query := "select * from users where username = ?"
	rows := db.QueryRow(query, loginUserRequest.Username)
	err := rows.Scan(&id, &username, &passhash)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "undefinded servererror"})
			return
		}
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(loginUserRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	claims := &models.Claims{
		Username: loginUserRequest.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	fmt.Printf("jwt key is %v", jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"expires": expirationTime,
	})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func getWorkById(c *gin.Context) {
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

	var work models.Work
	err := db.QueryRow(query, id).Scan(
		&work.Id, &work.Dir, &work.Width, &work.Height, &work.Year,
		&work.NameRu, &work.NameEn, &work.BaseId, &work.StrId,
		&work.ImgCount, &work.Descr, &work.Type, &work.BaseRu, &work.BaseEn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "work not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
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

func getSaleById(c *gin.Context) {
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

	var sale models.SaleWithBase
	err := db.QueryRow(query, id).Scan(
		&sale.Id, &sale.Dir, &sale.Width, &sale.Height, &sale.Year,
		&sale.Price, &sale.NameRu, &sale.NameEn, &sale.BaseId, &sale.StrId,
		&sale.ImgCount, &sale.Descr, &sale.BaseRu, &sale.BaseEn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "sale not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Get materials for this sale
	materialsQuery := `
		SELECT m.id, m.material_ru, m.material_en
		FROM materials m
		JOIN sale_materials sm ON m.id = sm.material_id
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
		models.SaleWithBase
		MaterialsIds []int `json:"materials_ids"`
	}{
		SaleWithBase: sale,
		MaterialsIds: materialsIds,
	}

	c.JSON(http.StatusOK, saleWithMaterials)
}

func getNewsById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	query := "SELECT * FROM news WHERE id = ?"
	var news models.News
	err := db.QueryRow(query, id).Scan(
		&news.Id, &news.Datetime, &news.TitleRu, &news.TitleEn,
		&news.SubtitleRu, &news.SubtitleEn, &news.Dir, &news.ImgBack,
		&news.ImgBackfull, &news.ImagesCount, &news.VideosCount,
		&news.TextRu, &news.TextEn)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, news)
}

// Update work by id
func updateWork(c *gin.Context) {
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
			base_id = ?, str_id = ?, img_count = ?, descr = ?, type = ?
		WHERE id = ?`,
		config.AppConfigInstance.Directories.WorksDbDirPrefix+strings.Replace(work.NameEn, " ", "_", -1)+"/",
		work.Width, work.Height, work.Year, work.NameRu, work.NameEn,
		work.BaseId, strings.Replace(work.NameEn, " ", "_", -1), work.ImgCount, work.Descr, work.Type, id)

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

// Update sale by id
func updateSale(c *gin.Context) {
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
			base_id = ?, str_id = ?, img_count = ?, descr = ?
		WHERE id = ?`,
		config.AppConfigInstance.Directories.SaleDbDirPrefix+strings.Replace(sale.NameEn, " ", "_", -1)+"/",
		sale.Width, sale.Height, sale.Year, sale.Price, sale.NameRu, sale.NameEn,
		sale.BaseId, strings.Replace(sale.NameEn, " ", "_", -1), sale.ImgCount, sale.Descr, id)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete existing material associations
	_, err = tx.Exec("DELETE FROM sale_materials WHERE sale_id IN (SELECT id FROM sales WHERE id = ?)", id)
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
			"INSERT INTO sale_materials (sale_id, material_id) VALUES (?, ?)",
			saleId, material_id)

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
	if len(sale.RemovedIndices) > 0 {
		// Remove specified images from directory
		for _, index := range sale.RemovedIndices {
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
					for _, removedIndex := range sale.RemovedIndices {
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
		if len(sale.RemovedIndices) == 0 {
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
		dirPath := config.AppConfigInstance.Directories.SaleDirSave + strings.Replace(sale_for_dir.NameEn, " ", "_", -1) + "/"
		err = os.RemoveAll(dirPath)
		fmt.Printf("Delete dir is: %v\n", dirPath)
		if err != nil {
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sale updated successfully"})
}

// Update news by id
func updateNews(c *gin.Context) {
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
		// news.Datetime, news.TitleRu, news.TitleEn, news.SubtitleRu, news.SubtitleEn,
		// news.Dir, news.ImgBack, news.ImgBackfull, news.ImagesCount, news.VideosCount,
		// news.TextRu, news.TextEn, id)
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

// func protected(c *gin.Context) {
// 	claims, _ := c.Get("claims")
// 	claimsObj, ok := claims.(*models.Claims)
// 	if !ok || claimsObj == nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "protected content",
// 		"user":    claims.(*models.Claims).Username,
// 	})
// }

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assuming environment variables are set externally")
	}

	// jwtKey = os.Getenv("JWT_SECRET")

	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Access configuration values
	appConfig := config.AppConfigInstance
	log.Printf("Starting %s on port %d", appConfig.App.Name, appConfig.App.Port)
	log.Printf("dir is %s and %s", appConfig.Directories.WorksDbDirPrefix, appConfig.Directories.WorksDirSave)

	r_gin := gin.Default()
	// r.Run(fmt.Sprintf(":%d", appConfig.App.Port))

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true // Allow sending cookies/auth headers
	config.AllowAllOrigins = true
	r_gin.Use(cors.New(config))

	// Auth routes
	r_gin.POST("/register", register)
	r_gin.POST("/login", login)
	// r_gin.GET("/protected", authMiddleware(), protected)

	r_gin.GET("/sales", getSales)
	r_gin.GET("/works", getWorks)
	r_gin.GET("/works/:id", getWorkById)
	r_gin.GET("/news", getNews)
	r_gin.GET("/news/:id", getNewsById)
	r_gin.GET("/materials", getMaterials)
	r_gin.GET("/bases", getBases)
	r_gin.GET("/sales/:id", getSaleById)
	r_gin.POST("/works", addWork)
	r_gin.PUT("/works/:id", authMiddleware(), updateWork)
	r_gin.POST("/sales", createSale)
	r_gin.PUT("/sales/:id", authMiddleware(), updateSale)
	r_gin.POST("/news", addNews)
	r_gin.DELETE("/works/:id", authMiddleware(), deleteWork)
	r_gin.DELETE("/news/:id", authMiddleware(), deleteNews)
	r_gin.PUT("/news/:id", authMiddleware(), updateNews)

	defer db.Close()
	if err := r_gin.Run("localhost:8000"); err != nil {
		log.Fatal(err)
	}

}
