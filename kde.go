//+build linux

package wallpaper

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

func parseKDEConfig() (string, error) {
	filename, err := getKDEConfigFile()
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

func writeKDEConfig(wallpaper string) error {
	filename, err := getKDEConfigFile()
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, regexp.MustCompile(`(?m)^Image=.*`).ReplaceAll(file, []byte("Image="+wallpaper)), 0666)
}

func getKDEConfigFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".config", "plasma-org.kde.plasma.desktop-appletsrc"), nil
}
