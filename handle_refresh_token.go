package main

import (
	"net/http"

	"github.com/jaxhemopo/Chirpy/internal/auth"
)

func (cfg *apiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	type responseParams struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not get bearer token", err)
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}
	newToken, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, responseParams{
		Token: newToken,
	})

}
