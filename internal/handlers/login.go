package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/SaadAhmedGit/formify/internal/database"
	"github.com/SaadAhmedGit/formify/internal/models"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginClaims struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
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
		Email:    payload.Email,
		FullName: user.FullName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	tokenString, err := createJWT(claims, env.JWT_SECRET)

	// Set token cookies
	tokenCookie := http.Cookie{
		Name:  "token",
		Value: tokenString,

		Expires:  expirationTime,
		HttpOnly: false,

		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}
	tokenSetCookie := http.Cookie{
		Name:  "token_set",
		Value: "true",

		Expires:  expirationTime,
		HttpOnly: false,

		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}

	http.SetCookie(w, &tokenCookie)
	http.SetCookie(w, &tokenSetCookie)
	RespondWithJSON(w, http.StatusOK, JSONResponse{
		"user": map[string]string{
			"email":     user.Email,
			"full_name": user.FullName,
		},
	})
}

func HandleVerify(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, JSONResponse{
			"error": "Invalid request",
		})
		return
	}

	token, err := verifyToken(tokenCookie.Value)
	if err != nil {
		RespondWithJSON(w, http.StatusUnauthorized, JSONResponse{
			"error": "Invalid token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		RespondWithJSON(w, http.StatusUnauthorized, JSONResponse{
			"error": "Invalid token",
		})
		return
	}

	RespondWithJSON(w, http.StatusOK, JSONResponse{
		"user": map[string]string{
			"email":     claims["email"].(string),
			"full_name": claims["full_name"].(string),
		},
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	tokenSetCookie, err := r.Cookie("token_set")
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, JSONResponse{
			"error": "Invalid request",
		})
		return
	}

	tokenCookie.Expires = time.Now().AddDate(0, 0, -1)
	tokenSetCookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, tokenCookie)
	http.SetCookie(w, tokenSetCookie)

	RespondWithJSON(w, http.StatusOK, JSONResponse{
		"message": "Successfully logged out",
	})
}
