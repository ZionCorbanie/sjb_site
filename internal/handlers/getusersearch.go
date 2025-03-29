package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetUserSearchHandler struct{}

func NewUserSearchHandler() *GetUserSearchHandler {
	return &GetUserSearchHandler{}
}

func (h *GetUserSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Users()
	err := templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
