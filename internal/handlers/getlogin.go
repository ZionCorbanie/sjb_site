package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetLoginHandler struct{}

func NewGetLoginHandler() *GetLoginHandler {
	return &GetLoginHandler{}
}

func (h *GetLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    target := r.URL.Query().Get("redirect")

    if target == "" {
        target = "/"
    }

	c := templates.Login("Login", target)
	err := templates.Layout(c, "My website").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
