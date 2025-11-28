package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	utils "EverDownload/utils"

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
		FormatID string `json:"format_id"`
		Quality  string `json:"quality"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Ext      string `json:"ext"`
	} `json:"medias"`
	Error bool `json:"error"`
}

type YTDLPOutput struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Uploader   string `json:"uploader"`
	Thumbnail  string `json:"thumbnail"`
	WebpageURL string `json:"webpage_url"`
	Formats    []struct {
		FormatID string `json:"format_id"`
		Ext      string `json:"ext"`
		Format   string `json:"format"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Acodec   string `json:"acodec"`
		Vcodec   string `json:"vcodec"`
		FPS      int    `json:"fps"`
		Filesize int64  `json:"filesize"`
	} `json:"formats"`
}

var ctx = context.Background()
var rdb *redis.Client

var formatIDRegex = regexp.MustCompile(`^[a-zA-Z0-9+_-]+$`)

func isValidFormatID(id string) bool {
	return id != "" && formatIDRegex.MatchString(id)
}

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		_ = godotenv.Load()
	}

	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
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
		if videoURL == "" || !utils.ValidateURL(videoURL) {
			http.Error(w, "Invalid or unsupported video URL", http.StatusBadRequest)
			return
		}

		videoData, err := fetchVideoMetaData(videoURL)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching video meta data: %v", err), http.StatusInternalServerError)
			return
		}

		sanitizedTitle := strings.ReplaceAll(videoData.Title, "/", "-")

		fmt.Fprintf(w, `
			<div class="mt-6 mb-20 p-4 rounded-lg shadow-2xl" x-data="{ selectedFormat: '%s', pageUrl: '%s' }">
			<h3 class="text-lg font-bold mb-4">Video Details</h3>
			<img src="%s" alt="Video Thumbnail" class="w-full rounded-md mb-4" />
			<p class="text-white mb-2"><strong>Title:</strong> %s</p>
			<p class="text-white mb-2"><strong>Author:</strong> %s</p>
			<div class="mt-4">
				<label for="qualitySelect" class="block mb-2">Select Quality</label>
				<select id="qualitySelect" x-model="selectedFormat" class="w-full p-2 bg-neutral-800 text-white rounded-md border">`,
			videoData.Medias[0].FormatID,
			videoData.URL,
			videoData.Thumbnail,
			videoData.Title,
			videoData.Author,
		)

		for _, media := range videoData.Medias {
			label := media.Quality
			if strings.Contains(label, "video only") || strings.Contains(label, "audio only") {
				parts := strings.Split(label, " ")
				res := ""
				for _, p := range parts {
					if strings.HasSuffix(p, "p") || strings.Contains(p, "x") {
						res = p
						break
					}
				}
				if res == "" {
					res = fmt.Sprintf("%dp", media.Height)
				}
				if strings.Contains(media.Quality, "audio only") {
					label = "Audio only"
				} else {
					label = res
				}
			}
			fmt.Fprintf(w, `<option value="%s">%s</option>`, media.FormatID, label)
		}

		fmt.Fprintf(w, `
			</select>
		</div>
		<a 
			x-bind:href="'/download?url=' + encodeURIComponent(pageUrl) + '&filename=%s.mp4&format=' + encodeURIComponent(selectedFormat)" 
			class="block mb-32 w-full mt-4 bg-red-900 text-center text-white p-3 rounded-md hover:bg-blue-600"
			download
		>
			Download Video
		</a>
		</div>`, sanitizedTitle)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		pageURL := r.URL.Query().Get("url")
		if pageURL == "" {
			http.Error(w, "Missing video page URL", http.StatusBadRequest)
			return
		}
		formatID := r.URL.Query().Get("format")
		if !isValidFormatID(formatID) {
			http.Error(w, "Invalid format", http.StatusBadRequest)
			return
		}
		fileName := r.URL.Query().Get("filename")
		if fileName == "" {
			fileName = "video.mp4"
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		w.Header().Set("Content-Type", "video/mp4")

		cmd := exec.Command(
			"yt-dlp",
			"-f", formatID,
			"--merge-output-format", "mp4",
			"--prefer-ffmpeg",
			"--no-mtime",
			"-o", "-",
			pageURL,
		)
		cmd.Stdout = w
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			http.Error(w, "Failed to download video", http.StatusInternalServerError)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fetchVideoMetaData(videoURL string) (*VideoResponse, error) {
	parsedURL, err := url.ParseRequestURI(videoURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, err
	}

	cacheKey := fmt.Sprintf("video_meta:%s", videoURL)
	if cacheData, err := rdb.Get(ctx, cacheKey).Result(); err == nil {
		var v VideoResponse
		if json.Unmarshal([]byte(cacheData), &v) == nil {
			return &v, nil
		}
	}

	cmd := exec.Command("yt-dlp", "-j", videoURL)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ytdlpData YTDLPOutput
	if err := json.Unmarshal(output, &ytdlpData); err != nil {
		return nil, err
	}

	videoResp := &VideoResponse{
		URL:       ytdlpData.WebpageURL,
		ID:        ytdlpData.ID,
		Author:    ytdlpData.Uploader,
		Title:     ytdlpData.Title,
		Thumbnail: ytdlpData.Thumbnail,
	}

	for _, f := range ytdlpData.Formats {
		if f.FormatID == "" {
			continue
		}
		if f.Vcodec == "none" && f.Acodec == "none" {
			continue
		}
		videoResp.Medias = append(videoResp.Medias, struct {
			FormatID string `json:"format_id"`
			Quality  string `json:"quality"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Ext      string `json:"ext"`
		}{
			FormatID: f.FormatID,
			Quality:  f.Format,
			Width:    f.Width,
			Height:   f.Height,
			Ext:      f.Ext,
		})
	}

	cacheData, _ := json.Marshal(videoResp)
	rdb.Set(ctx, cacheKey, cacheData, 5*time.Minute)
	return videoResp, nil
}
