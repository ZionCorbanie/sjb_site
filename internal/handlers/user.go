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

type UserHandler struct {
	userStore store.UserStore
	groupUserStore store.GroupUserStore
}

type UserHandlerParams struct {
	UserStore store.UserStore
	GroupUserStore store.GroupUserStore
}

func NewUserHandler(params UserHandlerParams) *UserHandler {
	return &UserHandler{
		userStore: params.UserStore,
        groupUserStore: params.GroupUserStore,
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	c := templates.Users()
	err := templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

//Search returns list of users from search
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	users, err := h.userStore.SearchUsers(search)

	if err != nil {
		return
	}

	c := templates.RenderUsers(users)
	c.Render(r.Context(), w)
}

func (h *UserHandler) Show(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	user, err := h.userStore.GetUserById(userId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

    groups, err:= h.groupUserStore.GetGroupUserByUser(userId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := templates.User(user, groups)
	s := templates.SidebarUser()
	err = templates.Layout(templates.Sidebar(templates.Card(c), templates.Card(s)), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Edit(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	user, err := h.userStore.GetUserById(userId)
	if err != nil {
		err = templates.NotFound().Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
		return
	}

	c := templates.UserEdit(user)
	err = templates.Layout(templates.Card(c), "Sint Jansbrug").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
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

	SendPopup(w, "Profiel aangepast")
	w.Header().Add("Hx-Redirect", fmt.Sprintf("/webalmanak/leden/%d", userId))
}
