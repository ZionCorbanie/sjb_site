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

	"github.com/go-chi/chi/v5"
)

type PatchUserEditHandler struct {
	userStore store.UserStore
}

type PatchUserEditHandlerParams struct {
	UserStore store.UserStore
}

func NewPatchtUserEditHandler(params PatchUserEditHandlerParams) *PatchUserEditHandler {
	return &PatchUserEditHandler{
		userStore: params.UserStore,
	}
}

func (h *PatchUserEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)

	if middleware.GetUser(r.Context()).ID != uint(userId) && !middleware.IsAdmin(r.Context()) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

    //parse form
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

	email := r.FormValue("email")
	address := r.FormValue("address")
	phone := r.FormValue("phone")

	userPatch := store.User{
		ID:           uint(userId),
		Email:        email,
		Adres:        address,
		PhoneNumber: phone,
	}

    //handle file upload
    file, header, err := r.FormFile("image")
    if err == nil {
        defer file.Close()

        fileName := fmt.Sprint("%d%s",userId,filepath.Ext(header.Filename))
        uploadDir := "static/uploads/user"
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

        userPatch.Image = "/"+filePath
        fmt.Printf("File uploaded successfully: %s\n", filePath)
    }


	err = h.userStore.PatchUser(userPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	w.Header().Add("Hx-Redirect", fmt.Sprintf("/webalmanak/leden/%d", userId))
}
