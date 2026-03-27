package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sjb_site/internal/handlers"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PromoHandler struct{
    promoStore store.PromoStore
}

type PromoHandlerParams struct {
    PromoStore store.PromoStore
}

func NewPromoHandler(params PromoHandlerParams) *PromoHandler {
	return &PromoHandler{
        promoStore: params.PromoStore,
    }
}

func (h *PromoHandler) New(w http.ResponseWriter, r *http.Request) {
	promos, err := h.promoStore.GetAllPromos()
	if err != nil {
		http.Error(w, "Error getting promos", http.StatusInternalServerError)
		return
	}

	err = templates.Layout(templates.Card(templates.Promos(promos)), "Promo maken").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}


func (h *PromoHandler) Edit(w http.ResponseWriter, r *http.Request) {
	promoId := chi.URLParam(r, "promoId")

	promo, err := h.promoStore.GetPromo(promoId)
	if err != nil {
		http.Error(w, "Error getting promos", http.StatusInternalServerError)
		return
	}

	err = templates.EditPromo(promo).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *PromoHandler) Put(w http.ResponseWriter, r *http.Request) {
    promoId, err := strconv.ParseInt(chi.URLParam(r, "promoId"), 10, 64)
	if err != nil {
		promoId = 0
	}

    err = r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    title := r.FormValue("title")
    content := r.FormValue("content")
    startDate, err := time.Parse("2006-01-02", r.Form.Get("startDate"))
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
    endDate, err := time.Parse("2006-01-02", r.Form.Get("endDate"))
	if err != nil {
		http.Error(w, "Invalid end date", http.StatusBadRequest)
		return
	}

    if title == "" || content == "" {
        http.Error(w, "Promo heeft titel en content nodig", http.StatusBadRequest)
        return
    }

    promo := &store.Promo{
		ID: uint(promoId),
		Title: title,
		Description: content,
		StartDate: startDate,
		EndDate: endDate,
	}
    file, header, err := r.FormFile("image")
    if err == nil {
        defer file.Close()
        
        fileName := fmt.Sprintf("%d%s",uuid.New().ID(),filepath.Ext(header.Filename))
        uploadDir := "static/uploads/promo"
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

        promo.Image = "/"+filePath
    }else {
        promo.Image = "/static/img/placeholder-group.png"
    }

    err = h.promoStore.SavePromo(promo)
    if err != nil {
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

	promos, err := h.promoStore.GetAllPromos()
	if err != nil {
		http.Error(w, "Error getting promos", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST"{
		handlers.SendPopup(w, "Promo aangemaakt")
	}else{
		handlers.SendPopup(w, "Promo aangepast")
	}
	templates.Promos(promos).Render(r.Context(), w)
}

func (h *PromoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	promoId := chi.URLParam(r, "promoId")
    
    err := h.promoStore.DeletePromo(promoId)
	if err != nil {
		http.Error(w, "Error deleting promos", http.StatusInternalServerError)
		return
	}

    promos, err := h.promoStore.GetAllPromos()
    err = templates.Promos(promos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	handlers.SendPopup(w, "Promo verwijderd")
}

func (h *PromoHandler) DeleteInactive(w http.ResponseWriter, r *http.Request) {
	err := h.promoStore.DeleteInactivePromos()
	if err != nil {
		http.Error(w, "Error deleting inactive promos", http.StatusInternalServerError)
		return
	}

	promos, err := h.promoStore.GetAllPromos()
	err = templates.Promos(promos).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	handlers.SendPopup(w, "Alle inactieve posts verwijderd")
}
