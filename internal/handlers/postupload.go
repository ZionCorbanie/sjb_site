package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sjb_site/internal/middleware"
)

type PostUploadHandler struct{}

func NewPostUploadHandler() *PostUploadHandler {
	return &PostUploadHandler{}
}

func (h *PostUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

    //handle file upload
    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Failed to get image", http.StatusBadRequest)
        return
    }

    defer file.Close()

    uploadDir := "static/uploads/posts"
    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        os.MkdirAll(uploadDir, os.ModePerm)
    }

    filePath := filepath.Join(uploadDir, header.Filename)

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

    fmt.Fprintf(w, "{\"path\": \"/%s\"}", filePath)
}
