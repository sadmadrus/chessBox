package game

import (
	"net/http"
)

// Creator — http.HandlerFunc для создания новой игры.
func Creator(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	http.Error(w, "Not implemented.", http.StatusNotImplemented)
}
