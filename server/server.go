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
	"github.com/rs/cors"
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

	offset := r.URL.Query()["offset"]
	limit := r.URL.Query()["limit"]

	query := "select s.*, b.base_ru, b.base_en from sales s join bases b on s.base_id = b.id"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(offset) > 0 {
			query = query + " offset " + offset[0]
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
		err := rows.Scan(&p.Id, &p.Dir, &p.Width, &p.Height, &p.Year, &p.Price, &p.NameRu, &p.NameEn, &p.BaseId, &p.StrId, &p.ImgCount, &p.BaseRu, &p.BaseEn)
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

	offset := r.URL.Query()["offset"]
	limit := r.URL.Query()["limit"]

	query := "select s.*, b.base_ru, b.base_en from paintings s join bases b on s.base_id = b.id"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(offset) > 0 {
			query = query + " offset " + offset[0]
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
		err := rows.Scan(&p.Id, &p.Dir, &p.Width, &p.Height, &p.Year, &p.NameRu, &p.NameEn, &p.BaseId, &p.StrId, &p.ImgCount, &p.BaseRu, &p.BaseEn)
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
	offset := r.URL.Query()["offset"]
	limit := r.URL.Query()["limit"]

	query := "select t.*, b.base_ru, b.base_en from threeds t join bases b on t.base_id = b.id"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(offset) > 0 {
			query = query + " offset " + offset[0]
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
		err := rows.Scan(&p.Id, &p.StrId, &p.Dir, &p.NameRu, &p.NameEn, &p.BaseId, &p.Year, &p.ImgCount, &p.BaseRu, &p.BaseEn)
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
	offset := r.URL.Query()["offset"]
	limit := r.URL.Query()["limit"]

	query := "select i.*, b.base_ru, b.base_en from illustrations i join bases b on i.base_id = b.id"

	if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(offset) > 0 {
			query = query + " offset " + offset[0]
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
		err := rows.Scan(&p.Id, &p.StrId, &p.Dir, &p.NameRu, &p.NameEn, &p.BaseId, &p.Year, &p.ImgCount, &p.BaseRu, &p.BaseEn)
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

	id := r.URL.Query()["id"]
	offset := r.URL.Query()["offset"]
	limit := r.URL.Query()["limit"]
	id_ne := r.URL.Query()["id_ne"]

	query := "select * from news"

	if len(id) > 0 {
		query = query + " where id = " + id[0]
	} else if len(id_ne) > 0 && len(limit) > 0 {
		query = query + " where id != " + id_ne[0] + " limit " + limit[0]
	} else if len(limit) > 0 {
		query = query + " limit " + limit[0]

		if len(offset) > 0 {
			query = query + " offset " + offset[0]
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
	r := mux.NewRouter()
	// books = append(books, Book{ID: "1", Title: "Война и Мир", Author: &Author{Firstname: "Лев", Lastname: "Толстой"}})
	// books = append(books, Book{ID: "2", Title: "Преступление и наказание", Author: &Author{Firstname: "Фёдор", Lastname: "Достоевский"}})
	// r.HandleFunc("/books/{id}", getBook).Methods("GET")
	// r.HandleFunc("/books", createBook).Methods("POST")
	// r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application
		},
	})

	r.HandleFunc("/sales", getSales).Methods("GET")
	r.HandleFunc("/paintings", getPaintings).Methods("GET")
	r.HandleFunc("/threeds", getThreeds).Methods("GET")
	r.HandleFunc("/illustrations", getIllustrations).Methods("GET")
	r.HandleFunc("/news", getNews).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", corsOpts.Handler(r)))
}
