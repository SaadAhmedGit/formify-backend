package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/SaadAhmedGit/forms/internal/config"
	"github.com/SaadAhmedGit/forms/internal/database"
	"github.com/SaadAhmedGit/forms/internal/models"
)

type UserDTO struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpDTO = UserDTO

type PreSignUpDTO = UserDTO
type preSignUpClaims struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var env, _ = config.Env()

func HandlePreSignUp(w http.ResponseWriter, r *http.Request) {

	var payload PreSignUpDTO

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

	// Check if email is already taken
	_, err = models.FindUser(db, payload.Email)
	if err != nil {
		RespondWithJSON(w, http.StatusBadRequest, JSONResponse{
			"error": "Email is already taken",
		})
		return
	}

	//Create JWT
	expirationTime := time.Now().Add(time.Duration(env.PRESIGNUP_TOKEN_VALIDITY_HOURS) * time.Hour)
	claims := &preSignUpClaims{
		Email:    payload.Email,
		Password: payload.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tokenString, err := createJWT(claims, env.JWT_ACCOUNT_ACTIVATION)

	err = sendVerificationEmail(payload, tokenString)
	if err != nil {
		log.Printf("Failed to send verification email: %v", err)
		RespondWithJSON(w, http.StatusInternalServerError, JSONResponse{
			"error": "There was an issue with the server. Please try again at a later time",
		})
		return
	}

	RespondWithJSON(w, http.StatusOK, JSONResponse{
		"message": "Verification email sent.",
	})
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {

	var claims preSignUpClaims

	// TODO: Swap query with body when frontend is made
	query := r.URL.Query()
	tokenParam := query.Get("token")

	// Verify JWT token
	token, err := jwt.ParseWithClaims(tokenParam, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.JWT_ACCOUNT_ACTIVATION), nil
	})
	if err != nil || !token.Valid {
		RespondWithJSON(w, http.StatusBadRequest, JSONResponse{
			"error": "Your verification link is invalid or expired. Please try again.",
		})
		return
	}

	// Extract JWT token claims
	var newUser models.User
	newUser.Email = claims.Email
	newUser.HashedPassword = claims.Password

	db, err := database.NewDatabase()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		RespondWithJSON(w, http.StatusInternalServerError, JSONResponse{
			"error": "There was an issue with the server. Please try again at a later time",
		})
		return
	}

	// Create user
	err = models.CreateUser(db, newUser)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		RespondWithJSON(w, http.StatusInternalServerError, JSONResponse{
			"error": "There was an issue with the server. Please try again at a later time",
		})
		return
	}

	RespondWithJSON(w, http.StatusOK, JSONResponse{
		"message": "User created successfully.",
	})
}
