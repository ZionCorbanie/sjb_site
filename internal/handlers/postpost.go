package handlers

import (
	"fmt"
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

    title := r.FormValue("title")
    content := r.FormValue("content")

    if title == "" || content == "" {
        http.Error(w, "Title and content are required", http.StatusBadRequest)
        return
    }

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

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

	w.Header().Set("HX-Redirect", fmt.Sprintf("/post/%d", post.ID))
	w.WriteHeader(http.StatusOK)
}
