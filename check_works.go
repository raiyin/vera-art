package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT str_id, name_ru, name_en FROM works LIMIT 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Existing works:")
	for rows.Next() {
		var str_id, name_ru, name_en string
		err := rows.Scan(&str_id, &name_ru, &name_en)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("str_id: %s, name_ru: %s, name_en: %s\n", str_id, name_ru, name_en)
	}
}
