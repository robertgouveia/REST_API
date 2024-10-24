package main

import (
	"log"

	"github.com/robertgouveia/social/internal/env"
	"github.com/robertgouveia/social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Printf("Server has started at %s", app.config.addr)
	log.Fatal(app.run(mux))
}
