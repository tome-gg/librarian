package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	librarian "github.com/tome-gg/librarian/protocol/v1/librarian"
	validator "github.com/tome-gg/librarian/protocol/v1/librarian/validator"

	"os/exec"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tome",
		Usage: "The Tome.gg CLI for working with the Librarian protocol",
		Commands: []*cli.Command{
			{
				Name: "initalize",
				Aliases: []string{"init"},
				Usage: "Initializes a new Git repository using the Tome.gg template, and then immediately clones it into a target directory.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "name",
						Aliases: []string{"n"},
						Usage: "The name of the repository to be generated using the gh CLI tool.",
						Required: true,
					},
					&cli.StringFlag{
						Name: "destination",
						Aliases: []string{"dest"},
						Usage: "The directory path designating the target destination where the repository should be cloned locally.",
						Required: true,
					},
					&cli.BoolFlag{
						Name: "public",
						Usage: "Initializes the GitHub repository as public. Defaults to a private repository.",
					},
				},
				Action: func(ctx *cli.Context) error {
					repositoryName := ctx.String("name")
					directory := ctx.String("directory")
					isPublic := ctx.Bool("public")

					publicFlag := "--private"
					if repositoryName == "" {
						return fmt.Errorf("Invalid repository name")
					}

					if isPublic {
						publicFlag = "--public"
					}

					cmd := exec.Command("gh", "repo", "create", repositoryName, "--template", "tome-gg/template",  publicFlag)
					err := cmd.Run()
					if err != nil {
						logrus.Errorf("Initialize repository failed: %s", err)
						return err
					}

					cmd = exec.Command("gh", "repo", "clone", repositoryName, directory)
					err = cmd.Run()
					if err != nil {
						logrus.Errorf("Clone repository failed: %s", err)
						return err
					}

					return nil
				},
			},
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

					fmt.Printf(" 🚀 [SUCCESS] Repository %s is valid!\n", directoryPath)

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
