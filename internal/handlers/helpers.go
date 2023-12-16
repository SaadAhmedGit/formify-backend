package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type JSONResponse map[string]interface{}

func RespondWithJSON(w http.ResponseWriter, code int, payload JSONResponse) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func createJWT(claims jwt.Claims, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func sendVerificationEmail(payload PreSignUpDTO, tokenString string) error {
	from := mail.NewEmail(env.APP_NAME, env.EMAIL_FROM)
	subject := fmt.Sprintf("Verify your email for %s", env.APP_NAME)
	to := mail.NewEmail(payload.FullName, payload.Email)
	plainTextContent := fmt.Sprintf("Please verify your email by clicking on the following link: %s/signup?token=%s", env.DEV_SERVER_URL, tokenString)
	htmlContent := fmt.Sprintf("Please verify your email by clicking on the following link: %s/signup?token=%s", env.DEV_SERVER_URL, tokenString)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(env.SENDGRID_API_KEY)
	_, err := client.Send(message)
	if err != nil {
		return err
	}
	return nil
}
