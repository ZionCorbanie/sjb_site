package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PutPollHandler struct{
    store *dbstore.PollStore
}

type PutPollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewPutPollHandler(params PutPollHandlerParams) *PutPollHandler {
    return &PutPollHandler{
        store: params.PollStore,
    }
}

func (h *PutPollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)
    if err != nil {
        templates.PollError(err).Render(r.Context(), w)
        return
    }
    r.ParseForm()
    size := len(r.Form)-1
    options := make([]store.PollOption, size)

    for i := 0; i < size; i++ {
        options[i] = store.PollOption{
            Option: r.Form.Get(strconv.FormatInt(int64(i), 10)),
        }
    }

    poll := store.Poll{
        ID: uint(pollId),
        Title: r.Form.Get("title"),
        Options: options,
    }

    err = h.store.PutPoll(poll)
    if err != nil {
        templates.PollError(err).Render(r.Context(), w)
        return
    }

    polls, err := h.store.GetPolls()
    if err != nil {
        templates.PollError(err).Render(r.Context(), w)
        return
    }

    templates.Polls(polls).Render(r.Context(), w)
}
