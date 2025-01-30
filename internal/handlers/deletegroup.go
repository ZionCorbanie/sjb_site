package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"

	"github.com/go-chi/chi/v5"
)

type DeleteGroupHandler struct {
	store store.GroupStore
}

type DeleteGroupHandlerParams struct {
	GroupStore store.GroupStore
}

func NewDeleteGroupHandler(params DeleteGroupHandlerParams) *DeleteGroupHandler {
	return &DeleteGroupHandler{
		store: params.GroupStore,
	}
}

func (h *DeleteGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupId")

	group, err := h.store.GetGroup(groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !middleware.IsAdmin(r.Context()) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = h.store.DeleteGroup(groupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/admin/groepen/%s", group.GroupType))
}
