package ytdlp

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"
)

func concurrentFragments() string {
	n := runtime.NumCPU() * 4
	if n < 16 {
		n = 16
	}
	if n > 64 {
		n = 64
	}
	return strconv.Itoa(n)
}

var dlCounter atomic.Int64

func FetchInfo(ctx context.Context, videoURL string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-j",
		"--no-playlist",
		"--skip-download",
		"--no-warnings",
		videoURL,
	)
	return cmd.Output()
}

func Stream(ctx context.Context, videoURL, formatID string, w io.Writer) error {
	id := dlCounter.Add(1)
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("everdl-%d-%d.mp4", os.Getpid(), id))
	defer os.Remove(tmpFile)

	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-f", formatID,
		"--merge-output-format", "mp4",
		"--concurrent-fragments", concurrentFragments(),
		"--no-playlist",
		"--no-mtime",
		"--no-part",
		"--no-warnings",
		"--retries", "3",
		"--fragment-retries", "3",
		"--buffer-size", "16K",
		"--http-chunk-size", "10M",
		"-o", tmpFile,
		videoURL,
	)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp: %w", err)
	}

	f, err := os.Open(tmpFile)
	if err != nil {
		return fmt.Errorf("open merged file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 2*1024*1024)
	_, err = io.CopyBuffer(w, f, buf)
	return err
}

func StreamWithInfo(ctx context.Context, videoURL, formatID string, w io.Writer) (size int64, err error) {
	id := dlCounter.Add(1)
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("everdl-%d-%d.mp4", os.Getpid(), id))
	defer os.Remove(tmpFile)

	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-f", formatID,
		"--merge-output-format", "mp4",
		"--concurrent-fragments", concurrentFragments(),
		"--no-playlist",
		"--no-mtime",
		"--no-part",
		"--no-warnings",
		"--retries", "3",
		"--fragment-retries", "3",
		"--buffer-size", "16K",
		"--http-chunk-size", "10M",
		"-o", tmpFile,
		videoURL,
	)
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return 0, fmt.Errorf("yt-dlp: %w", err)
	}

	f, err := os.Open(tmpFile)
	if err != nil {
		return 0, fmt.Errorf("open merged file: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat merged file: %w", err)
	}
	size = stat.Size()

	buf := make([]byte, 2*1024*1024)
	_, err = io.CopyBuffer(w, f, buf)
	return size, err
}
