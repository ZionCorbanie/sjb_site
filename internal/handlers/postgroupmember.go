package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostGroupMemberHandler struct {
	groupUserStore store.GroupUserStore
	userStore      store.UserStore
	groupStore     store.GroupStore
}

type PostGroupMemberHandlerParams struct {
	GroupUserStore store.GroupUserStore
	UserStore      store.UserStore
	GroupStore     store.GroupStore
}

func NewPostGroupMemberHandler(params PostGroupMemberHandlerParams) *PostGroupMemberHandler {
	return &PostGroupMemberHandler{
		groupUserStore: params.GroupUserStore,
		userStore:      params.UserStore,
		groupStore:     params.GroupStore,
	}
}

func (h *PostGroupMemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupId, _ := strconv.ParseUint(chi.URLParam(r, "groupId"), 10, 64)
	userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	err := h.groupUserStore.AddUserToGroup(uint(userId), uint(groupId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupUsers, err := h.groupUserStore.GetGroupUsersByGroup(chi.URLParam(r, "groupId"))
	if err != nil {
		http.Error(w, "Error getting group users", http.StatusInternalServerError)
		return
	}
	group, err := h.groupStore.GetGroup(chi.URLParam(r, "groupId"))
	if err != nil {
		http.Error(w, "Error getting group", http.StatusInternalServerError)
		return
	}

	t := templates.UserTable(groupUsers, group)
	err = t.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
