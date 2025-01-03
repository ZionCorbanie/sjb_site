package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetUserHandler struct {
	userStore store.UserStore
}

type GetUserHandlerParams struct {
	UserStore store.UserStore
}

func NewUserHandler(params GetUserHandlerParams) *GetUserHandler {
	return &GetUserHandler{
		userStore: params.UserStore,
	}
}

func (h *GetUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	c := templates.User(user)
	s := templates.SidebarUser()
	err = templates.Layout(templates.Sidebar(c, s), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
