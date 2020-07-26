package main

import (
	"os/user"
)

func main() {
	client := NewGUI()

	usr, err := user.Current()
	if err != nil {
		panic("Error obtaining current user")
	}
	client.Start(usr.HomeDir + "/.syndie")
}
