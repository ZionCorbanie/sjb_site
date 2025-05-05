package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store"
	"time"

	"github.com/google/uuid"
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

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }


    title := r.FormValue("title")
    content := r.FormValue("content")
    published := r.FormValue("publiek") == "on"
    external := r.FormValue("extern") == "on"

    if title == "" || content == "" {
        http.Error(w, "Title and content are required", http.StatusBadRequest)
        return
    }

    post := &store.Post{
        Author: *author,
        Title: title,
        Content: content,
        Date: date,
        Published: published,
        External: external,
    }

    file, header, err := r.FormFile("image")
    if err == nil {
        defer file.Close()
        
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

        post.Image = "/"+filePath
    }else {
        post.Image = "/static/img/placeholder-group.png"
    }

    err = h.postStore.CreatePost(post)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

	sendPopup(w, "Post aangemaakt")
	w.Header().Set("HX-Redirect", fmt.Sprintf("/post/%d", post.ID))
	w.WriteHeader(http.StatusOK)
}
