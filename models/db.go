package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB va initialiser la base de données avec le nom du fichier
// Si la base de données n'existe pas, elle sera créée
// On crée également les tables users et sessions
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`
	_, err = DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table users: %v", err)
	}

	createSessionsTable := `
    CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        session_id TEXT UNIQUE NOT NULL,
        user_email TEXT NOT NULL,
        FOREIGN KEY(user_email) REFERENCES users(email) ON DELETE CASCADE
    );`
	_, err = DB.Exec(createSessionsTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table sessions: %v", err)
	}
}
