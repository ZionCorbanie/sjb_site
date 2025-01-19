package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostsHandler struct{
    postStore store.PostStore
}

type PostsHandlerParams struct {
    PostsStore store.PostStore
}

func NewPostsHandler(params PostsHandlerParams) *PostsHandler {
	return &PostsHandler{
        postStore: params.PostsStore,
    }
}

func (h *PostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")
    if page == "" {
        page = "0"
    }
    pageInt, err := strconv.Atoi(page)

    user := middleware.GetUser(r.Context())
    admin := middleware.IsAdmin(r.Context())
    external := user == nil

    posts, err := h.postStore.GetPostsRange(pageInt*5, 5, admin, external)
    if err != nil {
        http.Error(w, "Error getting post", http.StatusInternalServerError)
        return
    }

	c := templates.Posts(posts, pageInt)
	err = templates.Layout(c, "Posts").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
