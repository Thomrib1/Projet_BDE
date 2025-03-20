package models

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table users: %v", err)
	}

	createSessionsTableQuery := `
    CREATE TABLE IF NOT EXISTS sessions (
        session_id TEXT PRIMARY KEY,
        user_email TEXT NOT NULL,
        FOREIGN KEY (user_email) REFERENCES users(email)
    );`
	_, err = DB.Exec(createSessionsTableQuery)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table sessions: %v", err)
	}
}
