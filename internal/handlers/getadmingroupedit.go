package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetAdminGroupEditHandler struct {
	GroupStore store.GroupStore
}

type GetAdminGroupEditHandlerParams struct {
	GroupStore store.GroupStore
}

func NewAdminGroupEditHandler(params GetAdminGroupEditHandlerParams) *GetAdminGroupEditHandler {
	return &GetAdminGroupEditHandler{
		GroupStore: params.GroupStore,
	}
}

func (h *GetAdminGroupEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	GroupId := chi.URLParam(r, "groupId")
	Group, err := h.GroupStore.GetGroup(GroupId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}
	Groups, err := h.GroupStore.GetGroupsByType(Group.GroupType)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	s := templates.AdminSidebarGroup(Group, Groups)
	c := templates.GroupEditAdmin(Group)
	err = templates.Layout(templates.Sidebar(c, s), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
