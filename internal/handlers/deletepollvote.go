package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store/dbstore"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DeletePollVoteHandler struct{
    store *dbstore.PollStore
}

type DeletePollVoteHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewDeletePollVoteHandler(params DeletePollVoteHandlerParams) *DeletePollVoteHandler {
    return &DeletePollVoteHandler{
        store: params.PollStore,
    }
}

func (h *DeletePollVoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)

    user := middleware.GetUser(r.Context())

    if user == nil {
        http.Error(w, "Not authorized", http.StatusUnauthorized)
        return
    }

    err = h.store.DeleteVote(uint(pollId), user.ID)
    if err != nil {
        http.Error(w, "Error deleting vote", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/poll", http.StatusSeeOther)
}
