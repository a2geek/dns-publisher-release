package web

import (
	"fmt"
	"net/http"
)

func (s WebServer) getTags(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "TODO")
}
