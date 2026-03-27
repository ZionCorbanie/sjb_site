package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type RegisterHandler struct {
	userStore store.UserStore
}

type RegisterHandlerParams struct {
	UserStore store.UserStore
}

func NewRegisterHandler(params RegisterHandlerParams) *RegisterHandler {
	return &RegisterHandler{
		userStore: params.UserStore,
	}
}

func (h *RegisterHandler) Show(w http.ResponseWriter, r *http.Request) {
	c := templates.RegisterPage()
	err := templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}

//TODO deze functie mag op de uiteindelijke site niet bestaan
func (h *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {
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
