package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type CommentsHandler struct{
    postStore store.CommentStore
}

type CommentsHandlerParams struct {
    CommentStore store.CommentStore
}

func NewCommentsHandler(params CommentsHandlerParams) *CommentsHandler {
	return &CommentsHandler{
        postStore: params.CommentStore,
    }
}

func (h *CommentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
    comments, err := h.postStore.GetCommentsByPost(postId)
    if err != nil {
        http.Error(w, "Error getting comments", http.StatusInternalServerError)
        return
    }

    err = templates.Comments(comments).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
