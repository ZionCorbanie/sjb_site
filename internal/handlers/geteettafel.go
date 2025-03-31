package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
	"time"
)

type GetEettafelHandler struct {
}

func NewEettafelHandler() *GetEettafelHandler {
	return &GetEettafelHandler{
	}
}

func (h *GetEettafelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    pageInt := int(time.Now().Unix() / 60 / 60 / 24)

    //Get maandag of week or next week if saterday or sunday
    weekDay := int(time.Now().Weekday())
    offset := 1-weekDay

    if weekDay == 6 {
        offset += 7
    }

    pageInt += offset

    c := templates.Eettafel(pageInt)
    err := templates.Layout(templates.Card(c), "Eettafel").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
