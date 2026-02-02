package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// ============================================================================
// 1. DOMAIN LAYER (The Core Logic)
// ============================================================================

type Card struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Column   string `json:"column"` // 'todo', 'doing', 'done'
	EditedAt string `json:"edited_at"`
}

// Service Interface (Port)
type CardService interface {
	List() []Card
	Create(title, column string) Card
	Update(id int, title string) Card
	Find(id int) Card
}

// ============================================================================
// 2. ADAPTERS LAYER (Infrastructure)
// ============================================================================

// --- A. The Negotiating Renderer (Hexagonal Output Port) ---
type RenderContext struct {
	TemplateName string
	Data         any
}

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, ctx RenderContext)
}

type NegotiatingRenderer struct {
	tmpl *template.Template
}

func NewNegotiatingRenderer(pattern string) *NegotiatingRenderer {
	return &NegotiatingRenderer{
		tmpl: template.Must(template.ParseGlob(pattern)),
	}
}

func (n *NegotiatingRenderer) Render(w http.ResponseWriter, r *http.Request, ctx RenderContext) {
	// 1. Check Accept Header for JSON
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ctx.Data)
		return
	}

	// 2. Default to HTML (HTMX)
	w.Header().Set("Content-Type", "text/html")
	if ctx.TemplateName == "" { return }

	err := n.tmpl.ExecuteTemplate(w, ctx.TemplateName, ctx.Data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// --- B. In-Memory Repository (Database Adapter) ---
type InMemoryService struct {
	sync.Mutex
	cards  []Card
	nextID int
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		cards:  []Card{},
		nextID: 1,
	}
}

func (s *InMemoryService) List() []Card {
	return s.cards
}

func (s *InMemoryService) Create(title, column string) Card {
	s.Lock()
	defer s.Unlock()
	card := Card{
		ID:       s.nextID,
		Title:    title,
		Column:   column,
		EditedAt: time.Now().Format(time.Kitchen),
	}
	s.cards = append(s.cards, card)
	s.nextID++
	return card
}

func (s *InMemoryService) Update(id int, title string) Card { return Card{ID: id, Title: title} }
func (s *InMemoryService) Find(id int) Card                 { return Card{ID: id, Title: "Found"} }

// ============================================================================
// 3. HTTP HANDLERS (The Web Layer)
// ============================================================================

type CardHandler struct {
	service  CardService
	renderer Renderer
}

func NewCardHandler(s CardService, r Renderer) *CardHandler {
	return &CardHandler{service: s, renderer: r}
}

// GET / - Show the Board
func (h *CardHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	cards := h.service.List()
	h.renderer.Render(w, r, RenderContext{
		TemplateName: "index.html",
		Data:         map[string]any{"Cards": cards},
	})
}

// POST /cards - Create a Card
func (h *CardHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	newCard := h.service.Create(r.FormValue("title"), r.FormValue("column"))

	// Render ONLY the card partial
	h.renderer.Render(w, r, RenderContext{
		TemplateName: "card.html",
		Data:         newCard,
	})
}

// GET /cards/{id}/edit - Return the Form
func (h *CardHandler) HandleEditForm(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	card := h.service.Find(id)

	h.renderer.Render(w, r, RenderContext{
		TemplateName: "card_edit.html",
		Data:         card,
	})
}

// ============================================================================
// 4. MAIN WIRING
// ============================================================================

func main() {
	renderer := NewNegotiatingRenderer("./templates/*.html")
	service := NewInMemoryService()
	service.Create("Learn Go", "todo")

	handler := NewCardHandler(service, renderer)
	router := http.NewServeMux()

	router.HandleFunc("GET /", handler.HandleIndex)
	router.HandleFunc("POST /cards", handler.HandleCreate)
	router.HandleFunc("GET /cards/{id}/edit", handler.HandleEditForm)

	log.Println("Starting GoKanban on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
