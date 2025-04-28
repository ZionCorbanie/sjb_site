package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PrikbordHandler struct{
    prikbordStore store.PromoStore
}

type PrikbordHandlerParams struct {
    PrikbordStore store.PromoStore
}

func NewPrikbordHandler(params PrikbordHandlerParams) *PrikbordHandler {
	return &PrikbordHandler{
        prikbordStore: params.PrikbordStore,
    }
}

func (h *PrikbordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promos, err := h.prikbordStore.GetActivePromos()
	if err != nil {
		http.Error(w, "Error getting promos", http.StatusInternalServerError)
		return
	}

	err = templates.Prikbord(promos).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
