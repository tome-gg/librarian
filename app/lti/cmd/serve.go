package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tome-gg/librarian/app/lti/internal/entrypoint/api"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "api-librarian"
	app.Usage = "Launches the Tome.gg Librarian web service"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Usage: "the port number to listen on",
			Value: 8080,
		},
	}

	app.Action = func(c *cli.Context) error {
		port := c.Int("port")
		fmt.Printf("Starting LTI app on port %d...\n", port)
		return api.Start(port)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Failed to start LTI app: %v", err)
	}
}
