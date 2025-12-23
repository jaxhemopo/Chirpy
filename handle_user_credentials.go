package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jaxhemopo/Chirpy/internal/auth"
	"github.com/jaxhemopo/Chirpy/internal/database"
)

func (cfg *apiConfig) HandleUserCredentials(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if accessToken == "" {
		respondWithError(w, http.StatusUnauthorized, "bearer token is missing", nil)
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

	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode the parameters", err)
		return
	}

	nullString := sql.NullString{
		String: params.Email,
		Valid:  true,
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash new password", err)
		return
	}

	updateParams := database.UpdateUserCredentialsParams{
		ID:       userID,
		Email:    nullString,
		Password: hashedPassword,
	}

	user, err := cfg.db.UpdateUserCredentials(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user credentials", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       nullString.String,
		IsChirpyRed: user.IsChirpyRed,
	})
}
