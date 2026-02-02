package adapters

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

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
	if ctx.TemplateName == "" {
		return
	}

	err := n.tmpl.ExecuteTemplate(w, ctx.TemplateName, ctx.Data)
	if err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
