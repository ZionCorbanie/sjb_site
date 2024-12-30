package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
    "fmt"
	"github.com/go-chi/chi/v5"
)

type PatchLidEditHandler struct {
	userStore store.UserStore
}

type PatchLidEditHandlerParams struct {
	UserStore store.UserStore
}

func NewPatchtLidEditHandler(params PatchLidEditHandlerParams) *PatchLidEditHandler {
	return &PatchLidEditHandler{
		userStore: params.UserStore,
	}
}

func (h *PatchLidEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

    if middleware.GetUser(r.Context()).ID != uint(userId) && !middleware.IsAdmin(r.Context()) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

	email := r.FormValue("email")
	address := r.FormValue("address")
    phone := r.FormValue("phone")

    userPatch := store.User{
        ID: uint(userId),
        Email: email,
        Adres: address,
        Phone_number: phone,
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
