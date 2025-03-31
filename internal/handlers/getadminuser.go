package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetAdminUserHandler struct {
	userStore store.UserStore
}

type GetAdminUserHandlerParams struct {
	UserStore store.UserStore
}

func NewAdminUserHandler(params GetAdminUserHandlerParams) *GetAdminUserHandler {
	return &GetAdminUserHandler{
		userStore: params.UserStore,
	}
}

func (h *GetAdminUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	s := templates.AdminSidebarUser(user)
	c := templates.UserEditAdmin(user)
	err = templates.Layout(templates.Card(templates.Sidebar(c, s)), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
