package wallpaper

import (
	"os/exec"
	"path"
	"strings"
)

func getXFCEProps(key string) ([]string, error) {
	output, err := exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--list").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.Trim(string(output), "\n"), "\n")
	var desktops []string

	for _, line := range lines {
		if path.Base(line) == key {
			desktops = append(desktops, line)
		}
	}

	return desktops, nil
}

func getXFCE() (string, error) {
	desktops, err := getXFCEProps("last-image")
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
	desktops, err := getXFCEProps("last-image")
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

func setXFCEMode(mode Mode) error {
	styles, err := getXFCEProps("image-style")
	if err != nil {
		return err
	}

	for _, style := range styles {
		err = exec.Command("xfconf-query", "--channel", "xfce4-desktop", "--property", style, "--set", mode.getXFCEString()).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (mode Mode) getXFCEString() string {
	switch mode {
	case Center:
		return "1"
	case Crop:
		return "5"
	case Fit:
		return "4"
	case Span:
		return "5"
	case Stretch:
		return "3"
	case Tile:
		return "2"
	default:
		panic("invalid wallpaper mode")
	}
}
