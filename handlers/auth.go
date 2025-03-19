package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"strings"

	"projet_BDE/models"
)

type User struct {
	Name     string
	Email    string
	Password string
}

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

// Page d'inscription
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

	user, err := models.GetUserByEmail(email)
	if err != nil || !models.CheckPassword(password, user.Password) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Création de la session en base de données
	sessionID := GenerateSessionID()
	_, err = models.DB.Exec("INSERT INTO sessions (session_id, user_email) VALUES (?, ?)", sessionID, email)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Création du cookie de session
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

// Inscription des utilisateurs
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

	err := models.CreateUser(name, email, password)
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// La fonction Logout va supprimer la session de l'utilisateur
// On redirige ensuite l'utilisateur vers la page d'accueil
// Et setcookie va supprimer le cookie de session
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		// Suppression de la session en base
		_, _ = models.DB.Exec("DELETE FROM sessions WHERE session_id = ?", cookie.Value)

		// Suppression du cookie
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
	var email string
	err := models.DB.QueryRow("SELECT user_email FROM sessions WHERE session_id = ?", sessionID).Scan(&email)
	return err == nil
}

// GetUserFromSession récupère l'utilisateur à partir de la session
// On utilise l'email stocké dans la session pour récupérer l'utilisateur
func GetUserFromSession(sessionID string) models.User {
	var user models.User
	err := models.DB.QueryRow("SELECT id, name, email FROM users WHERE email = (SELECT user_email FROM sessions WHERE session_id = ?)", sessionID).
		Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return models.User{}
	}
	return user
}
