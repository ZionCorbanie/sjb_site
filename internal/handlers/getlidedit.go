package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)


type GetLidEditHandler struct {
	userStore  store.UserStore
}

type GetLidEditHandlerParams struct {
	UserStore         store.UserStore
}

func NewLidEditHandler(params GetLidEditHandlerParams) *GetLidEditHandler {
	return &GetLidEditHandler{
        userStore: params.UserStore,
    }
}

func (h *GetLidEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
    fmt.Println(userId)
    user, err := h.userStore.GetUserById(uint(userId))
	if err != nil {
        err = templates.NotFound().Render(r.Context(), w)
        if err != nil {
            http.Error(w, "Error rendering template", http.StatusInternalServerError)
            return
        }
		return
	}

	c := templates.LidEdit(user)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
