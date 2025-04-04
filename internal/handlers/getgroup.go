package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetGroupHandler struct {
	GroupStore     store.GroupStore
	GroupUserStore store.GroupUserStore
}

type GetGroupHandlerParams struct {
	GroupStore     store.GroupStore
	GroupUserStore store.GroupUserStore
}

func NewGroupHandler(params GetGroupHandlerParams) *GetGroupHandler {
	return &GetGroupHandler{
		GroupStore:     params.GroupStore,
		GroupUserStore: params.GroupUserStore,
	}
}

func (h *GetGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")
	group, err := h.GroupStore.GetGroup(groupId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	users, err := h.GroupUserStore.GetGroupUserByGroup(groupId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

    groups, title, err := h.GroupStore.GetSimelarGroups(group)
    if err != nil {
        err = templates.NotFound().Render(r.Context(), w)
        if err != nil {
            http.Error(w, "Error rendering template", http.StatusInternalServerError)
            return
        }
        return
    }

    fmt.Println("groups", groups)

    user := middleware.GetUser(r.Context())
    isVoorzitter := false
    for _, u := range *users {
        if u.UserID == user.ID {
            isVoorzitter = true
            break
        }
    }

	c := templates.Group(group, users, isVoorzitter)
	s := templates.SidebarGroup(groups, title)
	err = templates.BannerLayout(templates.Sidebar(templates.Card(c), templates.Card(s)), group.Image, group.Name).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
