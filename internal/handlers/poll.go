package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PollHandler struct{
    pollStore *dbstore.PollStore
}

type PollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewPollHandler(params PollHandlerParams) *PollHandler {
    return &PollHandler{
        pollStore: params.PollStore,
    }
}

func (h *PollHandler) Show(w http.ResponseWriter, r *http.Request) {
    active, err := h.pollStore.GetActivePoll()
    if err != nil {
        http.Error(w, "Error getting active poll", http.StatusInternalServerError)
        return
    }

    user := middleware.GetUser(r.Context())
    poll, voted := h.pollStore.GetPollVotes(active.ID, user.ID)
    if poll == nil {
        http.Error(w, "Poll not found", http.StatusNotFound)
        return
    }

    totalVotes := 0
    for _, option := range poll.Options {
        totalVotes += option.VoteCount
    }

    templates.Poll(poll, voted, totalVotes).Render(r.Context(), w)
}

func (h *PollHandler) Vote(w http.ResponseWriter, r *http.Request) {

    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)

    r.ParseForm()
    option, err := strconv.ParseInt(r.Form.Get("option"), 10, 64)
    if err != nil {
        http.Error(w, "Invalid option id", http.StatusBadRequest)
        return
    }

    user := middleware.GetUser(r.Context())

    if user == nil {
        http.Error(w, "Not authorized", http.StatusUnauthorized)
        return
    }

    err = h.pollStore.Vote(uint(pollId), uint(option), user.ID)
    if err != nil {
        http.Error(w, "Error voting", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/poll", http.StatusSeeOther)
}

func (h *PollHandler) RemoveVote(w http.ResponseWriter, r *http.Request) {
    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)

    user := middleware.GetUser(r.Context())

    if user == nil {
        http.Error(w, "Not authorized", http.StatusUnauthorized)
        return
    }

    err = h.pollStore.DeleteVote(uint(pollId), user.ID)
    if err != nil {
        http.Error(w, "Error deleting vote", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/poll", http.StatusSeeOther)
}
