package handlers

import (
	"net/http"
)

func SendPopup(w http.ResponseWriter, message string) {
	w.Header().Set("HX-Trigger", `{"message": "` + message + `"}`)
}
