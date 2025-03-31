package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"
)

type HomeHandler struct{
    postStore store.PostStore
    menuStore store.MenuStore
}

type HomeHandlerParams struct {
    PostStore store.PostStore
    MenuStore store.MenuStore
}

func NewHomeHandler(params *HomeHandlerParams) *HomeHandler {
	return &HomeHandler{
        postStore: params.PostStore,
        menuStore: params.MenuStore,
    }
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, _ := r.Context().Value(middleware.UserKey).(*store.User)

    admin := middleware.IsAdmin(r.Context())
    external := user == nil

    posts, err := h.postStore.GetPostsRange(0, 3, admin, external)
    if err != nil {
        http.Error(w, "Error getting posts", http.StatusInternalServerError)
        return
    }
    menuId := time.Now().Unix()/(60*60*24)
    menu := templates.MenuDay(h.menuStore.GetMenu(strconv.FormatInt(menuId, 10)))

	c := templates.Index(user, posts, menu)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
