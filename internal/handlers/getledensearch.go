package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetLedenSearchHandler struct{}

func NewLedenSearchHandler() *GetLedenSearchHandler {
	return &GetLedenSearchHandler{}
}

func (h *GetLedenSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Users()
	err := templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
