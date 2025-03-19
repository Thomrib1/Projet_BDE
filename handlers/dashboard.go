package handlers

import (
	"html/template"
	"net/http"
)

// La fonction DashboardPage va afficher la page de tableau de bord
// Si l'utilisateur n'est pas connecté, il sera redirigé vers la page de connexion
// Sinon, il pourra voir son nom et son adresse e-mail
func DashboardPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || !IsLoggedIn(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user := GetUserFromSession(cookie.Value)
	data := struct {
		Name  string
		Email string
	}{
		Name:  user.Name,
		Email: user.Email,
	}

	tmpl, _ := template.ParseFiles("templates/dashboard.html")
	tmpl.Execute(w, data)
}
