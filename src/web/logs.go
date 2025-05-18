package web

import (
	"fmt"
	"net/http"
)

func getLogs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "TODO")
}
