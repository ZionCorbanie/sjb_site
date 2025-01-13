package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PostUserManagementHandler struct {
	userStore store.UserStore
}

type PostUserManagementHandlerParams struct {
	UserStore store.UserStore
}

func NewPostUserManagementHandler(params PostUserManagementHandlerParams) *PostUserManagementHandler {
	return &PostUserManagementHandler{
		userStore: params.UserStore,
	}
}

func (h *PostUserManagementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

	c := templates.RenderUserManagement(users)
	c.Render(r.Context(), w)
	return
}
