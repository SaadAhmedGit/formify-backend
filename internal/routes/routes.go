package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/SaadAhmedGit/formify/internal/handlers"
)

var (
	router = initRouter()
)

func NewRouter() http.Handler {
	return router
}

func initRouter() *mux.Router {
	mux := mux.NewRouter()

	// Auth
	mux.HandleFunc("/api/auth/login", handlers.HandleLogin).Methods("POST")
	mux.HandleFunc("/api/auth/verify", handlers.HandleVerify).Methods("POST")
	mux.HandleFunc("/api/auth/logout", handlers.HandleLogout).Methods("POST")
	mux.HandleFunc("/api/auth/signup", handlers.HandleSignUp).Methods("GET") // TODO: Change to POST
	mux.HandleFunc("/api/auth/pre-signup", handlers.HandlePreSignUp).Methods("POST")

	return mux
}
