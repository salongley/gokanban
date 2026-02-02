package adapters

import (
	"gokanban/internal/service"
	"net/http"
	"strconv"
)

// ============================================================================
// 3. HTTP HANDLERS (The Web Layer)
// ============================================================================

type CardHandler struct {
	service  service.CardService
	renderer Renderer
}

func NewCardHandler(s service.CardService, r Renderer) *CardHandler {
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

// PATCH /cards/{id}/move - Move card to a new column
func (h *CardHandler) HandleMove(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	newColumn := r.FormValue("column") // 'todo', 'doing', or 'done'

	// Update the data in memory

	card := h.service.Find(id)

	card.Column = newColumn
	h.service.Update(id, card)
	// We don't need to render HTML back because SortableJS has
	// already visually moved the DOM element. We just return 200 OK.
	w.WriteHeader(http.StatusOK)
}
