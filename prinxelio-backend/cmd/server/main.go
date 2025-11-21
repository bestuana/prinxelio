package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"prinxelio.com/backend/pkg/api"
	"prinxelio.com/backend/pkg/database"
)

func main() {
	fmt.Println("Starting application...")
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err = database.SeedCategories(db); err != nil {
		log.Fatalf("Failed to seed categories: %v", err)
	}
	if err = database.SeedProducts(db); err != nil {
		log.Fatalf("Failed to seed products: %v", err)
	}

	srv := api.NewServer(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3020"
	}
	fmt.Printf("HTTP server running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, srv.Mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
