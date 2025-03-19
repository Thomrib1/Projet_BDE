package models

import (
	"crypto/rand"
	"encoding/base64"
)

type User struct {
	Name     string
	Email    string
	Password string
}

// La Fonction GenerateSessionID va générer un identifiant de session aléatoire
// qui sera utilisé pour identifier les utilisateurs connectés.
// La fonction utilise la bibliothèque crypto/rand pour générer des octets aléatoires
// et les encode en base64 pour obtenir une chaîne de caractères.
func GenerateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}
