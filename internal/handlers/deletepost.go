package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"strings"

	"github.com/go-chi/chi/v5"
)

type DeletePostHandler struct{
    postStore store.PostStore
}

type DeletePostHandlerParams struct {
    PostStore store.PostStore
}

func NewDeletePostHandler(params DeletePostHandlerParams) *DeletePostHandler {
    return &DeletePostHandler{
        postStore: params.PostStore,
    }
}

func (h *DeletePostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")

	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

    post, err := h.postStore.GetPost(chi.URLParam(r, "postId"))
    if strings.Contains(post.Image, "uploads/posts") {
        err = os.Remove(post.Image[1:])
        if err != nil {
            fmt.Printf("Error deleting file: %s\n", post.Image)
        }
    }

    err = h.postStore.DeletePost(postId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sendPopup(w, "Post verwijderd")
	w.Header().Add("Hx-Redirect", "/posts")
}
