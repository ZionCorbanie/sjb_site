package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetCreateMenuHandler struct{}

func NewGetCreateMenuHandler() *GetCreateMenuHandler {
	return &GetCreateMenuHandler{}
}

func (h *GetCreateMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c := templates.CreateMenu()
	err := templates.Layout(c, "Menu toevoegen").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
