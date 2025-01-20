package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PatchAdminGroupEditHandler struct {
	GroupStore store.GroupStore
}

type PatchAdminGroupEditHandlerParams struct {
	GroupStore store.GroupStore
}

func NewPatchAdminGroupEditHandler(params PatchAdminGroupEditHandlerParams) *PatchAdminGroupEditHandler {
	return &PatchAdminGroupEditHandler{
		GroupStore: params.GroupStore,
	}
}

func (h *PatchAdminGroupEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupId, _ := strconv.ParseUint(chi.URLParam(r, "groupId"), 10, 64)

	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	name := r.FormValue("name")
	validateErr := h.GroupStore.ValidateInput(name, groupId)
	if validateErr != nil {
		templates.GroupError(validateErr).Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	description := r.FormValue("description")
	website := r.FormValue("website")
	//image := r.FormValue("image")

	GroupPatch := store.Group{
		ID:          uint(groupId),
		Email:       email,
		Website:     website,
		Name:        name,
		Description: description,
	}

	err := h.GroupStore.PatchGroup(GroupPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	templates.GroupError(nil).Render(r.Context(), w)
}
