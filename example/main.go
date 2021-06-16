package main

import (
	"fmt"

	"github.com/reujab/wallpaper"
)

func main() {
	background, err := wallpaper.Get()

	if err != nil {
		panic(err)
	}

	fmt.Println("Current wallpaper:", background)

	err = wallpaper.SetFromFile("/usr/share/backgrounds/gnome/adwaita-day.jpg")
	if err != nil {
		panic(err)
	}

	err = wallpaper.SetFromURL("https://i.imgur.com/pIwrYeM.jpg")
	if err != nil {
		panic(err)
	}
}
