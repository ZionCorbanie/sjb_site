package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type PromoHandler struct{
    promoStore store.PromoStore
}

type PromoHandlerParams struct {
    PromoStore store.PromoStore
}

func NewPromoHandler(params PromoHandlerParams) *PromoHandler {
	return &PromoHandler{
        promoStore: params.PromoStore,
    }
}

func (h *PromoHandler) Show(w http.ResponseWriter, r *http.Request) {
	promos, err := h.promoStore.GetActivePromos()
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

func (h *PromoHandler) Popup(w http.ResponseWriter, r *http.Request) {
	promoId := chi.URLParam(r, "promoId")
	promo, err := h.promoStore.GetPromo(promoId)
	if err != nil{
		http.Error(w, "Promo not found", http.StatusNotFound)
		return
	}

	err = templates.Popup(templates.PromoPopup(promo)).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
