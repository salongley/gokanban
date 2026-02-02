package main

import (
	"gokanban/internal/adapters"
	"log"
	"net/http"
)

// ============================================================================
// 4. MAIN WIRING
// ============================================================================

func main() {
	renderer := adapters.NewNegotiatingRenderer("./templates/*.html")
	service := adapters.NewInMemoryService()
	// Seed data for all three columns
	service.Create("Research HTMX", "todo")
	service.Create("Write Go Handlers", "doing")
	service.Create("Setup Project", "done")

	handler := adapters.NewCardHandler(service, renderer)
	router := http.NewServeMux()

	router.HandleFunc("GET /", handler.HandleIndex)
	router.HandleFunc("POST /cards", handler.HandleCreate)
	router.HandleFunc("GET /cards/{id}/edit", handler.HandleEditForm)
	router.HandleFunc("PATCH /cards/{id}/move", handler.HandleMove)

	log.Println("Starting GoKanban on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
