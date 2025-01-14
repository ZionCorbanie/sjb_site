package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
)

type HomeHandler struct{
    postStore store.PostStore
}

type HomeHandlerParams struct {
    PostStore store.PostStore
}

func NewHomeHandler(params *HomeHandlerParams) *HomeHandler {
	return &HomeHandler{
        postStore: params.PostStore,
    }
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, _ := r.Context().Value(middleware.UserKey).(*store.User)

    posts, err := h.postStore.GetPostsRange(0, 10)
    if err != nil {
        http.Error(w, "Error getting posts", http.StatusInternalServerError)
        return
    }

	c := templates.Index(user, posts)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
