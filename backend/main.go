package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	docs "personas-api/docs"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Configurar host de Swagger según entorno (Azure inyecta WEBSITE_HOSTNAME)
	if azureHost := os.Getenv("WEBSITE_HOSTNAME"); azureHost != "" {
		docs.SwaggerInfo.Host = azureHost
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	// ── Conexión a PostgreSQL ─────────────────────────────────
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL no configurada")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Error abriendo conexión a PostgreSQL: %v", err)
	}
	defer db.Close()

	// Retry loop: PostgreSQL tarda unos segundos en estar listo en Docker
	for i := 1; i <= 15; i++ {
		if err := db.Ping(); err == nil {
			log.Println("✅ Conectado a PostgreSQL")
			break
		}
		log.Printf("⏳ Esperando PostgreSQL... intento %d/15", i)
		time.Sleep(2 * time.Second)
		if i == 15 {
			log.Fatal("❌ No se pudo conectar a PostgreSQL tras 15 intentos")
		}
	}

	// ── Inicializar esquema ───────────────────────────────────
	if err := InitDB(db); err != nil {
		log.Fatalf("❌ Error inicializando esquema: %v", err)
	}

	// ── Seed inicial (ON CONFLICT DO NOTHING → idempotente) ──
	store := NewPostgresStore(db)
	seeds := []Persona{
		{Nombre: "Juan Pérez", RUT: "21614199-2", FechaNacimiento: "1990-05-20", Ciudad: "Santiago", Gustos: []string{"fútbol", "pizza", "películas"}},
		{Nombre: "Ana López", RUT: "18404852-4", FechaNacimiento: "1995-03-15", Ciudad: "Coquimbo", Gustos: []string{"lectura", "sushi", "videojuegos"}},
		{Nombre: "Diego Contreras", RUT: "12824745-k", FechaNacimiento: "2000-01-01", Ciudad: "La Serena", Gustos: []string{"ciberseguridad", "golang", "música"}},
	}
	for _, p := range seeds {
		if err := store.Add(p); err != nil {
			log.Printf("⚠️  Seed [%s]: %v", p.RUT, err)
		}
	}
	log.Println("✅ Seed completado")

	// ── Rutas HTTP ────────────────────────────────────────────
	handler := NewHandler(store)
	mux := http.NewServeMux()
	mux.Handle("/personas", handler)
	mux.Handle("/personas/", handler)
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor escuchando en :%s", port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
