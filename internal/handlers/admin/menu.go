package admin

import (
	"net/http"
	"sjb_site/internal/handlers"
	"sjb_site/internal/store"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"time"
)

type MenuHandler struct{
    store *dbstore.MenuStore
}

type MenuHandlerParams struct {
    MenuStore *dbstore.MenuStore
}

func NewMenuHandler(params MenuHandlerParams) *MenuHandler {
    return &MenuHandler{
        store: params.MenuStore,
    }
}

func (h *MenuHandler) New(w http.ResponseWriter, r *http.Request) {
    c := templates.CreateMenu()
	err := templates.Layout(templates.Card(c), "Menu toevoegen").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    date, err := time.Parse("2006-01-02", r.Form.Get("date"))
    if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    id := uint(date.Unix() / (60*60*24))
    menu := &store.Menu{
        ID: id,
        Date: date,
        Name: r.Form.Get("gerecht"),
        Basis: r.Form.Get("basis"),
        Vlees: r.Form.Get("vlees"),
        Vega: r.Form.Get("vega"),
        Toe: r.Form.Get("toe"),
    }

    err = h.store.CreateMenu(menu)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	handlers.SendPopup(w, "Menu toegevoegd")
}
