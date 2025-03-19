package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"
)

type User struct {
	Name     string
	Email    string
	Password string
}

// On va créer un système d'authentification pour le site
// Pour cela on crée des variables pour stocker les sessions et les utilisateurs temporairement
var sessions = make(map[string]string)
var users = make(map[string]User) // Changé de models.User à User

// GenerateSessionID génère un identifiant de session aléatoire
func GenerateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// La fonction LoginPage va afficher la page de connexion
// Pour cela on utilise la fonction ParseFiles qui va lire le fichier login.html
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

// SignupPage affiche la page d'inscription
func SignupPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signup.html")
	tmpl.Execute(w, nil)
}

// La func Authenticate va vérifier les informations de connexion
// Si les informations sont correctes, on crée une session pour l'utilisateur
// On redirige ensuite l'utilisateur vers le tableau de bord
// Sinon, on le redirige vers la page de connexion
func Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if !strings.HasSuffix(email, "@ynov.com") {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, exists := users[email]
	if !exists || password != user.Password {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	sessionID := GenerateSessionID()
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
