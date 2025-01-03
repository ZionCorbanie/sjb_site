package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetUserEditHandler struct {
	userStore store.UserStore
}

type GetUserEditHandlerParams struct {
	UserStore store.UserStore
}

func NewUserEditHandler(params GetUserEditHandlerParams) *GetUserEditHandler {
	return &GetUserEditHandler{
		userStore: params.UserStore,
	}
}

func (h *GetUserEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	user, err := h.userStore.GetUserById(userId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := templates.UserEdit(user)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
