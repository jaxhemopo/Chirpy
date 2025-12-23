package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jaxhemopo/Chirpy/internal/auth"
)

func (cfg *apiConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if accessToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not get bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid bearer token", err)
		return
	}

	chirpID := r.PathValue("chirpID")

	requestID, err := uuid.Parse(chirpID)

	chirp, err := cfg.db.GetChirp(r.Context(), requestID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you are not allowed to delete this chirp", nil)
		return
	} else {
		err = cfg.db.DeleteChirp(r.Context(), requestID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "chirp not found", err)
				return
			}
			respondWithError(w, http.StatusNotFound, "could not delete chirp", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
