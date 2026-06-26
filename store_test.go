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
