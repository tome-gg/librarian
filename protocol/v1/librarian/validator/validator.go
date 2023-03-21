package validator

import (
	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
)

type (
	// Validator defines the necessary validation operations of a validator.
	Validator interface {
		// Directory defines the process for validating a certain directory.
		Directory(dir *pkg.Directory) error

		// File defines the process for validating a certain file.
		File(dir *pkg.File) error
	}
)


func init() {
	_ = registerValidators()
}

var validators []Validator

func registerValidators() error {
	validators = []Validator{
		NewDSUValidator(),
		NewEvaluationValidator(),
	}
	return nil
}

// Init ...
func Init(root *pkg.Directory) *pkg.ValidationPlan {

	dirs := []*pkg.Directory{}
	files := []*pkg.File{}
	
	for _, d := range root.Directories {
		dirs = append(dirs, d)
		
		for _, f := range d.Files {
			files = append(files, &f)
		}
	}

	result := pkg.NewValidationPlan(dirs, files)

	return result
}

// ValidatePlan ...
func ValidatePlan(vp *pkg.ValidationPlan) []error {
	es := []error{}
	for _, d := range vp.Directories {
		for _, validator := range validators {
			err := validator.Directory(d)
			if err != nil {
				d.Error = err
				es = append(es, err)
			}
	
			
		}
	}

	for _, f := range vp.Files {
		for _, validator := range validators {
			err := validator.File(f)
			if err != nil {
				f.Error = err
				f.Directory.ErroneousFiles = append(f.Directory.ErroneousFiles, *f)
				es = append(es, err)
			}
		}
	}
	return es
}
