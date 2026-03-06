package main

import (
	"log"
	"net/http"
	"os"

	"EverDownload/internal/cache"
	"EverDownload/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" && os.Getenv("DOCKER_ENV") == "" {
		_ = godotenv.Load()
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	redisCache, err := cache.New(redisURL)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer redisCache.Close()

	h := handlers.New(redisCache)
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	mux.HandleFunc("/api/info", h.VideoInfo)
	mux.HandleFunc("/download", h.Download)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("EverDownload running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
