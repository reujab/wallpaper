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

func setKDE(uri string) error {
	return exec.Command("qdbus", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", `
		const monitors = desktops()
		for (var i = 0; i < monitors.length; i++) {
			monitors[i].currentConfigGroup = ["Wallpaper"]
			monitors[i].writeConfig("Image", `+strconv.Quote(uri)+`)
		}
	`).Run()
}
