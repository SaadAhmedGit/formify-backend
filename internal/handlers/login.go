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
		RespondWithJSON(w, http.StatusBadRequest, JSONResponse{
			"error": "Invalid request",
		})
		return
	}

	db, err := database.NewDatabase()
	if err != nil {
		RespondWithJSON(w, http.StatusInternalServerError, JSONResponse{
			"error": "There was an issue with the server. Please try again at a later time",
		})
		return
	}

	user, err := models.FindUser(db, payload.Email)
	if err != nil {
		RespondWithJSON(w, http.StatusUnauthorized, JSONResponse{
			"error": "The user does not exist",
		})
		return
	}

	if !models.UserAuthorized(user.HashedPassword, payload.Password) {
		RespondWithJSON(w, http.StatusUnauthorized, JSONResponse{
			"error": "Email or password is incorrect",
		})
		return
	}

	// Create JWT Claims
	expirationTime := time.Now().Add(time.Duration(env.TOKEN_VALIDITY_DAYS) * 24 * time.Hour)
	claims := LoginClaims{
		Email: payload.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	tokenString, err := createJWT(claims, env.JWT_SECRET)

	RespondWithJSON(w, http.StatusOK, JSONResponse{
		"token": tokenString,
	})
}
