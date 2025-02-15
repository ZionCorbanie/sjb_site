package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
    "time"
)

type PostCreatePostHandler struct{
    postStore store.PostStore
}

type PostCreatePostHandlerParams struct {
    PostStore store.PostStore
}

func NewPostCreatePostHandler(params PostCreatePostHandlerParams) *PostCreatePostHandler {
    return &PostCreatePostHandler{
        postStore: params.PostStore,
    }
}

func (h *PostCreatePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	author, _ := r.Context().Value(middleware.UserKey).(*store.User)
    date := time.Now()
    r.ParseForm()
    title := r.Form.Get("title")
    content := r.Form.Get("content")

    post := &store.Post{
        Author: *author,
        Title: title,
        Content: content,
        Date: date,
    }

    err := h.postStore.CreatePost(post)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }
}
