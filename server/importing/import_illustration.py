import json
import sqlite3

def get_base_id(base: str) -> int:
    base_dict = {
        "холст": 1,
        "бумага": 2,
        "холст на картоне": 3,
        "картон": 4,
        "Adobe Photoshop": 5,
        "Autodesk 3ds Max": 6
        }
    if base in base_dict:
        return base_dict[base]
    else:
        print(base)
        return 1

def get_material_id(material: str) -> int:
    material_dict = {
        "масло":1,
        "акварель":2,
        "текстурная паста":3,
        "уголь":4,
        "акрил":5,
        "гуашь":6,
        "карандаш":7,
        "фломастер":8,
        "пастель":9,
        "чернила":10
    }

    if material in material_dict:
        return material_dict[material]
    else:
        print(material)
        return 1

connection = sqlite3.connect('d:/projects/vera-art/server/db.sqlite')
cursor = connection.cursor()


# Open and read the JSON file
with open('d:/projects/vera-art/server/db.json', encoding="utf8") as file:
    data = json.load(file)

for i in data['illustration']:
    cursor.execute('INSERT INTO illustration (str_id, dir, year, name_ru, name_en, base_id) VALUES (?, ?, ?, ?, ?, ?)',
                   (i['id'], i['dir'], i['year'], i['name_ru'], i['name_en'], get_base_id(i['base_ru'])))


connection.commit()
connection.close()
