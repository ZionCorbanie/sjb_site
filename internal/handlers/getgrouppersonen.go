package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetGroupPersonenHandler struct {
	GroupStore     store.GroupStore
	GroupUserStore store.GroupUserStore
	UserStore      store.UserStore
}

type GetGroupPersonenHandlerParams struct {
	GroupStore     store.GroupStore
	GroupUserStore store.GroupUserStore
	UserStore      store.UserStore
}

func NewGroupPersonenHandler(params GetGroupPersonenHandlerParams) *GetGroupPersonenHandler {
	return &GetGroupPersonenHandler{
		GroupStore:     params.GroupStore,
		GroupUserStore: params.GroupUserStore,
		UserStore:      params.UserStore,
	}
}

func (h *GetGroupPersonenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupType := chi.URLParam(r, "groupType")
	groupId := chi.URLParam(r, "groupId")
	groupUsers, err := h.GroupUserStore.GetGroupUsersByGroup(groupId)
	if err != nil {
		http.Error(w, "Error getting group users", http.StatusInternalServerError)
		return
	}

	group, err := h.GroupStore.GetGroup(groupId)
	if err != nil {
		http.Error(w, "Error getting group", http.StatusInternalServerError)
		return
	}

	t := templates.UserTable(groupUsers, group)
	c := templates.GroupPersonen(t, groupType, groupId)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
