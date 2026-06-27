# Segunda Evaluación Sumativa — Introducción a DevOps

**Curso:** Introducción a DevOps — UCN Coquimbo  
**Integrantes:** Diego Contreras, Ian

---

## Parte 1: Integración Continua y Despliegue Continuo

Aplicación backend en **Go** con almacenamiento en memoria que gestiona registros de personas.

### Rutas disponibles

| Método | Ruta | Descripción |
|---|---|---|
| `GET` | `/personas` | Retorna todas las personas |
| `POST` | `/personas` | Agrega una persona |
| `DELETE` | `/personas/{rut}` | Elimina una persona por RUT |

### Modelo de datos

```json
{
  "nombre": "Juan Pérez",
  "rut": "12345678-9",
  "fecha_nacimiento": "1990-05-20",
  "ciudad": "Santiago"
}
```

### Ejecutar localmente

```bash
cd backend
go run .
```

El servidor levanta en `http://localhost:8080` por defecto. Se puede cambiar el puerto con la variable de entorno `PORT`.

### Ejecutar tests

```bash
cd backend
go test ./... -v
```

### Despliegue

La aplicación está desplegada en **Azure App Service**.  
URL: *(pendiente)*

---

## Parte 2: Orquestación de Contenedores

*(pendiente)*

---

## Parte 3: Cloud DevOps

*(pendiente)*
