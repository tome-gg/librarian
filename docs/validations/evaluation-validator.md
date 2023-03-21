# DSU validator

Source: `protocol/v1/librarian/validator/training-dsu-validator.go`

## Supported validations

1. Tome.gg evaluation YAML format
2. Tome.gg evaluation definition and meta format matching (evaluation YAML format definition matching, and meta format [i.e. DSU] definition matching)
3. Warning for empty evaluation set
4. Required fields for evaluation (`id`, `dimension`, `score`)
5. Evaluation must match an existing training reference
6. Checks for dimension registry

## Roadmap

None yet.