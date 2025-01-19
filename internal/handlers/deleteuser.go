package handlers

import (
	"net/http"
	"sjb_site/internal/store"

	"github.com/go-chi/chi/v5"
)

type DeleteUserHandler struct {
	userStore store.UserStore
}

type DeleteUserHandlerParams struct {
	UserStore store.UserStore
}

func NewDeleteUserHandler(params DeleteUserHandlerParams) *DeleteUserHandler {
	return &DeleteUserHandler{
		userStore: params.UserStore,
	}
}

func (h *DeleteUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	err := h.userStore.DeleteUser(userId)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Hx-Redirect", "/admin/leden/")
}
