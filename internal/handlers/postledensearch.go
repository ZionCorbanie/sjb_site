package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PostLedenSearchHandler struct {
	userStore store.UserStore
}

type PostLedenSearchHandlerParams struct {
	UserStore store.UserStore
}

func NewPostLedenSearchHandler(params PostLedenSearchHandlerParams) *PostLedenSearchHandler {
	return &PostLedenSearchHandler{
		userStore: params.UserStore,
	}
}

func (h *PostLedenSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

	c := templates.RenderUsers(users)
	c.Render(r.Context(), w)
	return
}
