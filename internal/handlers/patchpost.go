package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
    published := r.FormValue("publiek") == "on"
    external := r.FormValue("extern") == "on"

    oldPost, err := h.postStore.GetPost(chi.URLParam(r, "postId"))
    postPatch := store.Post{
        ID: uint(postId),
        Title: title,
        Content: content,
        Image: oldPost.Image,
        Published: published,
        External: external,
    }

    file, header, err := r.FormFile("image")
    if err == nil {
        defer file.Close()

        if strings.Contains(oldPost.Image, "uploads/posts") {
            err = os.Remove(oldPost.Image[1:])
            if err != nil {
                fmt.Printf("Error deleting file: %s\n", oldPost.Image)
            }
        }

        fileName := fmt.Sprintf("%d%s",uuid.New().ID(),filepath.Ext(header.Filename))
        uploadDir := "static/uploads/posts"
        os.MkdirAll(uploadDir, os.ModePerm)
        filePath := filepath.Join(uploadDir, fileName)

        dst, err := os.Create(filePath)
        if err != nil {
            http.Error(w, "Unable to save file", http.StatusInternalServerError)
            return
        }
        defer dst.Close()

        _, err = io.Copy(dst, file)
        if err != nil {
            http.Error(w, "Error saving file", http.StatusInternalServerError)
            return
        }

        postPatch.Image = "/"+filePath
    }

    err = h.postStore.PatchPost(postPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/post/%d", postId))
}
