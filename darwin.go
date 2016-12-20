// +build darwin

// Package wallpaper is UNTESTED on macOS.
package wallpaper

import "os/exec"
import "os/user"
import "path/filepath"
import "strconv"

// SetFromFile uses AppleScript to tell Finder to set the desktop wallpaper to specified file.
func SetFromFile(file string) error {
  return exec.
    Command("osascript", "-e", `tell application "Finder" to set desktop picture to POSIX file ` + strconv.Quote(file)).
    Run()
}

// SetFromURL downloads url and calls SetFromFile.
func SetFromURL(url string) error {
  file, err := downloadImage(url)

  if err != nil {
    return err
  }

  return SetFromFile(file)
}

func getCacheDir() (string, error) {
  usr, err := user.Current()

  if err != nil {
    return "", err
  }

  return filepath.Join(usr.HomeDir, "Library", "Caches"), nil
}
