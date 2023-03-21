package pkg

import (
	"fmt"
)

type (
	// Directory defines a reference to an existing folder.
	Directory struct {
		// Path defines the relative path from the root directory.
		Path string `json:"path"`
		// Directories defines the directories found in the current path.
		Directories []*Directory `json:"directories"`
		// Files defines the related files.
		Files []File `json:"files"`
		// Error defines whether an error was found during validation.
		Error error `json:"error"`
		// ErroneousFiles defines the list of files with errors.
		ErroneousFiles []File
	}

	// File defines a reference to an existing file.
	File struct {
		Directory *Directory 
		// Filepath defines where the file is found.
		Filepath string `json:"filepath"`
		// Error defines whether an error was found during validation.
		Error error `json:"error"`
	}


)


// Status returns the status of the directory, whether it was valid or not.
func (d *Directory) Status() string {

	subdirectoryStatuses := ""
	for _, dir := range d.Directories {
		subdirectoryStatuses += dir.Status()
	}
	
	if d.Error != nil {
		return fmt.Sprintf(" ❌ Path [%s] has a total of %d directories and %d files. Validation failed: %v\n%s", d.Path, len(d.Directories), len(d.Files), d.Error, subdirectoryStatuses)
		
	}

	return fmt.Sprintf(" ✅ Path [%s] has a total of %d directories and %d files.\n%s", d.Path, len(d.Directories), len(d.Files), subdirectoryStatuses)

	
	

	
}