package handlers

import (
	"html/template"
	"net/http"
)

// HomePage va nous afficher la page d'accueil
// Si l'utilisateur est connecté, il pourra voir un message de bienvenue
// Sinon, il sera invité à se connecter ou à s'inscrire
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	var isLoggedIn bool
	cookie, err := r.Cookie("session")
	if err == nil {
		isLoggedIn = IsLoggedIn(cookie.Value)
	}

	data := struct {
		LoggedIn bool
	}{
		LoggedIn: isLoggedIn,
	}

	tmpl, _ := template.ParseFiles("templates/home.html")
	tmpl.Execute(w, data)
}
