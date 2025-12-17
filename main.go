package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var apiCfg = &apiConfig{
	fileserverHits: atomic.Int32{},
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func main() {
	mux := http.NewServeMux()
	fshandler := http.StripPrefix("/app", http.FileServer((http.Dir("."))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fshandler))
	mux.HandleFunc("GET /api/healthz", HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", HandleMetrics)
	mux.HandleFunc("POST /admin/reset", HandleMetReset)
	mux.HandleFunc("POST /api/validate_chirp", HandleChirpValidation)

	Server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server starting on %s", Server.Addr)
	if err := http.ListenAndServe(Server.Addr, Server.Handler); err != nil {
		log.Fatal(err)
	}

}

func HandleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := apiCfg.fileserverHits.Load()
	htmlResponse := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Write([]byte(htmlResponse))
}

func HandleMetReset(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reset\n"))
}

func HandleChirpValidation(w http.ResponseWriter, r *http.Request) {

	type Chirp struct {
		Content string `json:"body"`
	}

	type ValidationResponse struct {
		Valid bool `json:"valid"`
	}

	newChirp := Chirp{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newChirp)
	if err != nil {
		log.Printf("Error decoding %s", err)
		w.WriteHeader(500)
		return
	}
	valResp := ValidationResponse{}

	type handleErr struct {
		Message string `json:"error"`
	}
	ErrHandler := handleErr{}

	if len(newChirp.Content) == 0 {
		ErrHandler.Message = "Chirp can not be empty"
		data, err := json.Marshal(ErrHandler)
		if err != nil {
			log.Printf("Error marshalling Json")
			w.WriteHeader(500)
		}
		w.Header().Set("Content-Type:", "application/json")
		w.WriteHeader(500)
		w.Write(data)
		return
	}
	if len(newChirp.Content) > 140 {
		ErrHandler.Message = "Chirp is too long"
		data, err := json.Marshal(ErrHandler)
		if err != nil {
			log.Printf("Error marshalling Json")
			w.WriteHeader(500)
		}
		w.Header().Set("Content-Type:", "application/json")
		w.WriteHeader(400)
		w.Write(data)
		return
	}
	valResp.Valid = true
	data, err := json.Marshal(valResp)
	if err != nil {
		log.Printf("Error marshalling Json")
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type:", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}
