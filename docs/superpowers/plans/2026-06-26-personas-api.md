# Personas API (Go) — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** API REST en Go (stdlib) con almacenamiento en memoria para gestionar personas, con 3 rutas (POST/GET/DELETE) y pruebas unitarias completas.

**Architecture:** Un solo paquete `main` compuesto por 4 archivos: `store.go` (datos en memoria con mutex), `handlers.go` (routing + handlers HTTP), `main.go` (entry point), y sus respectivos test files. Sin frameworks externos, solo stdlib.

**Tech Stack:** Go 1.21+, `net/http`, `net/http/httptest`, `encoding/json`, `sync`

---

## File Map

| Archivo | Responsabilidad |
|---|---|
| `go.mod` | Definición del módulo Go |
| `store.go` | Struct `Persona`, struct `PersonaStore`, métodos Add/GetAll/Delete |
| `store_test.go` | Tests unitarios del store |
| `handlers.go` | Router + handlers HTTP (getPersonas, addPersona, deletePersona) |
| `handlers_test.go` | Tests HTTP con `httptest` |
| `main.go` | Entry point: instancia store + handler, levanta servidor |

---

## Task 1: Inicializar módulo Go

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Inicializar el módulo**

```bash
cd /home/diego/Documentos/GitHub/PruebaDevOps-Diego-Ian
go mod init personas-api
```

Expected: se crea `go.mod` con contenido:
```
module personas-api

go 1.21
```

- [ ] **Step 2: Commit**

```bash
git add go.mod
git commit -m "chore: initialize Go module"
```

---

## Task 2: PersonaStore — TDD

**Files:**
- Create: `store.go`
- Create: `store_test.go`

- [ ] **Step 1: Escribir los tests del store**

Crear `store_test.go`:

```go
package main

import "testing"

func TestStore_AddAndGetAll(t *testing.T) {
	s := NewPersonaStore()
	s.Add(Persona{Nombre: "Juan", RUT: "12345678-9", FechaNacimiento: "1990-01-01", Ciudad: "Santiago"})
	personas := s.GetAll()
	if len(personas) != 1 {
		t.Fatalf("expected 1 persona, got %d", len(personas))
	}
	if personas[0].RUT != "12345678-9" {
		t.Fatalf("wrong RUT: %s", personas[0].RUT)
	}
}

func TestStore_GetAll_Empty(t *testing.T) {
	s := NewPersonaStore()
	personas := s.GetAll()
	if len(personas) != 0 {
		t.Fatalf("expected empty, got %d", len(personas))
	}
}

func TestStore_Delete_Exists(t *testing.T) {
	s := NewPersonaStore()
	s.Add(Persona{Nombre: "Juan", RUT: "12345678-9", FechaNacimiento: "1990-01-01", Ciudad: "Santiago"})
	ok := s.Delete("12345678-9")
	if !ok {
		t.Fatal("expected Delete to return true")
	}
	if len(s.GetAll()) != 0 {
		t.Fatal("expected store to be empty after delete")
	}
}

func TestStore_Delete_NotExists(t *testing.T) {
	s := NewPersonaStore()
	ok := s.Delete("99999999-9")
	if ok {
		t.Fatal("expected Delete to return false for unknown RUT")
	}
}
```

- [ ] **Step 2: Verificar que los tests fallan**

```bash
go test ./... -run TestStore
```

Expected: error de compilación — `Persona`, `NewPersonaStore` no definidos.

- [ ] **Step 3: Implementar store.go**

Crear `store.go`:

```go
package main

import "sync"

type Persona struct {
	Nombre          string `json:"nombre"`
	RUT             string `json:"rut"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Ciudad          string `json:"ciudad"`
}

type PersonaStore struct {
	mu       sync.RWMutex
	personas []Persona
}

func NewPersonaStore() *PersonaStore {
	return &PersonaStore{personas: []Persona{}}
}

func (s *PersonaStore) Add(p Persona) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.personas = append(s.personas, p)
}

func (s *PersonaStore) GetAll() []Persona {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Persona, len(s.personas))
	copy(result, s.personas)
	return result
}

func (s *PersonaStore) Delete(rut string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.personas {
		if p.RUT == rut {
			s.personas = append(s.personas[:i], s.personas[i+1:]...)
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Verificar que los tests pasan**

```bash
go test ./... -run TestStore -v
```

Expected:
```
--- PASS: TestStore_AddAndGetAll
--- PASS: TestStore_GetAll_Empty
--- PASS: TestStore_Delete_Exists
--- PASS: TestStore_Delete_NotExists
PASS
```

- [ ] **Step 5: Commit**

```bash
git add store.go store_test.go
git commit -m "feat: add in-memory PersonaStore with tests"
```

---

## Task 3: Handlers HTTP — TDD

**Files:**
- Create: `handlers.go`
- Create: `handlers_test.go`

- [ ] **Step 1: Escribir los tests de handlers**

Crear `handlers_test.go`:

```go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHandler() *Handler {
	return NewHandler(NewPersonaStore())
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
	body := `{"nombre":"Juan Pérez","rut":"12345678-9","fecha_nacimiento":"1990-05-20","ciudad":"Santiago"}`
	req := httptest.NewRequest(http.MethodPost, "/personas", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
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
	h.store.Add(Persona{Nombre: "Juan Pérez", RUT: "12345678-9", FechaNacimiento: "1990-05-20", Ciudad: "Santiago"})
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
	body := `{"nombre":"Ana López","rut":"98765432-1","fecha_nacimiento":"1995-03-15","ciudad":"Valparaíso"}`
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
```

- [ ] **Step 2: Verificar que los tests fallan**

```bash
go test ./... -run Test
```

Expected: error de compilación — `Handler`, `NewHandler` no definidos.

- [ ] **Step 3: Implementar handlers.go**

Crear `handlers.go`:

```go
package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
```

- [ ] **Step 4: Verificar que todos los tests pasan**

```bash
go test ./... -v
```

Expected: todos en PASS, sin errores de compilación.

- [ ] **Step 5: Commit**

```bash
git add handlers.go handlers_test.go
git commit -m "feat: add HTTP handlers with unit tests"
```

---

## Task 4: Entry point main.go

**Files:**
- Create: `main.go`

- [ ] **Step 1: Crear main.go**

```go
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
```

- [ ] **Step 2: Verificar que compila y corre**

```bash
go build ./...
go run . &
curl -s http://localhost:8080/personas
```

Expected:
```json
[]
```

```bash
kill %1
```

- [ ] **Step 3: Correr todos los tests una última vez**

```bash
go test ./... -v
```

Expected: todos los tests en PASS.

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add entry point and wire up server"
```

---

## Task 5: Verificación del flujo completo POST→GET→DELETE→GET

- [ ] **Step 1: Levantar la app**

```bash
go run . &
```

- [ ] **Step 2: Ejecutar el flujo completo**

```bash
# POST persona 1
curl -s -X POST http://localhost:8080/personas \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Juan Pérez","rut":"12345678-9","fecha_nacimiento":"1990-05-20","ciudad":"Santiago"}'

# POST persona 2
curl -s -X POST http://localhost:8080/personas \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Ana López","rut":"98765432-1","fecha_nacimiento":"1995-03-15","ciudad":"Valparaíso"}'

# GET todas
curl -s http://localhost:8080/personas

# DELETE persona 1
curl -s -X DELETE http://localhost:8080/personas/12345678-9

# GET todas (debe quedar solo Ana)
curl -s http://localhost:8080/personas
```

Expected final GET:
```json
[{"nombre":"Ana López","rut":"98765432-1","fecha_nacimiento":"1995-03-15","ciudad":"Valparaíso"}]
```

- [ ] **Step 3: Apagar la app**

```bash
kill %1
```
