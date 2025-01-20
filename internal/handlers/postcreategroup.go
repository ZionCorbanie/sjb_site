package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type PostCreateGroupHandler struct {
	groupStore store.GroupStore
}

type PostCreateGroupHandlerParams struct {
	GroupStore store.GroupStore
}

func NewPostCreateGroupHandler(params PostCreateGroupHandlerParams) *PostCreateGroupHandler {
	return &PostCreateGroupHandler{
		groupStore: params.GroupStore,
	}
}

func (h *PostCreateGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.Form.Get("name")
	validateErr := h.groupStore.ValidateInput(name, 0)
	if validateErr != nil {
		templates.GroupError(validateErr).Render(r.Context(), w)
		return
	}

	group := &store.Group{
		GroupType:   r.Form.Get("groupType"),
		Name:        name,
		Description: r.Form.Get("description"),
		Email:       r.Form.Get("email"),
		Website:     r.Form.Get("website"),
	}

	err := h.groupStore.CreateGroup(group)
	if err != nil {
		templates.GroupError(err).Render(r.Context(), w)
		return
	}

	templates.GroupError(nil).Render(r.Context(), w)
}
