# Segunda Evaluación Sumativa — Introducción a DevOps
**Fecha:** 22 de junio de 2026  
**Institución:** UCN Coquimbo — Ingeniería Civil en Computación e Informática / Ingeniería en TI

| | |
|---|---|
| **Puntaje total** | 80 puntos |
| **Nota 4.0** | 48 puntos |
| **Duración** | 3 días (remoto) |
| **Entrega** | 25 de junio de 2026 hasta las 23:59 (Campus Virtual) |
| **Resultados** | 6 de julio de 2026 |

**Modalidad:** Individual o en parejas.  
**Entregable:** Documento PDF subido a Campus Virtual. NO se reciben entregas atrasadas ni por otro medio.

---

## Parte 1: Integración Continua y Despliegue Continuo (30 puntos)

Desarrollar una aplicación backend (NodeJS, Python o PHP recomendados) **sin dependencias externas** que guarde registros de personas **en memoria** con los campos: nombre, RUT, fecha de nacimiento y ciudad donde vive.

| Tarea | Puntos |
|---|---|
| 3 rutas: POST (agregar), DELETE (eliminar), GET (obtener todas) + pruebas unitarias | 6 |
| Subir la app a repositorio público de GitHub con rama principal `main` | 1 |
| Establecer reglas para evitar commits directos a `main` | 3 |
| Crear Azure App Service y configurarlo para desplegar la rama `main` | 5 |
| Flujo de solicitudes: POST → GET → DELETE → GET en Azure (usando el dominio por defecto) | 5 |
| Agregar campo "gustos" por persona (comida, libros, juegos, etc.) aplicando Trunk-Based Development + actualizar tests | 5 |
| Verificar los cambios ejecutando el flujo de solicitudes nuevamente | 5 |

**Entregables:**
- Enlace al repositorio GitHub en el documento PDF.
- Video del proceso completo (se puede omitir el desarrollo de la app base) subido a Drive, YouTube, Vimeo, etc.
- No eliminar los archivos de GitHub Actions ni de App Services del repositorio.
- Eliminar los recursos de Azure al terminar para evitar costos.

---

## Parte 2: Orquestación de Contenedores (20 puntos)

Modificar la aplicación para usar una **base de datos relacional** (PostgreSQL, MySQL o MariaDB) en lugar de memoria.

| Tarea | Puntos |
|---|---|
| `Dockerfile` para construir la imagen de la aplicación | 5 |
| `docker-compose.yaml` — servicio de la base de datos seleccionada | 5 |
| `docker-compose.yaml` — visor de BD: pgAdmin (Postgres) / PhpMyAdmin (MySQL/MariaDB) | 5 |
| `docker-compose.yaml` — servicio de la app construida con el Dockerfile | 2 |
| Limitar memoria y CPU de la BD y del visor (cantidades a elección) | 3 |

**Entregables en el repositorio público:**
- `Dockerfile`
- `docker-compose.yaml`
- `example.env` (variables de entorno)

---

## Parte 3: Cloud DevOps (30 puntos)

Seleccionar **uno** de los tres proveedores: **Google Cloud Platform**, **Amazon Web Services** o **Microsoft Azure**. Investigar sus herramientas y servicios para aplicar DevOps en el ciclo de vida del software (planificación → monitoreo).

| Tarea | Puntos |
|---|---|
| Breve descripción de cada herramienta/servicio y la etapa del ciclo de vida donde se aplica | 5 |
| Diagrama de arquitectura simple para montar una aplicación usando esas herramientas (puede usarse la app de partes anteriores como ejemplo) | 10 |
| Cuadro comparativo (similitudes y diferencias) entre las herramientas del proveedor y las vistas en clases | 15 |

**Herramientas del curso a comparar:** GitHub, Jenkins, SonarQube, Docker, Kubernetes, Ansible, Terraform, ELK Stack, Prometheus, Grafana.

---

## Resultados de Aprendizaje evaluados

- **RA4:** Experiencia práctica con herramientas y tecnologías populares de DevOps.
- **RA5:** Automatización de flujos CI/CD.
- **RA7:** Mejores prácticas de aprovisionamiento de infraestructura y gestión de configuraciones.
- **RA8:** Tecnologías de contenedores y su papel en DevOps.
- **RA9:** Infraestructura como código (IaC) y sus beneficios.
- **RA10:** Consideraciones de seguridad y cumplimiento en DevOps.
