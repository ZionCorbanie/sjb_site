package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetUserManagementHandler struct{}

func NewGetUserManagementHandler() *GetUserManagementHandler {
	return &GetUserManagementHandler{}
}

func (h *GetUserManagementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := templates.UserManagement()
	err := templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
