package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestHandler usa MemoryStore para tests unitarios sin DB real
func newTestHandler() *Handler {
	return NewHandler(NewMemoryStore())
}

func TestGetPersonas_Empty(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/personas", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var personas []Persona
	if err := json.NewDecoder(w.Body).Decode(&personas); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(personas) != 0 {
		t.Fatalf("expected empty list, got %d items", len(personas))
	}
}

func TestAddPersona_OK(t *testing.T) {
	h := newTestHandler()
	body := `{"nombre":"Juan Pérez","rut":"12345678-9","fecha_nacimiento":"1990-05-20","ciudad":"Santiago","gustos":["fútbol","pizza"]}`
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var p Persona
	if err := json.NewDecoder(w.Body).Decode(&p); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(p.Gustos) != 2 {
		t.Fatalf("expected 2 gustos, got %d", len(p.Gustos))
	}
}

func TestAddPersona_InvalidBody(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAddPersona_MissingFields(t *testing.T) {
	h := newTestHandler()
	body := `{"nombre":"Juan Pérez"}`
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeletePersona_OK(t *testing.T) {
	h := newTestHandler()
	_ = h.store.Add(Persona{Nombre: "Juan Pérez", RUT: "12345678-9", FechaNacimiento: "1990-05-20", Ciudad: "Santiago", Gustos: []string{"fútbol"}})
	req := httptest.NewRequest(http.MethodDelete, "/personas/12345678-9", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeletePersona_NotFound(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodDelete, "/personas/99999999-9", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGetPersonas_AfterAdd(t *testing.T) {
	h := newTestHandler()
	body := `{"nombre":"Ana López","rut":"98765432-1","fecha_nacimiento":"1995-03-15","ciudad":"Valparaíso","gustos":["lectura","sushi"]}`
	postReq := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString(body))
	h.ServeHTTP(httptest.NewRecorder(), postReq)

	getReq := httptest.NewRequest(http.MethodGet, "/personas", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, getReq)

	var personas []Persona
	json.NewDecoder(w.Body).Decode(&personas)
	if len(personas) != 1 {
		t.Fatalf("expected 1 persona, got %d", len(personas))
	}
	if personas[0].RUT != "98765432-1" {
		t.Fatalf("unexpected RUT: %s", personas[0].RUT)
	}
	if len(personas[0].Gustos) != 2 {
		t.Fatalf("expected 2 gustos, got %d", len(personas[0].Gustos))
	}
}

func TestAddPersona_SinGustos_RetornaSliceVacio(t *testing.T) {
	h := newTestHandler()
	body := `{"nombre":"Carlos","rut":"11111111-1","fecha_nacimiento":"2000-01-01","ciudad":"Arica"}`
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var p Persona
	if err := json.NewDecoder(w.Body).Decode(&p); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if p.Gustos == nil {
		t.Fatal("gustos debe ser [] no null cuando no se envía")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPut, "/personas", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
