package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
	"strconv"
)

type GetInfoHandler struct{}

func NewInfoHandler() *GetInfoHandler {
	return &GetInfoHandler{}
}

func (h *GetInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	index, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	content := templates.InfoContent(item)

	err = templates.Info(index, content).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
