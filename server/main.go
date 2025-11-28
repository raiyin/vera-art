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
	"time"

	"github.com/raiyin/artserver/config"
	"github.com/raiyin/artserver/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var jwtKey = []byte("mydevsecret")
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

func saveSales(c *gin.Context) {

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
		config.AppConfigInstance.Directories.SaleDirPrefix+strings.Replace(sale.NameEn, " ", "_", -1)+"/",
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

			err = tx.Commit()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				log.Fatal(err)
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Error: %s\n", err)
				fmt.Printf("Error : %q\n", err)
			}
			log.Printf("Saved file %s from field %s", fileHeader.Filename, fieldName)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sale saved successfully"})
}

func getPaintings(c *gin.Context) {

	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "select p.*, b.base_ru, b.base_en from paintings p join bases b on p.base_id = b.id"

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
	paintings := []models.PaintingWithBase{}
	for rows.Next() {
		p := models.PaintingWithBase{}
		err := rows.Scan(&p.Id, &p.Dir, &p.Width, &p.Height, &p.Year, &p.NameRu, &p.NameEn, &p.BaseId, &p.StrId, &p.ImgCount, &p.Descr, &p.BaseRu, &p.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		paintings = append(paintings, p)
	}

	c.JSON(http.StatusOK, paintings)
}

func savePainting(c *gin.Context) {

	// 1. Get the JSON metadata from form field
	dataJson := c.PostForm("data")
	if dataJson == "" {
		c.JSON(400, gin.H{"error": "data is required"})
		return
	}

	fmt.Printf("json is: %q\n", dataJson)

	// 2. Parse the JSON
	var painting models.Painting
	if err := json.Unmarshal([]byte(dataJson), &painting); err != nil {
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
	err = tx.QueryRow("select MAX(id) from paintings").Scan(&maxID)
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
		"insert into paintings (id, dir, width, height, year, name_ru, name_en, base_id, str_id, img_count, descr) "+
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		maxID+1,
		config.AppConfigInstance.Directories.PaintingsDirPrefix+strings.Replace(painting.NameEn, " ", "_", -1)+"/",
		painting.Width,
		painting.Height,
		painting.Year,
		painting.NameRu,
		painting.NameEn,
		painting.BaseId,
		strings.Replace(painting.NameEn, " ", "_", -1),
		painting.ImgCount,
		painting.Descr)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Error: %q\n", err)
		return
	}

	for _, material_id := range painting.MaterialsIds {
		_, err = tx.Exec(
			"insert into painting_materials (painting_id, material_id) "+
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
	dirPath := config.AppConfigInstance.Directories.PaintingsDirSave + strings.Replace(painting.NameEn, " ", "_", -1) + "/"
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

			err = tx.Commit()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				tx.Rollback()
				log.Fatal(err)
				fmt.Printf("Error: %v\n", err)
				fmt.Printf("Error: %s\n", err)
				fmt.Printf("Error : %q\n", err)
			}
			log.Printf("Saved file %s from field %s", fileHeader.Filename, fieldName)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Painting saved successfully"})
}

func getThreeds(c *gin.Context) {
	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "select t.*, b.base_ru, b.base_en from threeds t join bases b on t.base_id = b.id"

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
	threeds := []models.Threed{}
	for rows.Next() {
		p := models.Threed{}
		err := rows.Scan(&p.Id, &p.StrId, &p.Dir, &p.NameRu, &p.NameEn, &p.BaseId, &p.Year, &p.ImgCount, &p.BaseRu, &p.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		threeds = append(threeds, p)
	}

	c.JSON(http.StatusOK, threeds)
}

func getIllustrations(c *gin.Context) {
	offset := c.Query("offset")
	limit := c.Query("limit")

	query := "select i.*, b.base_ru, b.base_en from illustrations i join bases b on i.base_id = b.id"

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
	illustrations := []models.Illustration{}
	for rows.Next() {
		p := models.Illustration{}
		err := rows.Scan(&p.Id, &p.StrId, &p.Dir, &p.NameRu, &p.NameEn, &p.BaseId, &p.Year, &p.ImgCount, &p.BaseRu, &p.BaseEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		illustrations = append(illustrations, p)
	}

	c.JSON(http.StatusOK, illustrations)
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

func saveNews(c *gin.Context) {
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
	// err = tx.QueryRow("select MAX(id) from news").Scan(&maxID)
	// if err != nil {
	// 	tx.Rollback()
	// 	log.Fatal(err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

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
	news.Dir = config.AppConfigInstance.Directories.NewsDirPrefix + dirPostfix

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
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if _, exists := users[user.Username]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	users[user.Username] = string(hashedPassword)
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	storedPassword, exists := users[user.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
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

func protected(c *gin.Context) {
	claims, _ := c.Get("claims")
	claimsObj, ok := claims.(*models.Claims)
	if !ok || claimsObj == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "protected content",
		"user":    claims.(*models.Claims).Username,
	})
}

func main() {

	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Access configuration values
	appConfig := config.AppConfigInstance
	log.Printf("Starting %s on port %d", appConfig.App.Name, appConfig.App.Port)

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
	r_gin.GET("/protected", authMiddleware(), protected)

	r_gin.GET("/sales", getSales)
	r_gin.GET("/paintings", getPaintings)
	r_gin.GET("/threeds", getThreeds)
	r_gin.GET("/illustrations", getIllustrations)
	r_gin.GET("/news", getNews)
	r_gin.GET("/materials", getMaterials)
	r_gin.GET("/bases", getBases)
	r_gin.POST("/paintings", savePainting)
	r_gin.POST("/sales", saveSales)
	r_gin.POST("/news", saveNews)

	defer db.Close()
	if err := r_gin.Run(":8000"); err != nil {
		log.Fatal(err)
	}

}
