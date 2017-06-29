// +build linux

package wallpaper

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-ini/ini"
)

func parseLXDEConfig() (wallpaper string, err error) {
	usr, err := user.Current()

	if err != nil {
		return
	}

	lxdeConfig, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, ".config", "pcmanfm", "LXDE", "desktop-items-0.conf"))

	if err != nil {
		return
	}

	cfg, err := ini.Load(lxdeConfig)

	if err != nil {
		return
	}

	key, err := cfg.
		Section("*").
		GetKey("wallpaper")

	if err != nil {
		return
	}

	wallpaper = key.String()

	return
}

func parseKDEConfig() (wallpaper string, err error) {
	kdeConfig, err := getKDEConfigFile()

	if err != nil {
		return
	}

	file, err := os.Open(kdeConfig)

	if err != nil {
		return
	}

	defer func() {
		err = file.Close()
	}()

	reader := bufio.NewReader(file)

	for {
		var line string

		line, err = reader.ReadString('\n')

		if err != nil {
			return
		}

		if len(line) >= 6 && line[:6] == "Image=" {
			wallpaper = strings.TrimSpace(removeProtocol(line[6:]))

			return
		}
	}
}

func writeKDEConfig(filename string) (err error) {
	// is there a more efficient way of doing this that doesn't require loading the
	// whole file into memory?

	kdeConfig, err := getKDEConfigFile()

	if err != nil {
		return
	}

	config, err := ioutil.ReadFile(kdeConfig)

	if err != nil {
		return
	}

	regex := regexp.MustCompile(`(?m)^Image=.*`)

	err = ioutil.WriteFile(kdeConfig, regex.ReplaceAll(config, []byte("Image="+filename)), 0666)

	if err != nil {
		return
	}

	return
}

func getKDEConfigFile() (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	return filepath.Join(usr.HomeDir, ".config", "plasma-org.kde.plasma.desktop-appletsrc"), nil
}
