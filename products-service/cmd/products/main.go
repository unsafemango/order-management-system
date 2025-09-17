package main

import (
	"log"
	"net/http"

	api "yray.com/products-service/internal/api"
	"yray.com/products-service/internal/env"
	"yray.com/products-service/internal/repo"
)

func main() {
	dsn := env.GetString("DB_DSN", "postgres://app:apppass@localhost:5432/retail?sslmode=disable")
	addr := env.GetString("HTTP_ADDR", ":8081")

	db, err := repo.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := repo.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	r := api.NewRouter(db)
	log.Printf("products service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
