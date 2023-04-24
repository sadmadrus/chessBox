package game

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Creator — http.HandlerFunc для создания новой игры.
func Creator(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	data, err := parseUrlEncoded(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	manager := data.Get("notify")
	white := data.Get("white")
	black := data.Get("black")
	if manager == "" || white == "" || black == "" {
		http.Error(w, "Required parameter missing.", http.StatusBadRequest)
		return
	}

	for _, p := range []string{"position", "timing", "move1", "timewhite", "timeblack"} {
		if _, ok := data[p]; ok {
			http.Error(w, "Not (yet) implemented.", http.StatusNotImplemented)
			return
		}
	}

	g, err := new(manager, white, black)
	if err != nil {
		// TODO проверка, не вернуть ли 408
		http.Error(w, fmt.Sprintf("Couldn't create the game: %v", err), http.StatusInternalServerError)
	}

	err = g.registerAndServe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't start the game: %v", err), http.StatusInternalServerError)
	}

	w.Header().Add("location", "/"+string(g.id))
	w.WriteHeader(http.StatusCreated)
	g.serveCurrentState(w)
}

// parseUrlEncoded возвращает данные из www-url-encoded.
func parseUrlEncoded(r *http.Request) (url.Values, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read the request's body: %w", err)
	}
	data, err := url.ParseQuery(string(b))
	if err != nil {
		return data, fmt.Errorf("could not parse data: %w", err)
	}
	return data, nil
}
