package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jaxhemopo/Chirpy/internal/auth"
	"github.com/jaxhemopo/Chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) HandleChirps(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	newChirp := Chirp{}
	err = decoder.Decode(&newChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode the parameters", err)
		return
	}

	if len(newChirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", err)
		return
	}

	cleanedChirp := HandleProfane(newChirp.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		Body:      cleanedChirp,
	})

	if err != nil {
		log.Println("error creating chirp:", err)
		respondWithError(w, http.StatusInternalServerError, "could not create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func HandleProfane(body string) string {

	badWords := map[string]string{
		"kerfuffle": "kerfuffle",
		"sharbert":  "sharbert",
		"fornax":    "fornax",
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		lowered := strings.ToLower(word)
		if _, exists := badWords[lowered]; exists {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Println("error retrieving chirps:", err)
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirps", err)
		return
	}
	chirpList := []Chirp{}
	for _, chirp := range chirps {
		i := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		chirpList = append(chirpList, i)
	}
	respondWithJSON(w, http.StatusOK, chirpList)
}

func (cfg *apiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	log.Println("Fetching chirp with ID:", idStr)
	chirpID, err := uuid.Parse(idStr)
	if err != nil {
		log.Println("invalid chirp ID format:", err)
		respondWithError(w, http.StatusNotFound, "invalid chirp ID format", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "chirp not found", err)
		} else {
			log.Println("error retrieving chirp:", err)
			respondWithError(w, http.StatusInternalServerError, "could not retrieve chirp", err)
		}
		return
	}

	responseChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, responseChirp)
}
