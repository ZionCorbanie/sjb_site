package handlers

import (
	"net/http"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type GetPollEditHandler struct{
    store *dbstore.PollStore
}

type GetPollEditHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewGetPollEditHandler(params GetPollEditHandlerParams) *GetPollEditHandler {
    return &GetPollEditHandler{
        store: params.PollStore,
    }
}

func (h *GetPollEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pollId := chi.URLParam(r, "pollId")
    poll, err := h.store.GetPoll(pollId)
	if err != nil {
		http.Error(w, "Error getting polls", http.StatusInternalServerError)
		return
	}

    err = templates.EditPoll(poll).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
