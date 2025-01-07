import json
import sqlite3

connection = sqlite3.connect('d:/projects/vera-art/server/db.sqlite')
cursor = connection.cursor()

# Open and read the JSON file
with open('d:/projects/vera-art/server/db.json', encoding="utf8") as file:
    data = json.load(file)

for i in data['news']:
    cursor.execute('INSERT INTO news (id, datetime, title_ru, title_en, subTitle_ru, subTitle_en, dir, img_back, img_backfull, imagescount, videoscount, text_ru, text_en) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',
                   (i['id'], i['datetime'], i['title_ru'], i['title_en'], i['subTitle_ru'], i['subTitle_en'], i['dir'], i['img_back'], i['img_backfull'], i['imagescount'], i['videoscount'], i['text_ru'], i['text_en']))

connection.commit()
connection.close()
