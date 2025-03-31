package handlers

import (
	"fmt"
	"net/http"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type PostPollActivateHandler struct{
    store *dbstore.PollStore
}

type PostPollActivateHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewPostPollActivateHandler(params PostPollActivateHandlerParams) *PostPollActivateHandler {
    return &PostPollActivateHandler{
        store: params.PollStore,
    }
}

func (h *PostPollActivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pollId := chi.URLParam(r, "pollId")

    fmt.Println("Poll ID: ", pollId)
    err := h.store.Activate(pollId)
    if  err != nil {
        http.Error(w, "Error making poll active", http.StatusInternalServerError)
        return
    }

    polls, err := h.store.GetPolls()

    if err != nil {
        http.Error(w, "Error getting polls", http.StatusInternalServerError)
        return
    }

    templates.Polls(polls).Render(r.Context(), w)
}
