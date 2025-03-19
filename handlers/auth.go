package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"piscine/projet_BDE/models"
)

// Variables pour stocker les sessions et les utilisateurs temporairement
var sessions = make(map[string]string)   // session_id -> email
var users = make(map[string]models.User) // email -> user

// LoginPage affiche la page de connexion
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

// SignupPage affiche la page d'inscription
func SignupPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signup.html")
	tmpl.Execute(w, nil)
}

// Authenticate gère la connexion des utilisateurs
func Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Vérification basique (email Ynov et mot de passe)
	if !strings.HasSuffix(email, "@ynov.com") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, exists := users[email]
	if !exists || password != user.Password {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Création de session
	sessionID := models.GenerateSessionID()
	sessions[sessionID] = email

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Register gère l'inscription des utilisateurs
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validations basiques
	if name == "" || email == "" || password == "" || !strings.HasSuffix(email, "@ynov.com") {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	// Enregistrement simple
	users[email] = models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Logout déconnecte l'utilisateur
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	if cookie != nil {
		delete(sessions, cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// IsLoggedIn vérifie si l'utilisateur est connecté
func IsLoggedIn(sessionID string) bool {
	_, exists := sessions[sessionID]
	return exists
}

// GetUserFromSession récupère l'utilisateur correspondant à la session
func GetUserFromSession(sessionID string) models.User {
	email, exists := sessions[sessionID]
	if !exists {
		return models.User{}
	}
	return users[email]
}
