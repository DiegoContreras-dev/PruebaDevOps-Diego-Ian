package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// @Summary Listar personas
// @Description Retorna todas las personas registradas
// @Produce json
// @Success 200 {array} Persona
// @Router /personas [get]
func swaggerGetPersonas() {}

// @Summary Agregar persona
// @Description Crea un nuevo registro de persona
// @Accept json
// @Produce json
// @Param persona body Persona true "Datos de la persona"
// @Success 201 {object} Persona
// @Failure 400 {object} map[string]string
// @Router /personas [post]
func swaggerAddPersona() {}

// @Summary Eliminar persona
// @Description Elimina una persona por su RUT
// @Produce json
// @Param rut path string true "RUT de la persona"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /personas/{rut} [delete]
func swaggerDeletePersona() {}

type Handler struct {
	store *PersonaStore
}

func NewHandler(store *PersonaStore) *Handler {
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
	personas := h.store.GetAll()
	if err := json.NewEncoder(w).Encode(personas); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
	}
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
	h.store.Add(p)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) deletePersona(w http.ResponseWriter, r *http.Request, rut string) {
	if !h.store.Delete(rut) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "persona no encontrada"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "persona eliminada"})
}
