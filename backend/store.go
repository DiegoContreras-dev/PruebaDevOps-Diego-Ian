package main

import "sync"

type Persona struct {
	Nombre          string   `json:"nombre"`
	RUT             string   `json:"rut"`
	FechaNacimiento string   `json:"fecha_nacimiento"`
	Ciudad          string   `json:"ciudad"`
	Gustos          []string `json:"gustos"`
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
