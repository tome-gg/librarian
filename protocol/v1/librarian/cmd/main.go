package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	librarian "github.com/tome-gg/librarian/protocol/v1/librarian"
	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
	validator "github.com/tome-gg/librarian/protocol/v1/librarian/validator"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf(" 📚 Tome.gg CLI; 🚀 version %s\n 🌎 Source: https://github.com/tome-gg/librarian\n 💜 Dreams of sustainability and freedom built from Manila\n\n", cCtx.App.Version)
	}
	app := &cli.App{
		Name:  "tome",
		Version: "0.4.5",
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
				Name:    "missing-evaluations",
				Aliases: []string{"missing"},
				Usage:   "Find DSU entries that don't have corresponding self evaluations",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Path to the directory to analyze",
						DefaultText: "No directory specified. Use --help to see available commands.",
					},
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Show all missing evaluations (default: show last 3 only)",
					},
				},
				Action: func(c *cli.Context) error {
					directoryPath := c.String("directory")
					showAll := c.Bool("all")

					if directoryPath == "" {
						wd, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to get current working directory: %s", err)
						}
						directoryPath = wd
					}

					// Trim trailing slash
					if directoryPath[len(directoryPath) - 1] == '/' {
						directoryPath = directoryPath[:len(directoryPath)-1]
					}

					directory, err := librarian.Parse(directoryPath)
					if err != nil {
						return fmt.Errorf("failed to parse directory: %s", err)
					}

					plan := validator.Init(directory)
					plan.Init()

					missingEvaluations, err := validator.FindMissingEvaluations(plan, !showAll)
					if err != nil {
						return fmt.Errorf("failed to find missing evaluations: %s", err)
					}

					if len(missingEvaluations) == 0 {
						fmt.Println("✅ All DSU entries have corresponding self evaluations!")
						return nil
					}

					if showAll {
						fmt.Printf("Found %d DSU entries without self evaluations (all entries):\n\n", len(missingEvaluations))
					} else {
						fmt.Printf("Found %d DSU entries without self evaluations (last 3, use --all for complete list):\n\n", len(missingEvaluations))
					}
					for _, entry := range missingEvaluations {
						fmt.Printf("UUID: %s\nDate: %s\n\n", entry.ID, entry.Datetime.Format("2006-01-02 15:04:05"))
					}

					return nil
				},
			},
			{
				Name:    "get-dsu",
				Aliases: []string{"get"},
				Usage:   "Retrieve a DSU entry by its UUID",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Path to the directory to search",
						DefaultText: "No directory specified. Use --help to see available commands.",
					},
					&cli.StringFlag{
						Name:     "uuid",
						Aliases:  []string{"u"},
						Usage:    "UUID of the DSU entry to retrieve",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					directoryPath := c.String("directory")
					uuid := c.String("uuid")

					if directoryPath == "" {
						wd, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to get current working directory: %s", err)
						}
						directoryPath = wd
					}

					// Trim trailing slash
					if directoryPath[len(directoryPath) - 1] == '/' {
						directoryPath = directoryPath[:len(directoryPath)-1]
					}

					directory, err := librarian.Parse(directoryPath)
					if err != nil {
						return fmt.Errorf("failed to parse directory: %s", err)
					}

					plan := validator.Init(directory)
					plan.Init()

					entry, err := validator.GetDSUByUUID(plan, uuid)
					if err != nil {
						return fmt.Errorf("failed to get DSU entry: %s", err)
					}

					fmt.Printf("UUID: %s\n", entry.ID)
					fmt.Printf("Date: %s\n", entry.Datetime.Format("2006-01-02 15:04:05"))
					fmt.Printf("Done Yesterday: %s\n", entry.DoneYesterday)
					fmt.Printf("Doing Today: %s\n", entry.DoingToday)
					if entry.Blockers != "" {
						fmt.Printf("Blockers: %s\n", entry.Blockers)
					}
					if entry.Remarks != "" {
						fmt.Printf("Remarks: %s\n", entry.Remarks)
					}

					return nil
				},
			},
			{
				Name:    "get-latest",
				Aliases: []string{"latest"},
				Usage:   "Retrieve the most recent DSU entry by date",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Path to the directory to search",
						DefaultText: "No directory specified. Use --help to see available commands.",
					},
				},
				Action: func(c *cli.Context) error {
					directoryPath := c.String("directory")

					if directoryPath == "" {
						wd, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to get current working directory: %s", err)
						}
						directoryPath = wd
					}

					// Trim trailing slash
					if directoryPath[len(directoryPath) - 1] == '/' {
						directoryPath = directoryPath[:len(directoryPath)-1]
					}

					directory, err := librarian.Parse(directoryPath)
					if err != nil {
						return fmt.Errorf("failed to parse directory: %s", err)
					}

					plan := validator.Init(directory)
					plan.Init()

					entry, err := validator.GetLatestDSU(plan)
					if err != nil {
						return fmt.Errorf("failed to get latest DSU entry: %s", err)
					}

					fmt.Printf("🚀 Latest DSU Entry:\n\n")
					fmt.Printf("UUID: %s\n", entry.ID)
					fmt.Printf("Date: %s\n", entry.Datetime.Format("2006-01-02 15:04:05"))
					fmt.Printf("Done Yesterday: %s\n", entry.DoneYesterday)
					fmt.Printf("Doing Today: %s\n", entry.DoingToday)
					if entry.Blockers != "" {
						fmt.Printf("Blockers: %s\n", entry.Blockers)
					}
					if entry.Remarks != "" {
						fmt.Printf("Remarks: %s\n", entry.Remarks)
					}

					return nil
				},
			},
			{
				Name:    "completion",
				Usage:   "Generate shell completion scripts",
				Subcommands: []*cli.Command{
					{
						Name:  "fish",
						Usage: "Generate fish completion script",
						Action: func(c *cli.Context) error {
							fmt.Println(`# Fish completion for tome
complete -c tome -f

# Main commands
complete -c tome -n "__fish_use_subcommand" -a "init" -d "Initialize a new Git repository using the Tome.gg template"
complete -c tome -n "__fish_use_subcommand" -a "initalize" -d "Initialize a new Git repository using the Tome.gg template"
complete -c tome -n "__fish_use_subcommand" -a "missing-evaluations" -d "Find DSU entries that don't have corresponding self evaluations"
complete -c tome -n "__fish_use_subcommand" -a "missing" -d "Find DSU entries that don't have corresponding self evaluations"
complete -c tome -n "__fish_use_subcommand" -a "get-dsu" -d "Retrieve a DSU entry by its UUID"
complete -c tome -n "__fish_use_subcommand" -a "get" -d "Retrieve a DSU entry by its UUID"
complete -c tome -n "__fish_use_subcommand" -a "get-latest" -d "Retrieve the most recent DSU entry by date"
complete -c tome -n "__fish_use_subcommand" -a "latest" -d "Retrieve the most recent DSU entry by date"
complete -c tome -n "__fish_use_subcommand" -a "validate" -d "Validate a directory using the Librarian protocol"
complete -c tome -n "__fish_use_subcommand" -a "completion" -d "Generate shell completion scripts"
complete -c tome -n "__fish_use_subcommand" -a "help" -d "Shows a list of commands or help for one command"

# Global options
complete -c tome -l help -s h -d "Show help"
complete -c tome -l version -s v -d "Print the version"

# Directory flag for commands that support it
complete -c tome -n "__fish_seen_subcommand_from missing-evaluations missing get-dsu get get-latest latest validate" -l directory -s d -d "Path to the directory" -r

# Missing evaluations flags
complete -c tome -n "__fish_seen_subcommand_from missing-evaluations missing" -l all -d "Show all missing evaluations (default: show last 3 only)"

# UUID flag for get-dsu command
complete -c tome -n "__fish_seen_subcommand_from get-dsu get" -l uuid -s u -d "UUID of the DSU entry to retrieve" -r

# Init command flags
complete -c tome -n "__fish_seen_subcommand_from init initalize" -l name -s n -d "The name of the repository" -r
complete -c tome -n "__fish_seen_subcommand_from init initalize" -l destination -d "The directory path for cloning" -r
complete -c tome -n "__fish_seen_subcommand_from init initalize" -l dest -d "The directory path for cloning" -r
complete -c tome -n "__fish_seen_subcommand_from init initalize" -l public -d "Initialize as public repository"

# Validate command flags
complete -c tome -n "__fish_seen_subcommand_from validate" -l verbose -d "Enable verbose logging"

# Completion subcommands
complete -c tome -n "__fish_seen_subcommand_from completion" -a "fish" -d "Generate fish completion script"`)
							return nil
						},
					},
				},
			},
			{
				Name:    "validate",
				Aliases: []string{},
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

					logrus.WithFields(logrus.Fields{
						"plan": plan,
					}).Debug("plan initialized")

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
			{
				Name:    "dimensions",
				Usage:   "Display all evaluation dimensions with their aliases, names, and labels",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "directory",
						Aliases:     []string{"d"},
						Usage:       "Target directory to scan for evaluation files",
						DefaultText: "current directory",
					},
				},
				Action: func(cCtx *cli.Context) error {
					directoryPath := cCtx.String("directory")
					if directoryPath == "" {
						currentDir, err := os.Getwd()
						if err != nil {
							return fmt.Errorf("failed to get current directory: %s", err)
						}
						directoryPath = currentDir
					}

					directory, err := librarian.Parse(directoryPath)
					if err != nil {
						return fmt.Errorf("failed to parse directory %s: %s", directoryPath, err)
					}

					plan := validator.Init(directory)
					plan.Init()

					// Find evaluation files and extract dimensions
					dimensions := make(map[string]struct {
						Alias      string
						Name       string
						Version    string
						Definition string
					})

					for _, file := range plan.Files {
						if !strings.Contains(file.Filepath, "evaluations") || !strings.HasSuffix(file.Filepath, ".yaml") {
							continue
						}

						fileBytes, err := os.ReadFile(file.Filepath)
						if err != nil {
							continue
						}

						var result pkg.EvaluationDefinition[pkg.StandardMeasurement]
						err = yaml.Unmarshal(fileBytes, &result)
						if err != nil {
							continue
						}

						if result.Tomegg.Type != "evaluations" {
							continue
						}

						for _, dim := range result.Meta.Dimensions {
							dimensions[dim.Alias] = struct {
								Alias      string
								Name       string
								Version    string
								Definition string
							}{
								Alias:      dim.Alias,
								Name:       dim.Name,
								Version:    dim.Version,
								Definition: dim.Definition,
							}
						}
					}

					if len(dimensions) == 0 {
						fmt.Println("No evaluation dimensions found in this repository.")
						return nil
					}

					fmt.Println("📏 Evaluation Dimensions")
					fmt.Println("========================")
					fmt.Println()
					fmt.Println("Field Definitions:")
					fmt.Println("  • Alias: Short identifier used in measurements (e.g., 'focus')")
					fmt.Println("  • Name:  Machine-readable snake_case name used in URLs/definitions")
					fmt.Println("  • Label: Human-readable title for display purposes")
					fmt.Println()

					// Sort dimensions by alias for consistent output
					aliases := make([]string, 0, len(dimensions))
					for alias := range dimensions {
						aliases = append(aliases, alias)
					}
					sort.Strings(aliases)

					for _, alias := range aliases {
						dim := dimensions[alias]
						
						// Create human-readable label from snake_case name
						labelWords := strings.Split(dim.Name, "_")
						for i, word := range labelWords {
							labelWords[i] = strings.Title(word)
						}
						label := strings.Join(labelWords, " ")

						fmt.Printf("Alias: %s\n", dim.Alias)
						fmt.Printf("Name:  %s\n", dim.Name)
						fmt.Printf("Label: %s\n", label)
						fmt.Printf("URL:   %s\n", dim.Definition)
						fmt.Println()
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
