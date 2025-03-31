package handlers

import (
	"net/http"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type DeletePollHandler struct{
    store *dbstore.PollStore
}

type DeletePollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewDeletePollHandler(params DeletePollHandlerParams) *DeletePollHandler {
    return &DeletePollHandler{
        store: params.PollStore,
    }
}

func (h *DeletePollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pollId := chi.URLParam(r, "pollId")
    
    err := h.store.DeletePoll(pollId)
	if err != nil {
		http.Error(w, "Error getting polls", http.StatusInternalServerError)
		return
	}

    polls, err := h.store.GetPolls()
    err = templates.Polls(polls).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
