package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PatchAdminUserEditHandler struct {
	userStore store.UserStore
}

type PatchAdminUserEditHandlerParams struct {
	UserStore store.UserStore
}

func NewPatchAdminUserEditHandler(params PatchAdminUserEditHandlerParams) *PatchAdminUserEditHandler {
	return &PatchAdminUserEditHandler{
		userStore: params.UserStore,
	}
}

func (h *PatchAdminUserEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	if middleware.GetUser(r.Context()).ID != uint(userId) && !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	email := r.FormValue("email")
	address := r.FormValue("address")
	phone := r.FormValue("phone")
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	username := r.FormValue("username")
	//image := r.FormValue("image")

	userPatch := store.User{
		ID:          uint(userId),
		Email:       email,
		Adres:       address,
		PhoneNumber: phone,
		FirstName:   firstname,
		LastName:    lastname,
		Username:    username,
	}

	err := h.userStore.PatchUser(userPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	w.Header().Add("Hx-Redirect", "/admin/leden/")
}
