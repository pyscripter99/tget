package internals

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func formatBytes(i int64) string {
	var unit string
	var size float64

	switch {
	case i >= 1e9:
		unit = "GB"
		size = float64(i) / 1e9
	case i >= 1e6:
		unit = "MB"
		size = float64(i) / 1e6
	case i >= 1e3:
		unit = "KB"
		size = float64(i) / 1e3
	default:
		unit = "B"
		size = float64(i)
	}

	return fmt.Sprintf("%.2f %s", size, unit)
}

func Download(url string, output string, userAgent string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Check if user agent is provided
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", "tget")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if output == "" {
		output = path.Base(req.URL.Path)
		output = strings.Replace(output, "/", "", 1)
		output = strings.ReplaceAll(output, "/", "-")
	}

	if output == "." || output == "/" {
		Warn(fmt.Errorf("invalid output path: '%s', using: '%s' instead", output, req.URL.Host))
		output = req.URL.Host
	}

	// Print info similar to wget, e.g.: file name, size, destination, headers
	fmt.Printf("Saving to: %s\n", output)
	fmt.Printf("File size: %d bytes (%s)\n", resp.ContentLength, formatBytes(resp.ContentLength))
	fmt.Printf("User Agent: %s\n", req.Header.Get("User-Agent"))

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)

	io.Copy(io.MultiWriter(f, bar), resp.Body)

	if resp.ContentLength == -1 {
		fmt.Print("\n") // Add newline after progress bar, as it is not added automatically
	}

	return nil
}
