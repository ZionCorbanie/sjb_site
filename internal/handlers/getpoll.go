package handlers

import (
	"net/http"
	"sjb_site/internal/middleware"
	"sjb_site/internal/store/dbstore"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	pollId, err := strconv.ParseInt(chi.URLParam(r, "pollId"), 10, 64)

    user := middleware.GetUser(r.Context())

    poll, voted := h.store.GetPollVotes(uint(pollId), user.ID)
	if err != nil {
		http.Error(w, "Error getting poll", http.StatusInternalServerError)
		return
	}

    templates.Poll(poll, voted, 2).Render(r.Context(), w)
}
