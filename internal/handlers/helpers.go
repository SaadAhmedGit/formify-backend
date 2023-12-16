package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
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
