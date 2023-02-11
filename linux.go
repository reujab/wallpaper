// +build linux

package wallpaper

import (
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
)

// Get returns the current wallpaper.
func Get() (string, error) {
	if isGNOMECompliant() {
		return parseDconf("gsettings", "get", "org.gnome.desktop.background", "picture-uri")
	}

	switch Desktop {
	case "KDE":
		return getKDE()
	case "X-Cinnamon":
		return parseDconf("dconf", "read", "/org/cinnamon/desktop/background/picture-uri")
	case "MATE":
		return parseDconf("dconf", "read", "/org/mate/desktop/background/picture-filename")
	case "XFCE":
		return getXFCE()
	case "LXDE":
		return getLXDE()
	case "Deepin":
		return parseDconf("dconf", "read", "/com/deepin/wrap/gnome/desktop/background/picture-uri")
	default:
		return "", ErrUnsupportedDE
	}
}

// SetFromFile sets wallpaper from a file path.
func SetFromFile(file string) error {
	if isGNOMECompliant() {
		err := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", strconv.Quote("file://"+file)).Run()
		if err != nil {
			return err
		}
		// Check if a separate dark mode background URI is available and set it too
		out, err := parseDconf("gsettings", "writable", "org.gnome.desktop.background", "picture-uri-dark")
		if err == nil && out == "true" {
			return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", strconv.Quote("file://"+file)).Run()
		}
		return err
	}

	switch Desktop {
	case "KDE":
		return setKDE(file)
	case "X-Cinnamon":
		return exec.Command("dconf", "write", "/org/cinnamon/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	case "MATE":
		return exec.Command("dconf", "write", "/org/mate/desktop/background/picture-filename", strconv.Quote(file)).Run()
	case "XFCE":
		return setXFCE(file)
	case "LXDE":
		return exec.Command("pcmanfm", "-w", file).Run()
	case "Deepin":
		return exec.Command("dconf", "write", "/com/deepin/wrap/gnome/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	default:
		err := exec.Command("swaybg", "-i", file).Start()
		// if the command completed successfully, return
		if err == nil {
			return nil
		}

		return exec.Command("feh", "-bg-fill", file).Run()
	}
}

// SetMode sets the wallpaper mode.
func SetMode(mode Mode) error {
	if isGNOMECompliant() {
		return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	}

	switch Desktop {
	case "KDE":
		return setKDEMode(mode)
	case "X-Cinnamon":
		return exec.Command("dconf", "write", "/org/cinnamon/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	case "MATE":
		return exec.Command("dconf", "write", "/org/mate/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	case "XFCE":
		return setXFCEMode(mode)
	case "LXDE":
		return exec.Command("pcmanfm", "--wallpaper-mode", mode.getLXDEString()).Run()
	case "Deepin":
		return exec.Command("dconf", "write", "/com/deepin/wrap/gnome/desktop/background/picture-options", strconv.Quote(mode.getGNOMEString())).Run()
	default:
		return ErrUnsupportedDE
	}
}

func getCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".cache"), nil
}
