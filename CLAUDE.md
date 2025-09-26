# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

The Tome.gg **Librarian** is a Go-based CLI tool that implements the tome.gg protocol for validating educational content repositories. It's designed to validate directory structures containing learning materials like mental models, evaluations, training content, and DSU (Daily Stand-Up) reports.

## Build and Development Commands

### Build Commands
```bash
# Build for local development/testing
make local-build

# Build for multiple platforms (creates binaries in root)
make build
```

### Testing
```bash
# Run tests
go test ./protocol/v1/librarian/...

# Run with verbose output
go test -v ./protocol/v1/librarian/...
```

### Running the CLI
```bash
# Run locally during development
go run ./protocol/v1/librarian/cmd/main.go

# Or use the local-build target
make local-build
```

### CLI Usage Examples
```bash
# Validate current directory
go run ./protocol/v1/librarian/cmd/main.go validate

# Validate specific directory
go run ./protocol/v1/librarian/cmd/main.go validate --directory /path/to/repository

# Validate with verbose logging (shows parsing and validation details)
go run ./protocol/v1/librarian/cmd/main.go validate --directory /path/to/repository --verbose

# Initialize new repository from template
go run ./protocol/v1/librarian/cmd/main.go init --name my-repo --destination ./target-dir

# Get help
go run ./protocol/v1/librarian/cmd/main.go --help
go run ./protocol/v1/librarian/cmd/main.go validate --help

# Generate shell completions
go run ./protocol/v1/librarian/cmd/main.go completion fish > ~/.config/fish/completions/tome.fish
```

## Shell Completions

### Fish Shell
To install fish shell autocompletions for tome:

```bash
# Generate and install fish completions
tome completion fish > ~/.config/fish/completions/tome.fish

# Or during development
go run ./protocol/v1/librarian/cmd/main.go completion fish > ~/.config/fish/completions/tome.fish
```

After installation, you'll get tab completion for:
- All commands: `init`, `missing-evaluations`, `get-dsu`, `get-latest`, `validate`, `completion`
- All aliases: `missing`, `get`, `latest`
- All flags: `--directory`, `--uuid`, `--name`, `--verbose`, etc.
- Context-aware completions based on the current command

```

## Core Architecture

### CLI Commands Structure
The main CLI is in `protocol/v1/librarian/cmd/main.go` and provides two primary commands:

1. **initialize/init** - Creates new repositories from the tome.gg template using GitHub CLI
2. **validate** - Validates directory structures against the librarian protocol

### Core Components

#### Librarian Parser (`protocol/v1/librarian/librarian.go`)
- Parses directory structures recursively
- Filters files by whitelist (`.md`, `.yaml`)
- Skips blacklisted directories (`.git`)
- Builds a hierarchical `Directory` structure containing files and subdirectories

#### Domain Models (`protocol/v1/librarian/pkg/`)
- `Directory` - Represents a file system directory with validation state
- `File` - Represents individual files with validation errors
- `ValidationPlan` - Orchestrates validation across directories and files

#### Validation System (`protocol/v1/librarian/validator/`)
- **Interface-based design**: All validators implement the `Validator` interface
- **Two-phase validation**: Validates directories first, then files
- **Current validators**:
  - `DSUValidator` - Validates Daily Stand-Up report structures
  - `EvaluationValidator` - Validates evaluation/assessment content
- **Extensible**: New validators can be easily added by implementing the `Validator` interface

### Validation Flow
1. Parse directory structure into `Directory` objects
2. Create `ValidationPlan` with all directories and files
3. Register validators (DSU, Evaluation)
4. Execute two-phase validation:
   - Phase 1: Directory structure validation
   - Phase 2: File content validation
5. Collect and report errors

## File Organization

- `protocol/v1/librarian/cmd/` - CLI entry point and command definitions
- `protocol/v1/librarian/pkg/` - Core domain models and data structures
- `protocol/v1/librarian/validator/` - Validation logic and specific validators
- `docs/` - Protocol documentation and validation specifications
- `install/` - Installation scripts for different platforms

## Key Dependencies

- `github.com/urfave/cli/v2` - CLI framework
- `github.com/sirupsen/logrus` - Structured logging
- `gopkg.in/yaml.v2` - YAML parsing
- `github.com/araddon/dateparse` - Date parsing utilities

## Validation Protocol

The librarian validates tome.gg protocol compliance for educational content repositories. Key validation areas:

- **DSU Reports** (`training/*.yaml`) - Daily stand-up and progress tracking content with UUID-based training entries
- **Evaluations** (`evaluations/*.yaml`) - Assessment and examination materials
- **Training Content** - Educational materials and mental models
- **Directory Structure** - Proper organization of learning materials

### Expected Repository Structure
```
repository/
├── tome.yaml                    # Main configuration
├── training/
│   ├── dsu-reports.yaml        # Main DSU reports
│   └── dsu-reports-q*.yaml     # Quarterly reports (optional)
├── evaluations/
│   ├── eval-self.yaml          # Self evaluations
│   └── meta/
│       └── meta-self.yaml      # Meta evaluations
├── mental-models/              # Learning frameworks (optional)
├── flash-cards/               # Memory aids (optional)
└── mappings/                  # Content relationships (optional)
```

### Validation Success Criteria
- All `.yaml` files must be valid YAML format
- DSU reports must contain valid UUID training entries
- Evaluation files must follow proper structure
- Directory structure must match expected tome.gg protocol

See `docs/validations/` for detailed validation specifications.

## Logging and Debugging

The CLI uses structured logging with logrus:
- Use `--verbose` flag for debug-level logging
- Default log level is WARN
- Logs include structured fields for better debugging

## Test Repository

For testing purposes, use the existing growth journal repository at:
`/Users/darrenkarlsapalo/git/github.com/darrensapalo/founder`

## Version and Releases

Current version: 0.4.2 (see `protocol/v1/librarian/cmd/main.go`)

### Release Process

**IMPORTANT**: Always bump the version number before creating a release:

1. **Update version** in `protocol/v1/librarian/cmd/main.go` (line 23):
   ```go
   Version: "x.y.z",  // Update this line
   ```

2. **Commit the version bump**:
   ```bash
   git add protocol/v1/librarian/cmd/main.go
   git commit -m "Bump version to vx.y.z"
   git push
   ```

3. **Create GitHub release**:
   ```bash
   gh release create vx.y.z --title "vx.y.z - Release Title" --notes "Release notes..."
   ```

4. **Build and update local binary**:
   ```bash
   env GOOS=darwin GOARCH=arm64 go build -o tome-darwin-arm-osx-m1 ./protocol/v1/librarian/cmd
   cp tome-darwin-arm-osx-m1 /usr/local/bin/tome
   ```

Cross-platform binaries are automatically built via GitHub Actions for:
- macOS ARM64 (M1/M2)
- Linux AMD64
- Windows AMD64