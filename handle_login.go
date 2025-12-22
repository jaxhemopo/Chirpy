package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jaxhemopo/Chirpy/internal/auth"
)

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		ExpiresIn int64  `json:"expires_in_seconds"`
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

	expiresIn := time.Duration(params.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = time.Hour * 1
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email.String,
		Token:     token,
	})

}
