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
	service.Create("Learn Go", "todo")

	handler := adapters.NewCardHandler(service, renderer)
	router := http.NewServeMux()

	router.HandleFunc("GET /", handler.HandleIndex)
	router.HandleFunc("POST /cards", handler.HandleCreate)
	router.HandleFunc("GET /cards/{id}/edit", handler.HandleEditForm)

	log.Println("Starting GoKanban on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
