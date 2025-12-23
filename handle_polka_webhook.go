package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jaxhemopo/Chirpy/internal/database"
)

func (cfg *apiConfig) HandlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID string `json:"user_id"`
	}
	type requestParams struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode the parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parsedId, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not parse user ID", err)
		return
	}

	err = cfg.db.SetIsChirpyRed(r.Context(), database.SetIsChirpyRedParams{
		IsChirpyRed: true,
		ID:          parsedId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not update user to chirpy red", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)

}
