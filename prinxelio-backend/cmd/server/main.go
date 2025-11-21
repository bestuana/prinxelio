package main

import (
    "fmt"
    "log"
    "net/http"

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
    fmt.Println("HTTP server running on :3010")
    if err := http.ListenAndServe(":3010", srv.Mux); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
