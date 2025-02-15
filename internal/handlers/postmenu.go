package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"time"
)

type PostCreateMenuHandler struct{
    store *dbstore.MenuStore
}

type PostCreateMenuHandlerParams struct {
    MenuStore *dbstore.MenuStore
}

func NewPostCreateMenuHandler(params PostCreateMenuHandlerParams) *PostCreateMenuHandler {
    return &PostCreateMenuHandler{
        store: params.MenuStore,
    }
}

func (h *PostCreateMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    date, err := time.Parse("2006-01-02", r.Form.Get("date"))
    if err != nil {
        templates.MenuError(err).Render(r.Context(), w)
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
        templates.MenuError(err).Render(r.Context(), w)
        return
    }

    templates.MenuError(nil).Render(r.Context(), w)
}
