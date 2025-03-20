package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}

	// Créer la table users si elle n'existe pas
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
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func CreateUser(name, email, password string) error {
	// Hacher le mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insérer l'utilisateur dans la base de données
	_, err = DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, string(hashedPassword))
	return err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return user, nil
	}
	return user, err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// La Fonction GenerateSessionID va générer un identifiant de session aléatoire
// qui sera utilisé pour identifier les utilisateurs connectés.
// La fonction utilise la bibliothèque crypto/rand pour générer des octets aléatoires
// et les encode en base64 pour obtenir une chaîne de caractères.
func GenerateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
