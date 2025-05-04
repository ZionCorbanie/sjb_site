package handlers

import (
	"net/http"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type DeletePromoHandler struct{
    store *dbstore.PromoStore
}

type DeletePromoHandlerParams struct {
    PromoStore *dbstore.PromoStore
}

func NewDeletePromoHandler(params DeletePromoHandlerParams) *DeletePromoHandler {
    return &DeletePromoHandler{
        store: params.PromoStore,
    }
}

func (h *DeletePromoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	promoId := chi.URLParam(r, "promoId")
    
    err := h.store.DeletePromo(promoId)
	if err != nil {
		http.Error(w, "Error deleting promos", http.StatusInternalServerError)
		return
	}

    promos, err := h.store.GetActivePromos()
    err = templates.Promos(promos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
