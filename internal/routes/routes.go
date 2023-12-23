package routes

import (
	"net/http"

	"github.com/SaadAhmedGit/formify/internal/handlers"
	"github.com/gorilla/mux"
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
	mux.HandleFunc("/pre-signup", handlers.HandlePreSignUp).Methods("POST")
	mux.HandleFunc("/signup", handlers.HandleSignUp).Methods("GET") // TODO: Change to POST
	mux.HandleFunc("/login", handlers.HandleLogin).Methods("GET")

	return mux
}
