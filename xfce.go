package wallpaper

import (
	"os/exec"
	"path"
	"strings"
)

func getXFCEDesktops() ([]string, error) {
	output, err := exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--list").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.Trim(string(output), "\n"), "\n")
	var desktops []string

	for _, line := range lines {
		if path.Base(line) == "last-image" {
			desktops = append(desktops, line)
		}
	}

	return desktops, nil
}

func getXFCE() (string, error) {
	desktops, err := getXFCEDesktops()
	if err != nil || len(desktops) == 0 {
		return "", err
	}

	output, err := exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--property", desktops[0]).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func setXFCE(file string) error {
	desktops, err := getXFCEDesktops()
	if err != nil {
		return err
	}
	for _, desktop := range desktops {
		err := exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--property", desktop, "--set", file).Run()
		if err != nil {
			return err
		}
	}
	return nil
}
