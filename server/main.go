// info for db
// table_link: 1-illustrations, 2-sales, 3-news
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/raiyin/artserver/config"
	"github.com/raiyin/artserver/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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

func main() {

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("dir is ", dir)

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
	log.Printf("dir is %s and %s", appConfig.Directories.RelWorksDir, appConfig.Directories.AbsWorksDir)

	r_gin := gin.Default()
	// r.Run(fmt.Sprintf(":%d", appConfig.App.Port))

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true // Allow sending cookies/auth headers
	config.AllowAllOrigins = true
	r_gin.Use(cors.New(config))

	// Auth routes
	r_gin.POST("/register", Register)
	r_gin.POST("/login", Login)

	r_gin.GET("/sales", GetSales)
	r_gin.GET("/works", GetWorks)
	r_gin.GET("/works/:id", GetWorkById)
	r_gin.GET("/news", GetNews)
	r_gin.GET("/news/:id", GetNewsById)
	r_gin.GET("/materials", getMaterials)
	r_gin.GET("/bases", getBases)
	r_gin.GET("/sales/:id", GetSaleById)
	r_gin.POST("/works", AddWork)
	r_gin.PUT("/works/:id", AuthMiddleware(), UpdateWork)
	r_gin.POST("/sales", CreateSale)
	r_gin.PUT("/sales/:id", AuthMiddleware(), UpdateSale)
	r_gin.POST("/news", AddNews)
	r_gin.DELETE("/works/:id", AuthMiddleware(), DeleteWork)
	r_gin.DELETE("/news/:id", AuthMiddleware(), DeleteNews)
	r_gin.PUT("/news/:id", AuthMiddleware(), UpdateNews)

	defer db.Close()
	if err := r_gin.Run("localhost:8000"); err != nil {
		log.Fatal(err)
	}
}
