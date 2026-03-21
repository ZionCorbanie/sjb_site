package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetJaarclubsHandler struct {
	GroupStore store.GroupStore
}

type GetJaarclubsHandlerParams struct {
	GroupStore store.GroupStore
}

func NewJaarclubsHandler(params GetJaarclubsHandlerParams) *GetJaarclubsHandler {
	return &GetJaarclubsHandler{
		GroupStore: params.GroupStore,
	}
}

func (h *GetJaarclubsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    jaar, err := strconv.ParseInt(chi.URLParam(r, "jaarlaag"), 10, 0)
    if err != nil {
        err = templates.Layout(templates.Card(templates.Jaarclubs()), "Sint Jansbrug - Jaarclubs").Render(r.Context(), w)
        if err != nil {
            http.Error(w, "Error rendering template", http.StatusInternalServerError)
            return
        }
        return
    }

	groups, err := h.GroupStore.GetJaarclubs(int(jaar))
	if err != nil || len(*groups) == 0 {
        w.WriteHeader(http.StatusOK)
		return
	}

	err = templates.JaarclubRow(groups, int(jaar)).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
