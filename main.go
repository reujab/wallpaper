package wallpaper

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Desktop contains the current desktop environment on Linux.
// Empty string on all other operating systems.
var Desktop = os.Getenv("XDG_CURRENT_DESKTOP")

var DesktopSession = os.Getenv("DESKTOP_SESSION")

// ErrUnsupportedDE is thrown when Desktop is not a supported desktop environment.
var ErrUnsupportedDE = errors.New("your desktop environment is not supported")

func downloadImage(url string) (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(cacheDir, cleanFilename(filepath.Base(url)))
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non-200 status code")
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return filename, nil
}
