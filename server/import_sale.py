import json
import sqlite3

def import_sale(data, connection):
    cursor.execute('INSERT INTO Users (username, email, age) VALUES (?, ?, ?)', ('newuser', 'newuser@example.com', 28))

    # Сохраняем изменения и закрываем соединение
    connection.commit()

connection = sqlite3.connect('db.sqlite')
cursor = connection.cursor()


# Open and read the JSON file
with open('db.json', encoding="utf8") as file:
    data = json.load(file)

for i in data['sale']:
    import_sale(i)

connection.close()
