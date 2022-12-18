package validator

import (
	validate "github.com/go-playground/validator/v10"
	librarian "github.com/tome-gg/librarian/protocol/v1/librarian"
)

type (
	// Validator defines the necessary validation operations of a validator.
	Validator interface {
		// Validate defines the process for validating a certain directory.
		Validate(dir *librarian.Directory) error
	}
)

// Validate performs validation on the given root directory.
func Validate(root *librarian.Directory) {
	v := validate.New()
	err := v.Struct(root)
	if err != nil {
		root.Error = err
	}
}
