package main

import "net/http"

func (cfg *apiConfig) HandleMetReset(w http.ResponseWriter, r *http.Request) {

	cfg.fileserverHits.Store(0)
	if cfg.platform == "dev" {
		err := cfg.db.ResetUsers(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not reset users", err)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))

}
