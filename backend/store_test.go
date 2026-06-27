package main

import "testing"

func TestMemoryStore_AddAndGetAll(t *testing.T) {
	s := NewMemoryStore()
	_ = s.Add(Persona{Nombre: "Juan", RUT: "12345678-9", FechaNacimiento: "1990-01-01", Ciudad: "Santiago", Gustos: []string{"fútbol"}})
	personas, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(personas) != 1 {
		t.Fatalf("expected 1 persona, got %d", len(personas))
	}
	if personas[0].RUT != "12345678-9" {
		t.Fatalf("wrong RUT: %s", personas[0].RUT)
	}
}

func TestMemoryStore_GetAll_Empty(t *testing.T) {
	s := NewMemoryStore()
	personas, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(personas) != 0 {
		t.Fatalf("expected empty, got %d", len(personas))
	}
}

func TestMemoryStore_Delete_Exists(t *testing.T) {
	s := NewMemoryStore()
	_ = s.Add(Persona{Nombre: "Juan", RUT: "21614199-2", FechaNacimiento: "1990-01-01", Ciudad: "Santiago", Gustos: []string{"música"}})
	ok, err := s.Delete("21614199-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Delete to return true")
	}
	personas, _ := s.GetAll()
	if len(personas) != 0 {
		t.Fatal("expected store to be empty after delete")
	}
}

func TestMemoryStore_Delete_NotExists(t *testing.T) {
	s := NewMemoryStore()
	ok, err := s.Delete("99999999-9")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Delete to return false for unknown RUT")
	}
}
