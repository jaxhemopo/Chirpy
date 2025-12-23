package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/jaxhemopo/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db             *database.Queries
	fileserverHits atomic.Int32
	platform       string
	secret         string
	PolkaAPIKey    string
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	pf := os.Getenv("PLATFORM")
	secret := os.Getenv("SERVER_SECRET")
	if secret == "" {
		log.Fatal("SERVER_SECRET is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY is not set")
	}

	const filepathRoot = "."
	const port = "8080"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to open DB connection")
	}
	dbQueries := database.New(db)

	var apiCfg = &apiConfig{
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
		platform:       pf,
		secret:         secret,
		PolkaAPIKey:    polkaKey,
	}

	mux := http.NewServeMux()
	fshandler := http.StripPrefix("/app", http.FileServer((http.Dir(filepathRoot))))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fshandler))
	mux.HandleFunc("GET /api/healthz", HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleMetReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleChirps)
	mux.HandleFunc("POST /api/users", apiCfg.HandleCreateUser)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChirp)
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleRevoke)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandlePolkaWebhook)
	mux.HandleFunc("PUT /api/users", apiCfg.HandleUserCredentials)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandleDeleteChirp)

	Server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server starting on %s", Server.Addr)
	if err := http.ListenAndServe(Server.Addr, Server.Handler); err != nil {
		log.Fatal(err)
	}

}
