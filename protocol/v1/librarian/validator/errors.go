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

// ErrMismatchedTomeggDefinition creates a specific error for mismatched tomegg definitions
func ErrMismatchedTomeggDefinition(expected, actual string) error {
	return fmt.Errorf("mismatched tomegg definition: expected '%s' but got '%s'", expected, actual)
}

// ErrMismatchedDimensionDefinition creates a specific error for mismatched dimension definitions
func ErrMismatchedDimensionDefinition(dimensionName, expected, actual string) error {
	return fmt.Errorf("mismatched definition for dimension '%s': expected '%s' but got '%s'", dimensionName, expected, actual)
}

// ErrMismatchedFormatDefinition creates a specific error for mismatched format definitions
func ErrMismatchedFormatDefinition(formatType, expected, actual string) error {
	return fmt.Errorf("mismatched definition for format '%s': expected '%s' but got '%s'", formatType, expected, actual)
}

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