package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"projet_BDE/handlers"
	"projet_BDE/models"
)

func main() {
	models.InitDB("bde.db")
	fmt.Println("Base de données initialisée avec succès")

	ensureDirectoryExists("static")
	ensureDirectoryExists("static/images")

	http.Handle("/templates/images/", http.StripPrefix("/templates/images/", http.FileServer(http.Dir("templates/images"))))
	http.HandleFunc("/", serveHTMLFile("templates/home.html"))
	http.HandleFunc("/equipe", serveHTMLFile("templates/equipe.html"))
	http.HandleFunc("/contact", serveHTMLFile("templates/contact.html"))
	http.HandleFunc("/calendrier", serveHTMLFile("templates/calendrier.html"))
	http.HandleFunc("/login", serveHTMLFile("templates/login.html"))
	http.HandleFunc("/admin", serveHTMLFile("templates/admin.html"))

	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/auth", handlers.Authenticate)
	http.HandleFunc("/dashboard", handlers.DashboardPage)
	http.HandleFunc("/logout", handlers.Logout)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))

	fmt.Println("Serveur démarré sur http://localhost:8080")
	fmt.Println("Fichiers statiques servis depuis le dossier: " + getCurrentDir() + "/static")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHTMLFile(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

func ensureDirectoryExists(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Le dossier %s n'existe pas, création...\n", dir)
		os.MkdirAll(dir, 0755)
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
