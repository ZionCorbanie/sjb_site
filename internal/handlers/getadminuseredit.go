package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetAdminUserEditHandler struct {
	userStore store.UserStore
}

type GetAdminUserEditHandlerParams struct {
	UserStore store.UserStore
}

func NewAdminUserEditHandler(params GetAdminUserEditHandlerParams) *GetAdminUserEditHandler {
	return &GetAdminUserEditHandler{
		userStore: params.UserStore,
	}
}

func (h *GetAdminUserEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	c := templates.UserEditAdmin(user)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
