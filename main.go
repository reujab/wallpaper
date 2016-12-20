package wallpaper

import "errors"
import "io"
import "net/http"
import "os"
import "path/filepath"

// Desktop contains the current desktop environment on Linux.
// Empty string on all other operating systems.
var Desktop = os.Getenv("XDG_CURRENT_DESKTOP")
// ErrUnsupportedDE is thrown when Desktop is not a supported desktop environment.
var ErrUnsupportedDE = errors.New("your desktop environment is not supported")

func downloadImage(url string) (filename string, err error) {
  cacheDir, err := getCacheDir()

  if err != nil {
    return
  }

  filename = filepath.Join(cacheDir, "bing-background.jpg")

  file, err := os.Create(filename)

  if err != nil {
    return
  }

  defer func() {
    err = file.Close()
  }()

  res, err := http.Get(url)

  if err != nil {
    return
  }

  defer func() {
    err = res.
      Body.
      Close()
  }()

  _, err = io.Copy(file, res.Body)

  if err != nil {
    return
  }

  return
}
