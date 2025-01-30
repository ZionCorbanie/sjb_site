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
	voorzitter, _ := strconv.ParseUint(r.FormValue("voorzitter"), 10, 64)
	secretaris, _ := strconv.ParseUint(r.FormValue("secretaris"), 10, 64)
	penningmeester, _ := strconv.ParseUint(r.FormValue("penningmeester"), 10, 64)

	for i, status := range statuses {
		id, _ := strconv.ParseUint(ids[i], 10, 64)
		title := titles[i]
		endDate, _ := time.Parse("2006-01-02", endDates[i])

		// Als end date is gezet is user meteen oud-lid, dus als oud-lid wordt geselecteerd wordt ook meteen end date ingevuld
		if statuses[i] != "oud_lid" {
			endDate = time.Time{}
		} else if !endDate.IsZero() {
			status = "oud_lid"
		} else {
			endDate = time.Now()
		}

		//Check of user voorzi, secri of penni is gemaakt
		function := ""
		if voorzitter == id {
			function = "voorzitter"
		} else if secretaris == id {
			function = "secretaris"
		} else if penningmeester == id {
			function = "penningmeester"
		}
		groupUserPatch := store.GroupUser{
			UserID:   uint(id),
			GroupID:  uint(groupId),
			Status:   status,
			Title:    title,
			EndDate:  &endDate,
			Function: function,
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
