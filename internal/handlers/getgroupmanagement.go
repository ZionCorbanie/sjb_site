package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetGroupManagementHandler struct {
	GroupStore store.GroupStore
}

type GetGroupManagementHandlerParams struct {
	GroupStore store.GroupStore
}

func NewGroupManagementHandler(params GetGroupManagementHandlerParams) *GetGroupManagementHandler {
	return &GetGroupManagementHandler{
		GroupStore: params.GroupStore,
	}
}

func (h *GetGroupManagementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	s := templates.GroupManagementSidebar()
	c := templates.GroupManagement(groups)
	err = templates.Layout(templates.Sidebar(c, s), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
