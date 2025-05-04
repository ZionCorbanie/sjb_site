package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PrikbordCreateHandler struct{
    promoStore store.PromoStore
}

type PrikbordCreateHandlerParams struct {
    PromoStore store.PromoStore
}

func NewPrikbordCreateHandler(params PrikbordCreateHandlerParams) *PrikbordCreateHandler {
	return &PrikbordCreateHandler{
        promoStore: params.PromoStore,
    }
}

func (h *PrikbordCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promos, err := h.promoStore.GetAllPromos()
	if err != nil {
		http.Error(w, "Error getting promos", http.StatusInternalServerError)
		return
	}

	err = templates.Layout(templates.Card(templates.Promos(promos)), "Promo maken").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
