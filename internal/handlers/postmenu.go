package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/store/dbstore"
	"time"
)

type PostMenuHandler struct{
    store *dbstore.MenuStore
}

type PostMenuHandlerParams struct {
    MenuStore *dbstore.MenuStore
}

func NewPostMenuHandler(params PostMenuHandlerParams) *PostMenuHandler {
    return &PostMenuHandler{
        store: params.MenuStore,
    }
}

func (h *PostMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	sendPopup(w, "Menu toegevoegd")
}
