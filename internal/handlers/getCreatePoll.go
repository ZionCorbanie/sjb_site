package handlers

import (
	"net/http"
	"sjb_site/internal/templates"
)

type GetCreatePollHandler struct{}

func NewGetCreatePollHandler() *GetCreatePollHandler {
	return &GetCreatePollHandler{}
}

func (h *GetCreatePollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    c := templates.CreatePoll()
	err := templates.Layout(templates.Card(c), "Poll maken").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
