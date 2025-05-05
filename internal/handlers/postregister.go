package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PostRegisterHandler struct {
	userStore store.UserStore
}

type PostRegisterHandlerParams struct {
	UserStore store.UserStore
}

func NewPostRegisterHandler(params PostRegisterHandlerParams) *PostRegisterHandler {
	return &PostRegisterHandler{
		userStore: params.UserStore,
	}
}

//TODO deze functie mag op de uiteindelijke site niet bestaan
func (h *PostRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := h.userStore.CreateUser(username, password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	c := templates.RegisterSuccess()
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}
