package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	store PersonaRepository
}

func NewHandler(store PersonaRepository) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path == "/personas" || r.URL.Path == "/personas/" {
		switch r.Method {
		case http.MethodGet:
			h.getPersonas(w, r)
		case http.MethodPost:
			h.addPersona(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
		return
	}

	if strings.HasPrefix(r.URL.Path, "/personas/") {
		rut := strings.TrimPrefix(r.URL.Path, "/personas/")
		if rut == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "rut requerido"})
			return
		}
		switch r.Method {
		case http.MethodDelete:
			h.deletePersona(w, r, rut)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "ruta no encontrada"})
}

func (h *Handler) getPersonas(w http.ResponseWriter, r *http.Request) {
	personas, err := h.store.GetAll()
	if err != nil {
		log.Printf("ERROR GetAll: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error interno del servidor"})
		return
	}
	json.NewEncoder(w).Encode(personas)
}

func (h *Handler) addPersona(w http.ResponseWriter, r *http.Request) {
	var p Persona
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "body inválido"})
		return
	}
	if p.Nombre == "" || p.RUT == "" || p.FechaNacimiento == "" || p.Ciudad == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "todos los campos son requeridos"})
		return
	}
	// Garantizar que gustos nunca sea null en la respuesta JSON
	if p.Gustos == nil {
		p.Gustos = []string{}
	}
	if err := h.store.Add(p); err != nil {
		log.Printf("ERROR Add: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error interno del servidor"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) deletePersona(w http.ResponseWriter, r *http.Request, rut string) {
	found, err := h.store.Delete(rut)
	if err != nil {
		log.Printf("ERROR Delete: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "error interno del servidor"})
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "persona no encontrada"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "persona eliminada"})
}
