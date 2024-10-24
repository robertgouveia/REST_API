package main

import (
	"log"

	"github.com/robertgouveia/social/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Printf("Server has started at %s", app.config.addr)
	log.Fatal(app.run(mux))
}
