package web

import (
	"encoding/json"
	"net/http"
)

func (s WebServer) getLogs(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(s.logCache.Lines())
	if err != nil {
		http.NotFound(w, r)
	}
	w.Write(b)
}
