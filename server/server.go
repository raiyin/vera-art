package main

import (
	"fmt"
	"net/http"
	"os"
)

type Part struct {
	NameRU     string `json:"name_ru"`
	NameEN     string `json:"name_en"`
	BaseRU     string `json:"base_ru"`
	BaseEN     string `json:"base_en"`
	MaterialRU string `json:"material_ru"`
	MaterialEN string `json:"material_en"`
	Size       string `json:"size"`
}

type Sale struct {
	ID         string `json:"id"`
	Dir        string `json:"dir"`
	NameRu     string `json:"name_ru"`
	NameEn     string `json:"name_en"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	BaseRu     string `json:"base_ru"`
	BaseEn     string `json:"base_en"`
	MaterialRu string `json:"material_ru"`
	MaterialEn string `json:"material_en"`
	Size       string `json:"size"`
	Year       int    `json:"year"`
	Price      int    `json:"price"`
	Parts      []Part `json:"parts"`
}

func main() {

	jsonFile, err := os.Open("db.json")
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	http.ListenAndServe(":80", nil)
}
