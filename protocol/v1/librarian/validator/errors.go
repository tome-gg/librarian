package validator

import "fmt"

// ErrInvalidTrainingType ...
var ErrInvalidTrainingType = fmt.Errorf("invalid training file type")

// ErrUnsupportedVersion ...
var ErrUnsupportedVersion = fmt.Errorf("unsupported training version")

// ErrUnsupportedFormat ...
var ErrUnsupportedFormat = fmt.Errorf("unsupported training format")

// ErrMismatchedDefinition ...
var ErrMismatchedDefinition = fmt.Errorf("mismatched definition")

// ErrNoDimension ...
var ErrNoDimension = fmt.Errorf("no dimension specified for evaluation")

// ErrUnregisteredDimension ...
var ErrUnregisteredDimension = fmt.Errorf("dimension not registered for evaluation")

// ErrNoMeasurements ...
var ErrNoMeasurements = fmt.Errorf("no measurements found")

// ErrRequiredField ...
func ErrRequiredField(id string, field string) error {
	return fmt.Errorf("required field %s for content entry %s", field, id)
}

// ErrTrainingNotFound ...
func ErrTrainingNotFound(id string) error {
	return fmt.Errorf("specified training %s was not found", id)
}