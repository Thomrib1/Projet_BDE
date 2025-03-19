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
// Sinon on le redirige vers la page de connexion
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
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if name == "" || email == "" || password == "" || !strings.HasSuffix(email, "@ynov.com") {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	users[email] = User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// La fonction Logout va supprimer la session de l'utilisateur
// On redirige ensuite l'utilisateur vers la page d'accueil
// Et setcookie va supprimer le cookie de session
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
// On vérifie si la session existe dans la map sessions
func IsLoggedIn(sessionID string) bool {
	_, exists := sessions[sessionID]
	return exists
}

// GetUserFromSession récupère l'utilisateur à partir de la session
// On utilise l'email stocké dans la session pour récupérer l'utilisateur
func GetUserFromSession(sessionID string) User {
	email, exists := sessions[sessionID]
	if !exists {
		return User{}
	}
	return users[email]
}
