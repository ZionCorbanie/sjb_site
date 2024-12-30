package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)


type GetLidHandler struct {
	userStore  store.UserStore
}

type GetLidHandlerParams struct {
	UserStore         store.UserStore
}

func NewLedenHandler(params GetLidHandlerParams) *GetLidHandler {
	return &GetLidHandler{
        userStore: params.UserStore,
    }
}

func (h *GetLidHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
    user, err := h.userStore.GetUserById(uint(userId))
	if err != nil {
        err = templates.NotFound().Render(r.Context(), w)
        if err != nil {
            http.Error(w, "Error rendering template", http.StatusInternalServerError)
            return
        }
		return
	}

	c := templates.Lid(user)
    s := templates.SidebarLid()
	err = templates.Layout(templates.Sidebar(c,s), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
