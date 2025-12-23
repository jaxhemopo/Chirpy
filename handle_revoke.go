package main

import (
	"net/http"

	"github.com/jaxhemopo/Chirpy/internal/auth"
)

func (cfg *apiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	getRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get bearer token", err)
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), getRefreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
