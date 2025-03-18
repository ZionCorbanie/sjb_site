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

type PatchPostHandler struct {
	postStore store.PostStore
}

type PatchPostHandlerParams struct {
	PostStore store.PostStore
}

func NewPatchPostHandler(params PatchPostHandlerParams) *PatchPostHandler {
	return &PatchPostHandler{
		postStore: params.PostStore,
	}
}

func (h *PatchPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	postId, _ := strconv.ParseUint(chi.URLParam(r, "postId"), 10, 64)

	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

    title := r.FormValue("title")
    content := r.FormValue("content")

    postPatch := store.Post{
        ID: uint(postId),
        Title: title,
        Content: content,
    }

    err := h.postStore.PatchPost(postPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/post/%d", postId))
}
