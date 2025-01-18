package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type PostCommentHandler struct{
    store store.CommentStore
}

type PostCommentHandlerParams struct {
    CommentStore store.CommentStore
}

func NewPostCommentHandler(params PostCommentHandlerParams) *PostCommentHandler {
    return &PostCommentHandler{
        store: params.CommentStore,
    }
}

func (h *PostCommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    date := time.Now()
    postID, err := strconv.ParseUint(chi.URLParam(r, "postId"), 10, 64)
    comment := &store.Comment{
        Content: r.Form.Get("content"),
        Date: date,
        Author: *middleware.GetUser(r.Context()),
        PostID: uint(postID),
    }

    err = h.store.CreateComment(comment)
    if err != nil {
        templates.CommentError(err.Error()).Render(r.Context(), w)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/comments/%d", postID), http.StatusSeeOther)
}
