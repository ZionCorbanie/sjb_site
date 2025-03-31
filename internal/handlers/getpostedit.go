package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetPostEditHandler struct {
	postStore store.PostStore
}

type GetPostEditHandlerParams struct {
	PostStore store.PostStore
}

func NewPostEditHandler(params GetPostEditHandlerParams) *GetPostEditHandler {
	return &GetPostEditHandler{
		postStore: params.PostStore,
	}
}

func (h *GetPostEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	post, err := h.postStore.GetPost(postId)

	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := templates.EditPost(post)
	err = templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
