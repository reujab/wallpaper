# wallpaper [![Documentation](https://godoc.org/github.com/reujab/wallpaper?status.svg)](https://godoc.org/github.com/reujab/wallpaper)

A cross-platform (Linux, Windows, and macOS) Golang library for getting and setting the desktop background.

## Installation

```sh
go get github.com/reujab/wallpaper
```

## Example

```go
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
	wallpaper.SetFromFile("/usr/share/backgrounds/gnome/adwaita-day.jpg")
	wallpaper.SetFromURL("https://i.imgur.com/pIwrYeM.jpg")
}
```

## Notes

* Enlightenment is not supported.
