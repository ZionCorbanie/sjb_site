package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type PostHandler struct{
    postStore store.PostStore
}

type PostHandlerParams struct {
    PostStore store.PostStore
}

func NewPostHandler(params PostHandlerParams) *PostHandler {
	return &PostHandler{
        postStore: params.PostStore,
    }
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
    post, err := h.postStore.GetPost(postId)
    if err != nil {
        http.Error(w, "Error getting post", http.StatusInternalServerError)
        return
    }

	c := templates.Post(post)
	err = templates.Layout(c, "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
