package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DeleteGroupMemberHandler struct {
	groupUserStore store.GroupUserStore
	groupStore     store.GroupStore
}

type DeleteGroupMemberHandlerParams struct {
	GroupUserStore store.GroupUserStore
	GroupStore     store.GroupStore
}

func NewDeleteGroupMemberHandler(params DeleteGroupMemberHandlerParams) *DeleteGroupMemberHandler {
	return &DeleteGroupMemberHandler{
		groupUserStore: params.GroupUserStore,
		groupStore:     params.GroupStore,
	}
}

func (h *DeleteGroupMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	group, err := h.groupStore.GetGroup(chi.URLParam(r, "groupId"))
	if err != nil {
		http.Error(w, "Error getting group", http.StatusInternalServerError)
		return
	}
	groupId, _ := strconv.ParseUint(chi.URLParam(r, "groupId"), 10, 64)
	userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	err = h.groupUserStore.DeleteGroupUser(uint(userId), uint(groupId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupUsers, err := h.groupUserStore.GetGroupUsersByGroup(chi.URLParam(r, "groupId"))
	if err != nil {
		http.Error(w, "Error getting group users", http.StatusInternalServerError)
		return
	}

	t := templates.UserTable(groupUsers, group)
	err = t.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
