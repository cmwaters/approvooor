package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.HideHelp = true
	app.Name = "BlobuSign Backend Server"
	app.Usage = "Manage publishing, requesting and signing documents"
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "Start the backend server",
			Action: func(c *cli.Context) error {
				fmt.Println("Node started")
				// Here you would call the function to start the node
				return nil
			},
		},
		{
			Name:  "display",
			Usage: "Display your public key",
			Action: func(c *cli.Context) error {
				fmt.Println("Showing the key")
				// Here you would retrieve and display the key
				return nil
			},
		},
		{
			Name:  "register",
			Usage: "Register a new private key for signing from a mnemonic",
			Action: func(c *cli.Context) error {
				fmt.Println("Adding a key")
				// Here you would add a key
				return nil
			},
		},
	}

	app.Run(os.Args)
}
