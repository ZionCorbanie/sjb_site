package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

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

func (h *PostHandler) Show(w http.ResponseWriter, r *http.Request) {
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
	err = templates.BannerLayout(templates.Card(c), post.Image, post.Title).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	//page := chi.URLParam(r, "page")
	page := r.URL.Query().Get("page")
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
	err = templates.Layout(templates.Card(c), "Posts").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
