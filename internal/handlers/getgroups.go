package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetGroupsHandler struct {
	GroupStore store.GroupStore
}

type GetGroupsHandlerParams struct {
	GroupStore store.GroupStore
}

func NewGroupsHandler(params GetGroupsHandlerParams) *GetGroupsHandler {
	return &GetGroupsHandler{
		GroupStore: params.GroupStore,
	}
}

func (h *GetGroupsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupType := chi.URLParam(r, "groupType")
	groups, err := h.GroupStore.GetGroupsByType(groupType)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := templates.Groups(groups)
	err = templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
