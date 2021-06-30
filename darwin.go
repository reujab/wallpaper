// +build darwin

package wallpaper

import (
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Get returns the path to the current wallpaper.
func Get() (string, error) {
	stdout, err := exec.Command("osascript", "-e", `tell application "Finder" to get POSIX path of (get desktop picture as alias)`).Output()
	if err != nil {
		return "", err
	}

	// is calling strings.TrimSpace() necessary?
	return strings.TrimSpace(string(stdout)), nil
}

// SetFromFile uses AppleScript to tell Finder to set the desktop wallpaper to specified file.
func SetFromFile(file string) error {
	return exec.Command("osascript", "-e", `tell application "System Events" to tell every desktop to set picture to `+strconv.Quote(file)).Run()
}

// SetMode does nothing on macOS.
func SetMode(mode Mode) error {
	return nil
}

func getCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, "Library", "Caches"), nil
}
