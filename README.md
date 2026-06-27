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
  "ciudad": "Santiago",
  "gustos": ["fútbol", "pizza", "películas"]
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

La aplicación está desplegada en **Azure App Service** con CI/CD via GitHub Actions.  
URL: https://pruebadevops-chahgwgxawgafnf4.canadacentral-01.azurewebsites.net  
Swagger UI: https://pruebadevops-chahgwgxawgafnf4.canadacentral-01.azurewebsites.net/swagger/index.html

### Link Video

https://drive.google.com/drive/folders/1gKgYxFyY-gt3xpPZ03l_YN3s26rj2Lk4?usp=drive_link

---

## Parte 2: Orquestación de Contenedores

La aplicación fue migrada para usar **PostgreSQL** como base de datos persistente, orquestada con Docker Compose junto a **pgAdmin** como visor de base de datos.

### Servicios

| Servicio | Imagen | Puerto | Descripción |
|---|---|---|---|
| `app` | Build local (`./backend`) | `8080` | API Go conectada a PostgreSQL |
| `postgres` | `postgres:16-alpine` | `5432` (interno) | Base de datos relacional |
| `pgadmin` | `dpage/pgadmin4:8` | `5050` | Visor web de la base de datos |

### Límites de recursos

| Servicio | CPU | Memoria |
|---|---|---|
| `app` | 0.50 | 128 MB |
| `postgres` | 0.50 | 256 MB |
| `pgadmin` | 0.25 | 256 MB |

### Variables de entorno

Copiar `example.env` como `.env` antes de levantar:

```bash
cp example.env .env
docker compose up --build -d
```

### Acceso local

| Servicio | URL | Credenciales |
|---|---|---|
| API | `http://localhost:8080/personas` | — |
| Swagger UI | `http://localhost:8080/swagger/index.html` | — |
| pgAdmin | `http://localhost:5050` | `admin@admin.com` / `admin` |

### Conectar pgAdmin a PostgreSQL

En pgAdmin → **Add New Server**:

| Campo | Valor |
|---|---|
| Host | `postgres` |
| Port | `5432` |
| Database | `personas_db` |
| Username | `personas` |
| Password | `secret` |

---

## Parte 3: Cloud DevOps

*(pendiente)*
