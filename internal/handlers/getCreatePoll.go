package handlers

import (
	"net/http"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
)

type GetCreatePollHandler struct{
    store *dbstore.PollStore
}

type GetCreatePollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewGetCreatePollHandler(params GetCreatePollHandlerParams) *GetCreatePollHandler {
    return &GetCreatePollHandler{
        store: params.PollStore,
    }
}

func (h *GetCreatePollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    polls, err := h.store.GetPolls()
	if err != nil {
		http.Error(w, "Error getting polls", http.StatusInternalServerError)
		return
	}
    c := templates.Polls(polls)
	err = templates.Layout(templates.Card(c), "Poll maken").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
