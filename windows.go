// +build windows

package wallpaper

import (
	"os"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724947.aspx
const (
	spiGetDeskWallpaper = 0x0073
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
	// the maximum length of a windows path is 256 utf16 characters
	var filename [256]uint16
	systemParametersInfo.Call(
		uintptr(spiGetDeskWallpaper),
		uintptr(cap(filename)),
		// the memory address of the first byte of the array
		uintptr(unsafe.Pointer(&filename[0])),
		uintptr(0),
	)
	return strings.Trim(string(utf16.Decode(filename[:])), "\x00"), nil
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

var illegalRe = regexp.MustCompile(`[<>:"/\\|?*]`)

// cleanFilename returns s with any illegal filename characters removed.
func cleanFilename(s string) string {
	return illegalRe.ReplaceAllString(s, "")
}
