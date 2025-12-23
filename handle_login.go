package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jaxhemopo/Chirpy/internal/auth"
	"github.com/jaxhemopo/Chirpy/internal/database"
)

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode the parameters", err)
		return
	}
	nullStr := sql.NullString{
		String: params.Email,
		Valid:  true,
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), nullStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get user password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not check password hash", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "invalid credentials", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 days
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not add token to databse", err)
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email.String,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})

}
