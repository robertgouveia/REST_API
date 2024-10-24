package main

import "log"

func main() {
	cfg := config{
		addr: ":3000",
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Printf("Server has started at %s", app.config.addr)
	log.Fatal(app.run(mux))
}
