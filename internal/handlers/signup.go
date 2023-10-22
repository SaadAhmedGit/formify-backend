package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

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
		log.Println("Failed to parse request body.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request"))
		return
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Println("Failed to connect to database.")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an issue with the server. Please try at a later time"))
		return
	}

	// Check if email is already taken
	UserExists, err := models.UserExists(db, payload.Email)
	if err != nil {
		log.Println("Failed to check if user exists.")
	}
	if UserExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("An account already exists on this email."))
		return
	}

	//Create JWT
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &preSignUpClaims{
		Email:    payload.Email,
		Password: payload.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(env.JWT_ACCOUNT_ACTIVATION))

	// Send verification email
	from := mail.NewEmail(env.APP_NAME, env.EMAIL_FROM)
	subject := fmt.Sprintf("Verify your email for %s", env.APP_NAME)
	to := mail.NewEmail(payload.FullName, payload.Email)
	plainTextContent := fmt.Sprintf("Please verify your email by clicking on the following link: %s/signup?token=%s", env.DEV_SERVER_URL, tokenString)
	htmlContent := fmt.Sprintf("Please verify your email by clicking on the following link: %s/signup?token=%s", env.DEV_SERVER_URL, tokenString)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(env.SENDGRID_API_KEY)
	_, err = client.Send(message)
	if err != nil {
		log.Println("Failed to send verification email.")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verification email sent. Please check your inbox."))
}

func HandleSignUp(w http.ResponseWriter, r *http.Request) {

	var claims preSignUpClaims

	// TODO: Swap to body when frontend is made
	query := r.URL.Query()
	tokenParam := query.Get("token")

	// Verify JWT token
	token, err := jwt.ParseWithClaims(tokenParam, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(env.JWT_ACCOUNT_ACTIVATION), nil
	})
	if err != nil || !token.Valid {
		log.Println("Invalid JWT token.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Your verification link is invalid."))
		return
	}

	// Extract JWT token claims
	var newUser models.User
	newUser.Email = claims.Email
	newUser.HashedPassword = claims.Password

	// Create user
	db, err := database.NewDatabase()
	if err != nil {
		log.Println("Failed to connect to database.")
	}

	err = models.CreateUser(db, newUser)
	if err != nil {
		log.Println("Failed to create user.")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully."))
}
