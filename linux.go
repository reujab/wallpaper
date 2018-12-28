// +build linux

package wallpaper

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	ini "gopkg.in/ini.v1"
	yaml "gopkg.in/yaml.v2"
)

// init guesses the current desktop by reading processes if $XDG_CURRENT_DESKTOP was not set.
func init() {
	if Desktop != "" {
		return
	}

	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return
	}

	for _, file := range files {
		// continue if not pid
		_, err := strconv.ParseUint(file.Name(), 10, 64)
		if err != nil {
			continue
		}

		// checks to see if process's binary is named `i3`
		bin, err := os.Readlink("/proc/" + file.Name() + "/exe")
		if err != nil {
			continue
		}

		if path.Base(bin) == "i3" {
			Desktop = "i3"
			return
		}
	}
}

// Get returns the current wallpaper.
func Get() (string, error) {
	if isGNOMECompliant() {
		return parseDconf("gsettings", "get", "org.gnome.desktop.background", "picture-uri")
	}

	switch Desktop {
	case "KDE":
		return parseKDEConfig()
	case "X-Cinnamon":
		return parseDconf("dconf", "read", "/org/cinnamon/desktop/background/picture-uri")
	case "MATE":
		return parseDconf("dconf", "read", "/org/mate/desktop/background/picture-filename")
	case "XFCE":
		output, err := exec.Command("xfconf-query", "-c", "xfce4-desktop", "-p", "/backdrop/screen0/monitor0/workspace0/last-image").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	case "LXDE":
		return parseLXDEConfig()
	case "Deepin":
		return parseDconf("dconf", "read", "/com/deepin/wrap/gnome/desktop/background/picture-uri")
	default:
		return "", ErrUnsupportedDE
	}
}

// SetFromFile sets wallpaper from a file path.
func SetFromFile(file string) error {
	if isGNOMECompliant() {
		return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", strconv.Quote("file://"+file)).Run()
	}

	switch Desktop {
	case "KDE":
		return setKDEBackground("file://" + file)
	case "X-Cinnamon":
		return exec.Command("dconf", "write", "/org/cinnamon/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	case "MATE":
		return exec.Command("dconf", "write", "/org/mate/desktop/background/picture-filename", strconv.Quote(file)).Run()
	case "XFCE":
		return exec.Command("xfconf-query", "-c", "xfce4-desktop", "-p", "/backdrop/screen0/monitor0/workspace0/last-image", "-s", file).Run()
	case "LXDE":
		return exec.Command("pcmanfm", "-w", file).Run()
	case "Deepin":
		return exec.Command("dconf", "write", "/com/deepin/wrap/gnome/desktop/background/picture-uri", strconv.Quote("file://"+file)).Run()
	case "i3":
		return exec.Command("feh", "--bg-fill", file).Run()
	default:
		return ErrUnsupportedDE
	}
}

// SetFromURL sets wallpaper from a URL.
//
// In GNOME, it sets org.gnome.desktop.background.picture-uri to the URL.
// In other desktops, it downloads the image and calls SetFromFile.
func SetFromURL(url string) error {
	switch Desktop {
	// only some GNOME-based desktops support urls for picture-uri
	case "GNOME", "ubuntu:GNOME":
		return exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", strconv.Quote(url)).Run()
	case "i3":
		return exec.Command("feh", "--bg-fill", url).Run()
	default:
		filename, err := downloadImage(url)
		if err != nil {
			return err
		}
		return SetFromFile(filename)
	}
}

func getCacheDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".cache"), nil
}

func removeProtocol(input string) string {
	if len(input) >= 7 && input[:7] == "file://" {
		return input[7:]
	}
	return input
}

func parseDconf(command string, args ...string) (string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}

	// unquote string
	var unquoted string
	// the output is quoted with single quotes, which cannot be unquoted using strconv.Unquote, but it is valid yaml
	err = yaml.UnmarshalStrict(output, &unquoted)
	if err != nil {
		return unquoted, err
	}

	return removeProtocol(unquoted), nil
}

func parseLXDEConfig() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	if DesktopSession == "" {
		DesktopSession = "LXDE"
	}

	cfg, err := ini.Load(filepath.Join(usr.HomeDir, ".config/pcmanfm/"+DesktopSession+"/desktop-items-0.conf"))
	if err != nil {
		return "", err
	}

	key, err := cfg.Section("*").GetKey("wallpaper")
	if err != nil {
		return "", err
	}
	return key.String(), err
}

func isGNOMECompliant() bool {
	return strings.Contains(Desktop, "GNOME") || Desktop == "Unity" || Desktop == "Pantheon"
}
