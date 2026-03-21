package admin

import (
	"fmt"
	"net/http"
	"sjb_site/internal/handlers"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)


type UserHandler struct {
	userStore store.UserStore
}

type UserHandlerParams struct {
	UserStore store.UserStore
}

func NewUserHandler(params UserHandlerParams) *UserHandler {
	return &UserHandler{
		userStore: params.UserStore,
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	c := templates.UserManagement()
	err := templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

	c := templates.RenderUserManagement(users)
	c.Render(r.Context(), w)
}

func (h *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
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

	handlers.SendPopup(w, fmt.Sprintf("Account van %s aangepast", firstname))
	w.Header().Add("Hx-Redirect", "/admin/leden/")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	err := h.userStore.DeleteUser(userId)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Hx-Redirect", "/admin/leden/")
}
