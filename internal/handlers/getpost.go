package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
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

    user := middleware.GetUser(r.Context())
    if !middleware.IsAdmin(r.Context()){
        if !post.Published{
            c := templates.NotFound()
            err = templates.Layout(c, "404").Render(r.Context(), w)
            if err != nil {
                http.Error(w, "Error rendering template", http.StatusInternalServerError)
                return
            }
            return
        }
        if user == nil && !post.External{
            target := r.URL.Path
            http.Redirect(w, r, fmt.Sprintf("/login?redirect=%s", target), http.StatusFound)
			return
        }
    }


	c := templates.Post(post)
	err = templates.BannerLayout(c, post.Image, post.Title).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
