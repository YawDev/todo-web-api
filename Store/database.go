package Store

import sqlite "todo-web-api/Store/Sqlite"

func DbConnection() {
	sqlite.Connect()
}
