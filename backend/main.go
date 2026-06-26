package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	store := NewPersonaStore()
	handler := NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/personas", handler)
	mux.Handle("/personas/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor escuchando en :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
