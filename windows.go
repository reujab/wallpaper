// +build windows

package wallpaper

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

// Get gets the current wallpaper.
func Get() (wallpaper string, err error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.READ)

	if err != nil {
		return
	}

	defer func() {
		err = key.Close()
	}()

	wallpaper, _, err = key.GetStringValue("Wallpaper")

	if err != nil {
		return
	}

	return
}

// SetFromFile sets the wallpaper for the current user to specified file by setting HKEY_CURRENT_USER\Control Panel\Desktop\Wallpaper.
//
// Note: this requires you to log out and in again.
func SetFromFile(file string) (err error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.WRITE)

	if err != nil {
		return
	}

	defer func() {
		err = key.Close()
	}()

	err = key.SetStringValue("Wallpaper", file)

	if err != nil {
		return
	}

	// this is supposed to update the wallpaper, but i only got a black background
	// err = exec.
	//   Command("rundll32", "user32.dll,UpdatePerUserSystemParameters").
	//   Run()
	//
	// if err != nil {
	//   return
	// }

	return
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
	return os.TempDir(), nil
}
