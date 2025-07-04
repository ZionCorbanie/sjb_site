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
)

type PatchUserHandler struct {
	userStore store.UserStore
}

type PatchUserHandlerParams struct {
	UserStore store.UserStore
}

func NewPatchtUserHandler(params PatchUserHandlerParams) *PatchUserHandler {
	return &PatchUserHandler{
		userStore: params.UserStore,
	}
}

func (h *PatchUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

        oldUser, err := h.userStore.GetUserById(chi.URLParam(r, "userId"))
        if strings.Contains(oldUser.Image, "uploads/user") {
            err = os.Remove(oldUser.Image[1:])
            if err != nil {
                fmt.Printf("Error deleting file: %s\n", oldUser.Image)
            }
        }
        

        fileName := fmt.Sprintf("%d%s",userId,filepath.Ext(header.Filename))
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
    }


	err = h.userStore.PatchUser(userPatch)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c := templates.RegisterError()
		c.Render(r.Context(), w)
		return
	}

	sendPopup(w, "Account van aangepast")
	w.Header().Add("Hx-Redirect", fmt.Sprintf("/webalmanak/leden/%d", userId))
}
