package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store/dbstore"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PostPollVoteHandler struct{
    store *dbstore.PollStore
}

type PostPollVoteHandlerParams struct {
    PollStore *dbstore.PollStore
}

func NewPostPollVoteHandler(params PostPollVoteHandlerParams) *PostPollVoteHandler {
    return &PostPollVoteHandler{
        store: params.PollStore,
    }
}

func (h *PostPollVoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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

    err = h.store.Vote(uint(pollId), uint(option), user.ID)
    if err != nil {
        http.Error(w, "Error voting", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/poll", http.StatusSeeOther)
}
