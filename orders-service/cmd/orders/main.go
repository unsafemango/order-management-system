package main

import (
	"log"
	"net/http"

	"yray.com/orders-service/internal/api"
	"yray.com/orders-service/internal/env"
	"yray.com/orders-service/internal/repo"
)

func main() {
	dsn := env.GetString("DB_DSN", "postgres://app:apppass@localhost:5432/retail?sslmode=disable")
	addr := env.GetString("HTTP_ADDR", ":8080")
	productsURL := env.GetString("PRODUCTS_URL", "http://localhost:8081/v1")

	db, err := repo.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := repo.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	r := api.NewRouter(db, productsURL)

	log.Printf("orders service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
