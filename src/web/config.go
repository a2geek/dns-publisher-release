package web

import (
	"encoding/json"
	"net/http"
	"strings"
)

func getConfig(w http.ResponseWriter, r *http.Request) {
	// Convert to bytes...
	b, err := json.Marshal(configuration)
	if err != nil {
		http.NotFound(w, r)
	}
	// ...so we can put it into a map...
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		http.NotFound(w, r)
	}
	// ... do magic...
	replaceSecrets(data)
	// ... and the map into bytes...
	b, err = json.Marshal(data)
	if err != nil {
		http.NotFound(w, r)
	}
	w.Write(b)
}

func replaceSecrets(data map[string]interface{}) {
	for k, v := range data {
		lower := strings.ToLower(k)
		if strings.Contains(lower, "secret") || strings.Contains(lower, "private") {
			v = "[redacted]"
			data[k] = v
		} else {
			m, ok := v.(map[string]interface{})
			if ok {
				replaceSecrets(m)
			}
		}
	}
}
