package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type GetWeekMenuHandler struct {
	WeekMenuStore     store.MenuStore
}

type GetWeekMenuHandlerParams struct {
	WeekMenuStore     store.MenuStore
}

func NewWeekMenuHandler(params GetWeekMenuHandlerParams) *GetWeekMenuHandler {
	return &GetWeekMenuHandler{
        WeekMenuStore:     params.WeekMenuStore,
	}
}

func (h *GetWeekMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "pageId")
    pageInt, err := strconv.Atoi(page)
    if err != nil {
    }


	menus, err := h.WeekMenuStore.GetMenuRange(pageInt, 5)
    for i, menu := range menus {
        if menu.ID == 0 {
            menu.ID = uint(pageInt+i)
            menu.Name = "Menu onbekend"
            menu.Date = time.Unix(int64((pageInt+i)*60*60*24), 0)
        }
    }

	err = templates.WeekMenu(menus).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
