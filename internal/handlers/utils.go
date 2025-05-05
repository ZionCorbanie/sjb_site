package handlers

import (
	"net/http"
)

func sendPopup(w http.ResponseWriter, message string) {
	w.Header().Set("HX-Trigger", `{"message": "` + message + `"}`)
}
