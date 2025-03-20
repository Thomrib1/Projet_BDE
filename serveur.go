package main

import (
	"fmt"
	"log"
	"net/http"
	"projet_BDE/handlers"
	"projet_BDE/models"
)

func main() {
	models.InitDB("bde.db")
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/signup", handlers.SignupPage)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/auth", handlers.Authenticate)
	http.HandleFunc("/dashboard", handlers.DashboardPage)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
