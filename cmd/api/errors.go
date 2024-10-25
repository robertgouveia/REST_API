package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) conflict(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusNotFound, "Resource not found")
}
