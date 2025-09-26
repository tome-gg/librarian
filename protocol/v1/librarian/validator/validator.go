package validator

import (
	"github.com/sirupsen/logrus"
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

var validators []Validator

func registerValidators(_ *pkg.Directory, plan *pkg.ValidationPlan) error {

	validators = []Validator{
		NewDSUValidator(plan),
		NewEvaluationValidator(plan),
	}
	return nil
}

// Init ...
func Init(root *pkg.Directory) *pkg.ValidationPlan {

	dirs := []*pkg.Directory{}
	files := []*pkg.File{}

	// Add root directory files
	for i := range root.Files {
		files = append(files, &root.Files[i])
	}

	// Recursively collect all directories and files
	collectDirectoriesAndFiles(root, &dirs, &files)

	for _, f := range files {
		logrus.Debugf("Files added: %s", f.Filepath)
	}

	plan := pkg.NewValidationPlan(dirs, files)

	registerValidators(root, plan)

	return plan
}

// collectDirectoriesAndFiles recursively collects all directories and files
func collectDirectoriesAndFiles(dir *pkg.Directory, dirs *[]*pkg.Directory, files *[]*pkg.File) {
	for _, d := range dir.Directories {
		*dirs = append(*dirs, d)

		// Add ALL files from this directory (not just the last one)
		for i := range d.Files {
			*files = append(*files, &d.Files[i])
		}

		// Recursively process subdirectories
		collectDirectoriesAndFiles(d, dirs, files)
	}
}

// ValidatePlan ...
func ValidatePlan(vp *pkg.ValidationPlan) []error {

	logrus.Debugf("Validating the plan: Step 1 - validate directories")
	
	es := []error{}
	for _, d := range vp.Directories {
		logrus.Debugf("Validating dir: %s", d.Path)
		for _, validator := range validators {
			err := validator.Directory(d)
			if err != nil {
				d.Error = err
				es = append(es, err)
			}
	
			
		}
	}

	logrus.Debugf("Validating the plan: Step 2 - validate files")

	for _, f := range vp.Files {
		logrus.Debugf("File to be validated: %s", f.Filepath)
	}

	for _, f := range vp.Files {
		logrus.Debugf("Validating file: %s", f.Filepath)
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
