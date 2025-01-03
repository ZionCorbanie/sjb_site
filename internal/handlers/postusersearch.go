package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PostUserSearchHandler struct {
	userStore store.UserStore
}

type PostUserSearchHandlerParams struct {
	UserStore store.UserStore
}

func NewPostUserSearchHandler(params PostUserSearchHandlerParams) *PostUserSearchHandler {
	return &PostUserSearchHandler{
		userStore: params.UserStore,
	}
}

func (h *PostUserSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

	c := templates.RenderUsers(users)
	c.Render(r.Context(), w)
	return
}
