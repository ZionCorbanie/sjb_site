package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"

	"github.com/go-chi/chi/v5"
)

type CalendarPopupHandler struct{
    calendarStore store.CalendarStore
}

type CalendarPopupHandlerParams struct {
    CalendarStore store.CalendarStore
}

func NewCalendarPopupHandler(params CalendarPopupHandlerParams) *CalendarPopupHandler {
	return &CalendarPopupHandler{
        calendarStore: params.CalendarStore,
    }
}

func (h *CalendarPopupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "eventId")
	
    event, err := h.calendarStore.GetCalendarItem(id)
    if err != nil {
        http.Error(w, "Error getting event", http.StatusInternalServerError)
        return
    }

    err = templates.AgendaPopup(event).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
