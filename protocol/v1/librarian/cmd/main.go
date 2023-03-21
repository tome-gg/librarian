package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	librarian "github.com/tome-gg/librarian/protocol/v1/librarian"
	validator "github.com/tome-gg/librarian/protocol/v1/librarian/validator"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tome",
		Usage: "The Tome.gg CLI for working with the Librarian protocol",
		Commands: []*cli.Command{
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "Validate a directory using the Librarian protocol",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Path to the directory to validate",
						DefaultText: "No directory specified. Use --help to see available commands.",
					},
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Enable verbose logging",
					},
				},
				Action: func(c *cli.Context) error {
					directoryPath := c.String("directory")
					verbose := c.Bool("verbose")

					if directoryPath == "" {
						wd, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to get current working directory: %s", err)
						}

						directoryPath = wd
					}

					// Trim trailing slash; limit 1 (will fail with multiple trailing slashes)
					if directoryPath[len(directoryPath) - 1] == '/' {
						directoryPath = directoryPath[:len(directoryPath)-1]
					}

					if verbose {
						logrus.SetLevel(logrus.DebugLevel)
					} else {
						logrus.SetLevel(logrus.WarnLevel)
					}

					logrus.WithFields(logrus.Fields{
						"path": directoryPath,
					}).Infof("validating directory")
					
					directory, err := librarian.Parse(directoryPath)

					if err != nil {
						return fmt.Errorf("failed to parse directory: %s", err)
					}

					plan := validator.Init(directory)
					plan.Init()

					errors := validator.ValidatePlan(plan)

					if len(errors) > 0 {
						logrus.WithFields(
							logrus.Fields{
								"errors": errors,
							}).Error("validation failed")

							return errors[0]
					}

					return nil
				},
			},
		},
	}

	// Set logrus output to stdout
	logrus.SetOutput(os.Stdout)

	// Set logrus timestamp format
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatalf("Error: %s", err)
	}
}
