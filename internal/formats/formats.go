package formats

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type rawFormat struct {
	FormatID       string  `json:"format_id"`
	Ext            string  `json:"ext"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	Acodec         string  `json:"acodec"`
	Vcodec         string  `json:"vcodec"`
	ABR            float64 `json:"abr"`
	TBR            float64 `json:"tbr"`
	Filesize       int64   `json:"filesize"`
	FilesizeApprox int64   `json:"filesize_approx"`
}

type rawInfo struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Uploader   string      `json:"uploader"`
	Thumbnail  string      `json:"thumbnail"`
	WebpageURL string      `json:"webpage_url"`
	Formats    []rawFormat `json:"formats"`
}

type Format struct {
	FormatID       string `json:"format_id"`
	Label          string `json:"label"`
	Ext            string `json:"ext"`
	Height         int    `json:"height"`
	Filesize       int64  `json:"filesize"`
	FilesizeApprox bool   `json:"filesize_approx"`
	HasVideo       bool   `json:"has_video"`
	HasAudio       bool   `json:"has_audio"`
}

type VideoInfo struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Author    string   `json:"author"`
	Thumbnail string   `json:"thumbnail"`
	URL       string   `json:"url"`
	Formats   []Format `json:"formats"`
}

func hasVideo(f rawFormat) bool { return f.Vcodec != "" && f.Vcodec != "none" }
func hasAudio(f rawFormat) bool { return f.Acodec != "" && f.Acodec != "none" }

func extFamily(ext string) string {
	switch ext {
	case "mp4", "m4a", "m4v":
		return "mp4"
	case "webm":
		return "webm"
	default:
		return ext
	}
}

func effectiveSize(filesize, approx int64) (int64, bool) {
	if filesize > 0 {
		return filesize, false
	}
	if approx > 0 {
		return approx, true
	}
	return 0, true
}

func fmtSize(bytes int64, approx bool) string {
	prefix := ""
	if approx {
		prefix = "~"
	}
	const mb = 1024 * 1024
	const gb = 1024 * mb
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%s%.1f GB", prefix, float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%s%.0f MB", prefix, float64(bytes)/float64(mb))
	default:
		return fmt.Sprintf("%s%.0f KB", prefix, float64(bytes)/1024)
	}
}

func ParseInfo(data []byte) (*VideoInfo, error) {
	var raw rawInfo
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var videoOnly, audioOnly, combined []rawFormat
	for _, f := range raw.Formats {
		hv, ha := hasVideo(f), hasAudio(f)
		if !hv && !ha {
			continue
		}
		switch {
		case hv && ha:
			combined = append(combined, f)
		case hv:
			videoOnly = append(videoOnly, f)
		case ha:
			audioOnly = append(audioOnly, f)
		}
	}

	bestAudio := map[string]rawFormat{}
	for _, a := range audioOnly {
		fam := extFamily(a.Ext)
		existing, ok := bestAudio[fam]
		if !ok || a.ABR > existing.ABR || (a.ABR == existing.ABR && a.TBR > existing.TBR) {
			bestAudio[fam] = a
		}
	}

	seen := map[string]bool{}
	var formats []Format

	for _, f := range combined {
		key := fmt.Sprintf("%d:%s", f.Height, extFamily(f.Ext))
		if seen[key] {
			continue
		}
		seen[key] = true
		sz, approx := effectiveSize(f.Filesize, f.FilesizeApprox)
		label := fmt.Sprintf("%dp %s", f.Height, strings.ToUpper(f.Ext))
		if sz > 0 {
			label += " " + fmtSize(sz, approx)
		}
		formats = append(formats, Format{
			FormatID: f.FormatID, Label: label, Ext: f.Ext,
			Height: f.Height, Filesize: sz, FilesizeApprox: approx,
			HasVideo: true, HasAudio: true,
		})
	}

	for _, f := range videoOnly {
		key := fmt.Sprintf("%d:%s", f.Height, extFamily(f.Ext))
		if seen[key] {
			continue
		}
		seen[key] = true

		mergedID := f.FormatID
		outExt := f.Ext
		var sz int64
		approx := true

		for _, pref := range []string{"mp4", "webm"} {
			if a, ok := bestAudio[pref]; ok {
				mergedID = f.FormatID + "+" + a.FormatID
				vsz, _ := effectiveSize(f.Filesize, f.FilesizeApprox)
				asz, _ := effectiveSize(a.Filesize, a.FilesizeApprox)
				sz = vsz + asz
				approx = f.FilesizeApprox > 0 || a.FilesizeApprox > 0
				outExt = "mp4"
				break
			}
		}

		label := fmt.Sprintf("%dp %s", f.Height, strings.ToUpper(outExt))
		if sz > 0 {
			label += " " + fmtSize(sz, approx)
		}
		formats = append(formats, Format{
			FormatID: mergedID, Label: label, Ext: outExt,
			Height: f.Height, Filesize: sz, FilesizeApprox: approx,
			HasVideo: true, HasAudio: len(bestAudio) > 0,
		})
	}

	added := 0
	for _, pref := range []string{"mp4", "webm", "mp3"} {
		if added >= 2 {
			break
		}
		a, ok := bestAudio[pref]
		if !ok {
			continue
		}
		sz, approxSz := effectiveSize(a.Filesize, a.FilesizeApprox)
		abrStr := ""
		if a.ABR > 0 {
			abrStr = fmt.Sprintf(" %.0fk", a.ABR)
		}
		label := fmt.Sprintf("Audio only (%s%s)", strings.ToUpper(a.Ext), abrStr)
		if sz > 0 {
			label += " " + fmtSize(sz, approxSz)
		}
		formats = append(formats, Format{
			FormatID: a.FormatID, Label: label, Ext: a.Ext,
			Filesize: sz, FilesizeApprox: approxSz, HasAudio: true,
		})
		added++
	}

	sort.Slice(formats, func(i, j int) bool {
		return formats[i].Height > formats[j].Height
	})

	return &VideoInfo{
		ID:        raw.ID,
		Title:     raw.Title,
		Author:    raw.Uploader,
		Thumbnail: raw.Thumbnail,
		URL:       raw.WebpageURL,
		Formats:   formats,
	}, nil
}
