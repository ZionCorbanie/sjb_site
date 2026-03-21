package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type MenuHandler struct {
	menuStore     store.MenuStore
}

type MenuHandlerParams struct {
	MenuStore     store.MenuStore
}

func NewMenuHandler(params MenuHandlerParams) *MenuHandler {
	return &MenuHandler{
        menuStore:     params.MenuStore,
	}
}

func (h *MenuHandler) Day(w http.ResponseWriter, r *http.Request) {
	menuId := chi.URLParam(r, "menuId")
	menu := h.menuStore.GetMenu(menuId)

    err := templates.MenuDay(menu).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *MenuHandler) Week(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "pageId")
    pageInt, err := strconv.Atoi(page)
    if err != nil {
    }

	menus, err := h.menuStore.GetMenuRange(pageInt, 5)
    for i, menu := range *menus {
        if menu.ID == 0 {
            menu.ID = uint(pageInt+i)
            menu.Name = "Onbekend"
            menu.Date = time.Unix(int64((pageInt+i)*60*60*24), 0)
        }
    }

	err = templates.WeekMenu(menus).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
