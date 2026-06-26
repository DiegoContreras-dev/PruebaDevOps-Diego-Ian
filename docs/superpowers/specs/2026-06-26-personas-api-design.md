# Diseño: API de Personas en Go — Parte 1 CI/CD

**Fecha:** 2026-06-26  
**Contexto:** Segunda Evaluación Sumativa — Introducción a DevOps

---

## Objetivo

App backend en Go con almacenamiento **en memoria** (sin dependencias externas) que gestiona registros de personas. Servirá como base para las Partes 2 y 3 del examen.

---

## Estructura de archivos

```
/
├── main.go           # entry point, registra rutas y levanta el server
├── store.go          # almacén en memoria (slice + sync.RWMutex)
├── handlers.go       # handlers HTTP
├── handlers_test.go  # pruebas unitarias con net/http/httptest
└── go.mod
```

---

## Modelo de datos

```go
type Persona struct {
    Nombre          string `json:"nombre"`
    RUT             string `json:"rut"`
    FechaNacimiento string `json:"fecha_nacimiento"` // formato: YYYY-MM-DD
    Ciudad          string `json:"ciudad"`
}
```

---

## Rutas

| Método   | Ruta               | Descripción                          | Éxito | Error                    |
|----------|--------------------|--------------------------------------|-------|--------------------------|
| `POST`   | `/personas`        | Agrega una persona (body JSON)       | 201   | 400 (body inválido), 500 |
| `GET`    | `/personas`        | Retorna todas las personas como JSON | 200   | 500                      |
| `DELETE` | `/personas/{rut}`  | Elimina la persona con ese RUT       | 200   | 404 (no existe)          |

### Códigos de respuesta usados

| Código | Significado      | Cuándo                                   |
|--------|------------------|------------------------------------------|
| 200    | OK               | GET exitoso, DELETE exitoso              |
| 201    | Created          | POST exitoso                             |
| 400    | Bad Request      | JSON inválido o campos vacíos            |
| 404    | Not Found        | DELETE con RUT que no existe             |
| 405    | Method Not Allowed | Método HTTP no permitido en la ruta   |
| 500    | Internal Server Error | Error al codificar/decodificar JSON |

---

## Diseño del Store

```go
type PersonaStore struct {
    mu       sync.RWMutex
    personas []Persona
}
```

- `Add(p Persona) error` — agrega persona
- `GetAll() []Persona` — retorna copia del slice
- `Delete(rut string) bool` — elimina por RUT, retorna si existía

El mutex garantiza seguridad ante requests concurrentes.

---

## Testing

Usando `net/http/httptest` — sin levantar servidor real:

- `TestGetPersonas_Empty` — GET retorna array vacío
- `TestAddPersona_OK` — POST válido retorna 201
- `TestAddPersona_InvalidBody` — POST con JSON roto retorna 400
- `TestDeletePersona_OK` — DELETE persona existente retorna 200
- `TestDeletePersona_NotFound` — DELETE RUT inexistente retorna 404
- `TestGetPersonas_AfterAdd` — GET retorna persona recién agregada

---

## Consideraciones

- El identificador de eliminación es el **RUT** (clave natural).
- El servidor escucha en el puerto definido por la variable de entorno `PORT` (default `8080`) para compatibilidad con Azure App Service.
- Todos los responses son `Content-Type: application/json`.
