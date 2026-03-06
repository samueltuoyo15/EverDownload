package ytdlp

import (
	"context"
	"io"
	"os/exec"
)

func FetchInfo(ctx context.Context, videoURL string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-j",
		"--no-playlist",
		"--skip-download",
		videoURL,
	)
	return cmd.Output()
}

func Stream(ctx context.Context, videoURL, formatID string, w io.Writer) error {
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-f", formatID,
		"--merge-output-format", "mp4",
		"--prefer-ffmpeg",
		"--no-playlist",
		"--no-mtime",
		"-o", "-",
		videoURL,
	)
	cmd.Stdout = w
	return cmd.Run()
}
