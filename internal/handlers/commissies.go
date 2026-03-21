package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type GetCommissiesHandler struct {
	GroupStore store.GroupStore
}

type GetCommissiesHandlerParams struct {
	GroupStore store.GroupStore
}

func NewCommissiesHandler(params GetCommissiesHandlerParams) *GetCommissiesHandler {
	return &GetCommissiesHandler{
		GroupStore: params.GroupStore,
	}
}

func (h *GetCommissiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groups, err := h.GroupStore.GetCommissies()
	if err != nil || len(*groups) == 0 {
        w.WriteHeader(http.StatusOK)
		return
	}

	err = templates.Layout(templates.Card(templates.Groups(groups)), "Commissies").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
