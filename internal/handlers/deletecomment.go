package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"

	"github.com/go-chi/chi/v5"
)

type DeleteCommentHandler struct{
    store store.CommentStore
}

type DeleteCommentHandlerParams struct {
    CommentStore store.CommentStore
}

func NewDeleteCommentHandler(params DeleteCommentHandlerParams) *DeleteCommentHandler {
    return &DeleteCommentHandler{
        store: params.CommentStore,
    }
}

func (h *DeleteCommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "commentId")
    postID := chi.URLParam(r, "postId")

    //get post
    comment, err := h.store.GetComment(commentID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if comment.Author.ID != middleware.GetUser(r.Context()).ID || !middleware.IsAdmin(r.Context()) {
        w.WriteHeader(http.StatusForbidden)
        return
    }

    err = h.store.DeleteComment(commentID)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/comments/%s", postID), http.StatusSeeOther)
}
