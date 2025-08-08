package main

import (
	utils "EverDownload/utils"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type VideoRequest struct {
	URL string `json:"url"`
}
type VideoResponse struct {
	URL       string `json:"url"`
	Source    string `json:"source"`
	ID        string `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Medias    []struct {
		URL     string `json:"url"`
		Quality string `json:"quality"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Ext     string `json:"ext"`
	} `json:"medias"`
	Error bool `json:"error"`
}

var ctx = context.Background()
var rdb *redis.Client

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		_ = godotenv.Load()
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable not set")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_URL"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		DB:        0,
		TLSConfig: &tls.Config{},
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error Parsing Form", http.StatusBadRequest)
			return
		}

		videoURL := r.FormValue("videoURL")
		if videoURL == "" {
			http.Error(w, "videoURL is required", http.StatusBadRequest)
			return
		}

		videoData, err := fetchVideoMetaData(videoURL, secretKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching video meta data: %v", err), http.StatusInternalServerError)
			return
		}

		sanitizedTitle := strings.ReplaceAll(videoData.Title, "/", "-")

		fmt.Fprintf(w, `
			<div class="mt-6 mb-20 p-4 rounded-lg shadow-2xl" x-data="{ selectedUrl: '%s' }">
			<h3 class="text-lg font-bold mb-4">Video Details</h3>
			<img src="%s" alt="Video Thumbnail" class="w-full rounded-md mb-4" />
			<p class="text-white mb-2"><strong>Title:</strong> %s</p>
			<p class="text-white mb-2"><strong>Source:</strong> %s</p>
			<p class="text-white mb-2"><strong>Author:</strong> %s</p>
			<div class="mt-4">
				<label for="qualitySelect" class="block mb-2">Select Quality</label>
				<select id="qualitySelect" x-model="selectedUrl" class="w-full p-2 bg-neutral-800 text-white rounded-md border">`,
			videoData.Medias[0].URL,
			videoData.Thumbnail,
			videoData.Title,
			videoData.Source,
			videoData.Author,
		)

		for _, media := range videoData.Medias {
			qualityLabel := strings.TrimSpace(media.Quality)
			if qualityLabel == "" {
				qualityLabel = fmt.Sprintf("%dx%d %s", media.Width, media.Height, media.Ext)
			}
			fmt.Fprintf(w, `<option value="%s">%s</option>`, media.URL, qualityLabel)
		}

		fmt.Fprintf(w, `
			</select>
		</div>
		<a 
			x-bind:href="'/download?url=' + encodeURIComponent(selectedUrl) + '&filename=%s.mp4'" 
			class="block mb-32 w-full mt-4 bg-red-900 text-center text-white p-3 rounded-md hover:bg-blue-600"
			download
		>
			Download Video
		</a>
		</div>`, sanitizedTitle)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "URL parameter is required", http.StatusBadRequest)
			return
		}
		fileName := r.URL.Query().Get("filename")
		if fileName == "" {
			fileName = "video.mp4"
		}

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Failed to fetch video", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
		io.Copy(w, resp.Body)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fetchVideoMetaData(videoURL, apiKey string) (*VideoResponse, error) {
	parsedURL, err := url.ParseRequestURI(videoURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, errors.New("invalid video URL")
	}

	host := strings.ToLower(parsedURL.Hostname())
	if !utils.AllowedHosts[host] {
		return nil, errors.New("Unsupported video source")
	}

	cacheKey := fmt.Sprintf("video_meta:%s", videoURL)
	if cacheData, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		var v VideoResponse
		if json.Unmarshal([]byte(cacheData), &v) == nil {
			return &v, nil
		}
	}
	payload := VideoRequest{URL: videoURL}
	API_HOST := os.Getenv("HOST")
	API_ENDPOINT := os.Getenv("API_ENDPOINT")
	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", API_ENDPOINT, bytes.NewBuffer(jsonPayload))
	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", API_HOST)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result VideoResponse
	if json.Unmarshal(body, &result) != nil || result.Error {
		return &result, nil
	}
	cachedData, _ := json.Marshal(result)
	rdb.Set(ctx, cacheKey, cachedData, 5*time.Minute)
	return &result, nil
}

func getIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
