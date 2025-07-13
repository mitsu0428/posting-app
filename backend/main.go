package main

import (
	"log"
	"net/http"
	"os"
	"posting-app/di"
	"posting-app/infrastructure"
)

func main() {
	container := di.BuildContainer()

	err := container.Invoke(func(params di.DIParams) error {
		if err := infrastructure.RunMigrations(params.Database); err != nil {
			return err
		}

		router := params.Handler.SetupRoutes()

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("Server starting on port %s", port)
		return http.ListenAndServe(":"+port, router)
	})

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}