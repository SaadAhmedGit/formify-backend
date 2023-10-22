package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/SaadAhmedGit/forms/internal/database"
	"github.com/SaadAhmedGit/forms/internal/models"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload LoginDTO

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request",
		})
		return
	}

	db, err := database.NewDatabase()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "There was an issue with the server. Please try again at a later time",
		})
		return
	}

	user, err := models.FindUser(db, payload.Email, payload.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Email or password is incorrect",
		})
		return
	}

	if models.UserAuthorized(user.HashedPassword, payload.Password) == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Email or password is incorrect",
		})
		return
	}

	// Create JWT
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &LoginClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(env.JWT_SECRET))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
	})
}
