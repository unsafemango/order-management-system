package api

import "net/http"

func (a *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]string{"message": "health check passed"})
}
