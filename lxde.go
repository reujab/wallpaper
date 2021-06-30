package wallpaper

import (
	"os/user"
	"path/filepath"

	ini "gopkg.in/ini.v1"
)

func getLXDE() (string, error) {
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

func (mode Mode) getLXDEString() string {
	switch mode {
	case Center:
		return "center"
	case Crop:
		return "crop"
	case Fit:
		return "fit"
	case Span:
		return "screen"
	case Stretch:
		return "stretch"
	case Tile:
		return "tile"
	default:
		panic("invalid wallpaper mode")
	}
}
