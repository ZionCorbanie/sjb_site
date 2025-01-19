package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PatchUserEditHandler struct {
	userStore store.UserStore
}

type PatchUserEditHandlerParams struct {
	UserStore store.UserStore
}

func NewPatchtUserEditHandler(params PatchUserEditHandlerParams) *PatchUserEditHandler {
	return &PatchUserEditHandler{
		userStore: params.UserStore,
	}
}

func (h *PatchUserEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	if middleware.GetUser(r.Context()).ID != uint(userId) && !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	email := r.FormValue("email")
	address := r.FormValue("address")
	phone := r.FormValue("phone")

	validateErr := h.userStore.ValidateInput(email, address, userId)
	if validateErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	userPatch := store.User{
		ID:          uint(userId),
		Email:       email,
		Adres:       address,
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
