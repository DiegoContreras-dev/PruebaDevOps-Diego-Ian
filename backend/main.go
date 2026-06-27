package main

import (
	"log"
	"net/http"
	"os"
	docs "personas-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	if azureHost := os.Getenv("WEBSITE_HOSTNAME"); azureHost != "" {
		docs.SwaggerInfo.Host = azureHost
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	store := NewPersonaStore()
	store.Add(Persona{Nombre: "Juan Pérez", RUT: "21614199-2", FechaNacimiento: "1990-05-20", Ciudad: "Santiago", Gustos: []string{"fútbol", "pizza", "películas"}})
	store.Add(Persona{Nombre: "Ana López", RUT: "18404852-4", FechaNacimiento: "1995-03-15", Ciudad: "Coquimbo", Gustos: []string{"lectura", "sushi", "videojuegos"}})
	store.Add(Persona{Nombre: "Diego Contreras", RUT: "12824745-k", FechaNacimiento: "2000-01-01", Ciudad: "La Serena", Gustos: []string{"ciberseguridad", "golang", "música"}})
	handler := NewHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/personas", handler)
	mux.Handle("/personas/", handler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor escuchando en :%s", port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
