package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type PatchGroupMembersHandler struct {
	groupUserStore store.GroupUserStore
	groupStore     store.GroupStore
}

type PatchGroupMembersHandlerParams struct {
	GroupUserStore store.GroupUserStore
	GroupStore     store.GroupStore
}

func NewPatchGroupMembersHandler(params PatchGroupMembersHandlerParams) *PatchGroupMembersHandler {
	return &PatchGroupMembersHandler{
		groupUserStore: params.GroupUserStore,
		groupStore:     params.GroupStore,
	}
}

func (h *PatchGroupMembersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	groupId, _ := strconv.ParseUint(chi.URLParam(r, "groupId"), 10, 64)
	statuses := r.Form["status"]
	ids := r.Form["id"]
	titles := r.Form["title"]
	endDates := r.Form["endDate"]

	for i, status := range statuses {
		id, _ := strconv.ParseUint(ids[i], 10, 64)
		title := titles[i]
		endDate, _ := time.Parse("01-02-2006", endDates[i])
		groupUserPatch := store.GroupUser{
			UserID:  uint(id),
			GroupID: uint(groupId),
			Status:  status,
			Title:   title,
			EndDate: endDate,
		}

		err := h.groupUserStore.UpdateGroupUser(groupUserPatch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	group, err := h.groupStore.GetGroup(chi.URLParam(r, "groupId"))
	if err != nil {
		http.Error(w, "Error getting group", http.StatusInternalServerError)
		return
	}

	groupUsers, err := h.groupUserStore.GetGroupUsersByGroup(chi.URLParam(r, "groupId"))
	if err != nil {
		http.Error(w, "Error getting group users", http.StatusInternalServerError)
		return
	}

	t := templates.UserTable(groupUsers, group)
	err = t.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
