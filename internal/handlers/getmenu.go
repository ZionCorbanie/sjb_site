package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetMenuHandler struct {
	MenuStore     store.MenuStore
}

type GetMenuHandlerParams struct {
	MenuStore     store.MenuStore
}

func NewMenuHandler(params GetMenuHandlerParams) *GetMenuHandler {
	return &GetMenuHandler{
        MenuStore:     params.MenuStore,
	}
}

func (h *GetMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	menuId := chi.URLParam(r, "menuId")
	menu := h.MenuStore.GetMenu(menuId)

    err := templates.MenuDay(menu).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
