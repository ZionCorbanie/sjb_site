package handlers

import (
    "sjb_site/internal/templates"
    "net/http"
)

type GetLedenSearchHandler struct {}

func NewLedenSearchHandler() *GetLedenSearchHandler {
	return &GetLedenSearchHandler{}
}

func (h *GetLedenSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.Leden()
	err := templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
