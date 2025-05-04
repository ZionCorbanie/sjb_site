package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type PromoEditHandler struct{
    promoStore store.PromoStore
}

type PromoEditHandlerParams struct {
    PromoStore store.PromoStore
}

func NewPromoEditHandler(params PromoEditHandlerParams) *PromoEditHandler {
	return &PromoEditHandler{
        promoStore: params.PromoStore,
    }
}

func (h *PromoEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promoId := chi.URLParam(r, "promoId")

	promo, err := h.promoStore.GetPromo(promoId)
	if err != nil {
		http.Error(w, "Error getting promos", http.StatusInternalServerError)
		return
	}

	err = templates.EditPromo(promo).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
