# Librarian

The tome.gg **Librarian** is the system module responsible for defining the tome.gg protocol and the repository for data connectors/adaptors for loading and validating tome.gg content.

## Quick Start

### Validate a Repository
```bash
# Validate current directory
go run ./protocol/v1/librarian/cmd/main.go validate

# Validate specific directory
go run ./protocol/v1/librarian/cmd/main.go validate --directory /path/to/repository

# Validate with verbose logging
go run ./protocol/v1/librarian/cmd/main.go validate --directory /path/to/repository --verbose
```

### Initialize a New Repository
```bash
# Create a new tome.gg repository from template
go run ./protocol/v1/librarian/cmd/main.go init --name my-learning-repo --destination ./my-repo
```

### Build the CLI
```bash
# Build for local development
make local-build

# Build for all platforms
make build
```

## What Gets Validated

The librarian validates tome.gg protocol compliance for educational content repositories:

- **DSU Reports** (`training/*.yaml`) - Daily stand-up and progress tracking content
- **Evaluations** (`evaluations/*.yaml`) - Assessment and examination materials
- **Training Content** - Educational materials and mental models
- **Directory Structure** - Proper organization of learning materials

## Roadmap

1. ✅ Define system capability requirements
2. ✅ Define data definitions and validation constraints
3. ✅ Define version-controlled protocols
4. ✅ Implement validation system
5. 🔄 Implement basic Github integration

For more details, please see [docs](docs/).