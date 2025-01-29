package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type PostUserSearchGroupHandler struct {
	userStore store.UserStore
}

type PostUserSearchGroupHandlerParams struct {
	UserStore store.UserStore
}

func NewPostUserSearcGrouphHandler(params PostUserSearchGroupHandlerParams) *PostUserSearchGroupHandler {
	return &PostUserSearchGroupHandler{
		userStore: params.UserStore,
	}
}

func (h *PostUserSearchGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	groupType := chi.URLParam(r, "groupType")
	search := r.FormValue("search")
	if search == "" {
		return
	}
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

	c := templates.RenderUsersGroup(users, groupId, groupType)
	c.Render(r.Context(), w)
}
