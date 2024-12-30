package handlers

import (
	"sjb_site/internal/templates"
	"net/http"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Admin()
	err := templates.Layout(c, "Sjb admin").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
