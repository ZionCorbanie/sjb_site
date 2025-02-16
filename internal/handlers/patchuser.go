package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
)

type PatchUserHandler struct {
	userStore store.UserStore
}

type PatchUserHandlerParams struct {
	UserStore store.UserStore
}

func NewPatchtUserHandler(params PatchUserHandlerParams) *PatchUserHandler {
	return &PatchUserHandler{
		userStore: params.UserStore,
	}
}

func (h *PatchUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	if middleware.GetUser(r.Context()).ID != uint(userId) && !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	email := r.FormValue("email")
	address := r.FormValue("address")
	phone := r.FormValue("phone")

	userPatch := store.User{
		ID:           uint(userId),
		Email:        email,
		Adres:        address,
		PhoneNumber: phone,
	}

	err := h.userStore.PatchUser(userPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/webalmanak/leden/%d", userId))
}
