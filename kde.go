//+build linux

package wallpaper

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func getKDE() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(usr.HomeDir, ".config", "plasma-org.kde.plasma.desktop-appletsrc")
	if err != nil {
		return "", err
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= 6 && line[:6] == "Image=" {
			return strings.TrimSpace(removeProtocol(line[6:])), nil
		}
	}
	if scanner.Err() != nil {
		return "", scanner.Err()
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return "", errors.New("kde image not found")
}

func setKDE(path string) error {
	return evalKDE(`
		for (const desktop of desktops()) {
			desktop.currentConfigGroup = ["Wallpaper", "org.kde.image", "General"]
			desktop.writeConfig("Image", ` + strconv.Quote("file://"+path) + `)
		}
	`)
}

func setKDEMode(mode Mode) error {
	return evalKDE(`
		for (const desktop of desktops()) {
			desktop.currentConfigGroup = ["Wallpaper", "org.kde.image", "General"]
			desktop.writeConfig("FillMode", ` + mode.getKDEString() + `)
		}
	`)
}

func evalKDE(script string) error {
	return exec.Command("qdbus", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", script).Run()
}

func (mode Mode) getKDEString() string {
	switch mode {
	case Center:
		return "6"
	case Crop:
		return "2"
	case Fit:
		return "1"
	case Span:
		return "2"
	case Stretch:
		return "0"
	case Tile:
		return "3"
	default:
		panic("invalid walllpaper mode")
	}
}
