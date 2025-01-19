package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetCreatePostHandler struct{}

func NewGetCreatePostHandler() *GetCreatePostHandler {
	return &GetCreatePostHandler{}
}

func (h *GetCreatePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c := templates.CreatePost()
	err := templates.Layout(c, "Post toevoegen").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
