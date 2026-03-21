package admin

import (
	"net/http"
	"sjb_site/internal/handlers"
	"sjb_site/internal/store"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PollHandler struct{
    store *dbstore.PollStore
}

type PollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewPollHandler(params PollHandlerParams) *PollHandler {
    return &PollHandler{
        store: params.PollStore,
    }
}

func (h *PollHandler) New(w http.ResponseWriter, r *http.Request) {
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

func (h *PollHandler) Create(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    size := len(r.Form)-1

    options := make([]store.PollOption, size)

	for i := range size{
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

	handlers.SendPopup(w, "Poll aangemaakt")
    templates.Polls(polls).Render(r.Context(), w)
}


func (h *PollHandler) Edit(w http.ResponseWriter, r *http.Request) {
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

func (h *PollHandler) Put(w http.ResponseWriter, r *http.Request) {
    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)
    if err != nil {
        templates.PollError(err).Render(r.Context(), w)
        return
    }
    r.ParseForm()
    size := len(r.Form)-1
    options := make([]store.PollOption, size)

	for i := range size{
        options[i] = store.PollOption{
            PollID: uint(pollId),
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

	handlers.SendPopup(w, "Poll aangepast")
    templates.Polls(polls).Render(r.Context(), w)
}

func (h *PollHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pollId := chi.URLParam(r, "pollId")
    
    err := h.store.DeletePoll(pollId)
	if err != nil {
		http.Error(w, "Error deleting polls", http.StatusInternalServerError)
		return
	}

    polls, err := h.store.GetPolls()
    err = templates.Polls(polls).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	handlers.SendPopup(w, "Poll verwijderd")
}

func (h *PollHandler) Activate(w http.ResponseWriter, r *http.Request) {
	pollId := chi.URLParam(r, "pollId")

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
