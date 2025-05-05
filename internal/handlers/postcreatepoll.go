package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"strconv"
)

type PostCreatePollHandler struct{
    store *dbstore.PollStore
}

type PostCreatePollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewPostCreatePollHandler(params PostCreatePollHandlerParams) *PostCreatePollHandler {
    return &PostCreatePollHandler{
        store: params.PollStore,
    }
}

func (h *PostCreatePollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    size := len(r.Form)-1

    options := make([]store.PollOption, size)

    for i := 0; i < size; i++ {
        options[i] = store.PollOption{
            Option: r.Form.Get(strconv.FormatInt(int64(i), 10)),
        }
    }

    poll := &store.Poll{
        Title: r.Form.Get("title"),
        Options: options,
    }

    err := h.store.CreatePoll(poll)
    if err != nil {
        templates.PollError(err).Render(r.Context(), w)
        return
    }

    polls, err := h.store.GetPolls()
    if err != nil {
        templates.PollError(err).Render(r.Context(), w)
        return
    }

	sendPopup(w, "Poll gemaakt")
    templates.Polls(polls).Render(r.Context(), w)
}
