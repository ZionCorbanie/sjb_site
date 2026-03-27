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

type CommentHandler struct{
    commentStore store.CommentStore
}

type CommentHandlerParams struct {
    CommentStore store.CommentStore
}

func NewCommentHandler(params CommentHandlerParams) *CommentHandler {
	return &CommentHandler{
        commentStore: params.CommentStore,
    }
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
    comments, err := h.commentStore.GetCommentsByPost(postId)
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

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    date := time.Now()
    postID, err := strconv.ParseUint(chi.URLParam(r, "postId"), 10, 64)
    comment := &store.Comment{
        Content: r.Form.Get("content"),
        Date: date,
        Author: *middleware.GetUser(r.Context()),
        PostID: uint(postID),
    }

    err = h.commentStore.CreateComment(comment)
    if err != nil {
        templates.CommentError(err.Error()).Render(r.Context(), w)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/comments/%d", postID), http.StatusSeeOther)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
    commentID := chi.URLParam(r, "commentId")
    postID := chi.URLParam(r, "postId")

    comment, err := h.commentStore.GetComment(commentID)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if comment.Author.ID != middleware.GetUser(r.Context()).ID && !middleware.IsAdmin(r.Context()) {
        w.WriteHeader(http.StatusForbidden)
        return
    }

    err = h.commentStore.DeleteComment(commentID)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/comments/%s", postID), http.StatusSeeOther)
}
