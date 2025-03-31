package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetUserHandler struct {
	userStore store.UserStore
	groupUserStore store.GroupUserStore
}

type GetUserHandlerParams struct {
	UserStore store.UserStore
	GroupUserStore store.GroupUserStore
}

func NewUserHandler(params GetUserHandlerParams) *GetUserHandler {
	return &GetUserHandler{
		userStore: params.UserStore,
        groupUserStore: params.GroupUserStore,
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

    groups, err:= h.groupUserStore.GetGroupUserByUser(userId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := templates.User(user, groups)
	s := templates.SidebarUser()
	err = templates.Layout(templates.Sidebar(templates.Card(c), templates.Card(s)), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
