package handlers

import (
	"net/http"
	"sjb_site/internal/store"
	"sjb_site/internal/templates"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CalendarDayHandler struct{
    calendarStore store.CalendarStore
}

type CalendarDayHandlerParams struct {
    CalendarStore store.CalendarStore
}

func NewCalendarDayHandler(params CalendarDayHandlerParams) *CalendarDayHandler {
	return &CalendarDayHandler{
        calendarStore: params.CalendarStore,
    }
}

func (h *CalendarDayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	day, err := strconv.ParseInt(chi.URLParam(r, "day"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid day parameter", http.StatusBadRequest)
		return
	}
	
    events, err := h.calendarStore.GetCalendarItems(int(day))
    if err != nil {
        http.Error(w, "Error getting day", http.StatusInternalServerError)
        return
    }

	if events == nil || len(*events) == 0 {
		w.WriteHeader(http.StatusOK)
		//w.Write([]byte("No events found"))
		return
	}

    err = templates.CalendarDay(int(day),events).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
