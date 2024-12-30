package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/hash"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PostLedenSearchHandler struct {
	userStore         store.UserStore
	sessionStore      store.SessionStore
	passwordhash      hash.PasswordHash
	sessionCookieName string
}

type PostLedenSearchHandlerParams struct {
	UserStore         store.UserStore
	SessionStore      store.SessionStore
	PasswordHash      hash.PasswordHash
	SessionCookieName string
}

func NewPostLedenSearchHandler(params PostLedenSearchHandlerParams) *PostLedenSearchHandler {
	return &PostLedenSearchHandler{
		userStore:         params.UserStore,
		sessionStore:      params.SessionStore,
		passwordhash:      params.PasswordHash,
		sessionCookieName: params.SessionCookieName,
	}
}

func (h *PostLedenSearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

    fmt.Println(users)

    c := templates.RenderLeden(users)
    c.Render(r.Context(), w)
    return
}
