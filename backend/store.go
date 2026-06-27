package main

import (
	"database/sql"
	"sync"

	"github.com/lib/pq"
)

// Persona representa una persona en el sistema.
type Persona struct {
	Nombre          string   `json:"nombre"`
	RUT             string   `json:"rut"`
	FechaNacimiento string   `json:"fecha_nacimiento"`
	Ciudad          string   `json:"ciudad"`
	Gustos          []string `json:"gustos"`
}

// PersonaRepository define las operaciones sobre personas.
// Permite intercambiar MemoryStore (tests) por PostgresStore (producción).
type PersonaRepository interface {
	Add(p Persona) error
	GetAll() ([]Persona, error)
	Delete(rut string) (bool, error)
}

// ─────────────────────────────────────────────────────────────
// MemoryStore — implementación en memoria para tests unitarios
// ─────────────────────────────────────────────────────────────

type MemoryStore struct {
	mu       sync.RWMutex
	personas []Persona
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{personas: []Persona{}}
}

func (s *MemoryStore) Add(p Persona) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p.Gustos == nil {
		p.Gustos = []string{}
	}
	s.personas = append(s.personas, p)
	return nil
}

func (s *MemoryStore) GetAll() ([]Persona, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Persona, len(s.personas))
	copy(result, s.personas)
	return result, nil
}

func (s *MemoryStore) Delete(rut string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, p := range s.personas {
		if p.RUT == rut {
			s.personas = append(s.personas[:i], s.personas[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

// ─────────────────────────────────────────────────────────────
// PostgresStore — implementación con PostgreSQL (producción)
// ─────────────────────────────────────────────────────────────

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// InitDB crea la tabla personas si no existe.
func InitDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS personas (
			rut              TEXT PRIMARY KEY,
			nombre           TEXT NOT NULL,
			fecha_nacimiento TEXT NOT NULL,
			ciudad           TEXT NOT NULL,
			gustos           TEXT[] NOT NULL DEFAULT '{}'
		)
	`)
	return err
}

// Add inserta una persona. Ignora duplicados (idempotente en seed).
func (s *PostgresStore) Add(p Persona) error {
	if p.Gustos == nil {
		p.Gustos = []string{}
	}
	_, err := s.db.Exec(`
		INSERT INTO personas (rut, nombre, fecha_nacimiento, ciudad, gustos)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (rut) DO NOTHING`,
		p.RUT, p.Nombre, p.FechaNacimiento, p.Ciudad, pq.Array(p.Gustos),
	)
	return err
}

// GetAll retorna todas las personas ordenadas por nombre.
func (s *PostgresStore) GetAll() ([]Persona, error) {
	rows, err := s.db.Query(`
		SELECT rut, nombre, fecha_nacimiento, ciudad, gustos
		FROM personas
		ORDER BY nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	personas := []Persona{}
	for rows.Next() {
		var p Persona
		var gustos pq.StringArray
		if err := rows.Scan(&p.RUT, &p.Nombre, &p.FechaNacimiento, &p.Ciudad, &gustos); err != nil {
			return nil, err
		}
		p.Gustos = []string(gustos)
		if p.Gustos == nil {
			p.Gustos = []string{}
		}
		personas = append(personas, p)
	}
	return personas, rows.Err()
}

// Delete elimina una persona por RUT. Retorna false si no existe.
func (s *PostgresStore) Delete(rut string) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM personas WHERE rut = $1`, rut)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
