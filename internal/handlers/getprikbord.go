package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PrikbordHandler struct{
    prikbordStore store.PostStore
}

type PrikbordHandlerParams struct {
    PrikbordStore store.PostStore
}

func NewPrikbordHandler(params PrikbordHandlerParams) *PrikbordHandler {
	return &PrikbordHandler{
        prikbordStore: params.PrikbordStore,
    }
}

func (h *PrikbordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := templates.Prikbord().Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
