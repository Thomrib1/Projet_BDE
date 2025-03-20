package models

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func CreateUser(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, string(hashedPassword))
	return err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	return user, err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
