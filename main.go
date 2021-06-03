package main

import (
	"os/user"

	"github.com/kpetku/syndie-gui/ui"
)

func main() {
	client := ui.NewUI()

	usr, err := user.Current()
	if err != nil {
		panic("Error obtaining current user")
	}
	client.Start(usr.HomeDir + "/.syndie")
}
