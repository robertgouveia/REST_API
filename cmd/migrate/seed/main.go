package main

import (
	"log"

	"github.com/robertgouveia/social/internal/db"
	"github.com/robertgouveia/social/internal/store"
)

func main() {
	conn, err := db.New("postgres://user:adminpassword@localhost/social?sslmode=disable", 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	store := store.NewStorage(conn)

	db.Seed(store)
}
