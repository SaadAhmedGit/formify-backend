package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SaadAhmedGit/formify/internal/config"
	"github.com/SaadAhmedGit/formify/internal/routes"
	"github.com/gorilla/handlers"
)

func main() {

	env, err := config.Env()
	if err != nil {
		log.Fatalln("Failed to parse environment variables.")
	}

	router := routes.NewRouter()

	// CORS
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	credentialsOK := handlers.AllowCredentials()
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	corsRouter := handlers.CORS(headersOK, originsOK, credentialsOK, methodsOK)(router)
	loggedRouter := handlers.LoggingHandler(os.Stderr, corsRouter)

	fmt.Printf("Listening on port %d...\n", env.PORT)
	addr := fmt.Sprintf(":%d", env.PORT)

	err = http.ListenAndServe(addr, loggedRouter)
	if err != nil {
		log.Fatalf("Failed to listen on port %d", env.PORT)
	}
}
