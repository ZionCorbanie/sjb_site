package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetCreateGroupHandler struct{}

func NewGetCreateGroupHandler() *GetCreateGroupHandler {
	return &GetCreateGroupHandler{}
}

func (h *GetCreateGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.CreateGroup()
	err := templates.Layout(c, "Groep toevoegen").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
