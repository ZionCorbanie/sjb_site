package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"

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
	menu, err := h.MenuStore.GetMenu(menuId)
	if err != nil {
        menu = &store.Menu{}
        menuIdInt, _ := strconv.ParseInt(menuId, 10, 64)
        menu.ID = uint(menuIdInt)
        menu.Name = "Onbekend"

        menu.Date = time.Unix(menuIdInt*60*60*24, 0)
        if menu.Date.Weekday() == time.Saturday || menu.Date.Weekday() == time.Sunday {
            menu.Name = "Eettafel gesloten"
        }
	}

	err = templates.MenuDay(menu).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
