package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book
var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
		panic(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(25)
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	type ResponseId struct {
		Id int `json:"id"`
	}

	js, err := json.Marshal(ResponseId{Id: rand.Int()})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getSales(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query()["start"]
	limit := r.URL.Query()["limit"]

	query := "select * from sales"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(start) > 0 {
			query = query + " offset " + start[0]
		}
	}

	rows, err := db.Query(query)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	sales := []Sale{}
	for rows.Next() {
		p := Sale{}
		err := rows.Scan(&p.Id, &p.Dir, &p.Width, &p.Height, &p.Year, &p.Price, &p.NameRu, &p.NameEn, &p.BaseId, &p.StrId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sales = append(sales, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

func getPaintings(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query()["start"]
	limit := r.URL.Query()["limit"]

	query := "select * from paintings"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(start) > 0 {
			query = query + " offset " + start[0]
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	paintings := []Painting{}
	for rows.Next() {
		p := Painting{}
		err := rows.Scan(&p.Id, &p.Dir, &p.Width, &p.Height, &p.Year, &p.NameRu, &p.NameEn, &p.BaseId, &p.StrId)
		if err != nil {
			fmt.Println(err)
			continue
		}
		paintings = append(paintings, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paintings)
}

func getThreeds(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query()["start"]
	limit := r.URL.Query()["limit"]

	query := "select * from threeds"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(start) > 0 {
			query = query + " offset " + start[0]
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	threeds := []Threed{}
	for rows.Next() {
		p := Threed{}
		err := rows.Scan(&p.Id, &p.StrId, &p.Dir, &p.NameRu, &p.NameEn, &p.BaseId, &p.Year)
		if err != nil {
			fmt.Println(err)
			continue
		}
		threeds = append(threeds, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threeds)
}

func getIllustrations(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query()["start"]
	limit := r.URL.Query()["limit"]

	query := "select * from illustrations"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(start) > 0 {
			query = query + " offset " + start[0]
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	illustrations := []Illustration{}
	for rows.Next() {
		p := Illustration{}
		err := rows.Scan(&p.Id, &p.StrId, &p.Dir, &p.NameRu, &p.NameEn, &p.BaseId, &p.Year)
		if err != nil {
			fmt.Println(err)
			continue
		}
		illustrations = append(illustrations, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(illustrations)
}

func getNews(w http.ResponseWriter, r *http.Request) {

	start := r.URL.Query()["start"]
	limit := r.URL.Query()["limit"]

	query := "select * from news"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(start) > 0 {
			query = query + " offset " + start[0]
		}
	}

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	news := []News{}
	for rows.Next() {
		p := News{}
		err := rows.Scan(&p.Id, &p.Datetime, &p.TitleRu, &p.TitleEn, &p.SubtitleRu, &p.SubtitleEn, &p.Dir, &p.ImgBack, &p.ImgBackfull, &p.ImagesCount, &p.VideosCount, &p.TextRu, &p.TextEn)
		if err != nil {
			fmt.Println(err)
			continue
		}
		news = append(news, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {

	defer db.Close()
	// db, err := sql.Open("sqlite3", "db.sqlite")
	// if err != nil {
	// 	panic(err)
	// }
	r := mux.NewRouter()
	// books = append(books, Book{ID: "1", Title: "Война и Мир", Author: &Author{Firstname: "Лев", Lastname: "Толстой"}})
	// books = append(books, Book{ID: "2", Title: "Преступление и наказание", Author: &Author{Firstname: "Фёдор", Lastname: "Достоевский"}})
	r.HandleFunc("/sales", getSales).Methods("GET")
	r.HandleFunc("/paintings", getPaintings).Methods("GET")
	r.HandleFunc("/threeds", getThreeds).Methods("GET")
	r.HandleFunc("/illustrations", getIllustrations).Methods("GET")
	r.HandleFunc("/news", getNews).Methods("GET")
	// r.HandleFunc("/books/{id}", getBook).Methods("GET")
	// r.HandleFunc("/books", createBook).Methods("POST")
	// r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
	// db.Close()
}
