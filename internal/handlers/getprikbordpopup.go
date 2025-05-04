package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type PrikbordPopupHandler struct{
    prikbordPopupStore store.PromoStore
}

type PrikbordPopupHandlerParams struct {
    PrikbordPopupStore store.PromoStore
}

func NewPrikbordPopupHandler(params PrikbordPopupHandlerParams) *PrikbordPopupHandler {
	return &PrikbordPopupHandler{
        prikbordPopupStore: params.PrikbordPopupStore,
    }
}

func (h *PrikbordPopupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promoId := chi.URLParam(r, "promoId")
	promo, err := h.prikbordPopupStore.GetPromo(promoId)
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
