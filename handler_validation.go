package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func HandleChirpValidation(w http.ResponseWriter, r *http.Request) {

	type Chirp struct {
		Content string `json:"body"`
	}

	type returnVal struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	newChirp := Chirp{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode the parameters", err)
		return
	}

	if len(newChirp.Content) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp too long", err)
		return
	}

	cleanedChirp := HandleProfane(newChirp.Content)

	respondWithJSON(w, http.StatusOK, returnVal{
		CleanBody: cleanedChirp,
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
