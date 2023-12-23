package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SaadAhmedGit/formify/internal/config"
	"github.com/SaadAhmedGit/formify/internal/routes"
)

func main() {

	env, err := config.Env()
	if err != nil {
		log.Fatalln("Failed to parse environment variables.")
	}

	router := routes.NewRouter()

	fmt.Printf("Listening on port %d...\n", env.PORT)
	addr := fmt.Sprintf(":%d", env.PORT)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Failed to listen on port %d", env.PORT)
	}
}
