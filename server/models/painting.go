package models

import (
	"database/sql"

	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Painting struct {
	Id           int    `json:"id"`
	Dir          string `json:"dir"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Year         int    `json:"year"`
	NameRu       string `json:"name_ru"`
	NameEn       string `json:"name_en"`
	BaseId       int    `json:"base_id"`
	StrId        string `json:"str_id"`
	ImgCount     int    `json:"img_count"`
	Descr        string `json:"descr"`
	MaterialsIds []int  `json:"materials_ids"`
}

type PaintingWithBase struct {
	Id       int    `json:"id"`
	Dir      string `json:"dir"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Year     int    `json:"year"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   int    `json:"base_id"`
	StrId    string `json:"str_id"`
	ImgCount int    `json:"img_count"`
	Descr    string `json:"descr"`
	BaseRu   string `json:"base_ru"`
	BaseEn   string `json:"base_en"`
}

// Save painting to database
func (p *Painting) Save(db *sql.DB) error {
	// imagePathsJSON, err := json.Marshal(p.ImagePaths)
	// if err != nil {
	// 	return err
	// }

	res, err := db.Exec(`
		INSERT INTO paintings
		(dir, width, height, year, name_ru, name_en, base_id, str_id, img_count, descr)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Dir, p.Width, p.Height, p.Year, p.NameRu, p.NameEn, p.BaseId, strings.Replace(p.NameEn, " ", "_", -1), p.ImgCount, p.Descr)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = int(id)

	return nil
}
