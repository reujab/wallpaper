// +build windows

package wallpaper

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
const (
	spiSetDeskWallpaper = 0x0014

	uiParam = 0x0000

	spifUpdateINIFile = 0x01
	spifSendChange    = 0x02
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

// Get returns the current wallpaper.
func Get() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\Desktop`, registry.READ)
	if err != nil {
		return "", err
	}
	defer key.Close()

	wallpaper, _, err := key.GetStringValue("Wallpaper")
	if err != nil {
		return "", err
	}

	err = key.Close()
	if err != nil {
		return "", err
	}

	return wallpaper, nil
}

// SetFromFile sets the wallpaper for the current user.
func SetFromFile(filename string) error {
	filenameUTF16, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}

	systemParametersInfo.Call(
		uintptr(spiSetDeskWallpaper),
		uintptr(uiParam),
		uintptr(unsafe.Pointer(filenameUTF16)),
		uintptr(spifUpdateINIFile|spifSendChange),
	)
	return nil
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
