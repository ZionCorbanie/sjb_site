package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
)

type GetPollHandler struct{
    store *dbstore.PollStore
}

type GetPollHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewGetPollHandler(params GetPollHandlerParams) *GetPollHandler {
    return &GetPollHandler{
        store: params.PollStore,
    }
}

func (h *GetPollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    active, err := h.store.GetActivePoll()
    if err != nil {
        http.Error(w, "Error getting active poll", http.StatusInternalServerError)
        return
    }

    user := middleware.GetUser(r.Context())
    poll, voted := h.store.GetPollVotes(active.ID, user.ID)
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
