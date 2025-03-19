package handlers

import (
	"html/template"
	"net/http"

	"piscine/projet_BDE/models"
)

// On va créer un système d'authentification pour le site
// Pour cela on crée des  variables pour stocker les sessions et les utilisateurs temporairement
var sessions = make(map[string]string)
var users = make(map[string]models.User)

// La fonction LoginPage va afficher la page de connexion
// Pour cela on utilise la fonction ParseFiles qui va lire le fichier login.html
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/login.html")
	tmpl.Execute(w, nil)
}

// SignupPage affiche la page d'inscription
// ParseFiles qui lit le fichier signup.html
func SignupPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/signup.html")
	tmpl.Execute(w, nil)
}
