package pkg

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type (

	// ValidationPlan ...
	ValidationPlan struct {
		Directories []*Directory
		DirectoryWeights []int

		Files []*File
		FileWeights []int
		Metadata map[string]interface{}
	}
)

// NewValidationPlan ...
func NewValidationPlan(ds []*Directory, fs []*File) *ValidationPlan {
	// Remove redundancies from fs and ds
	uniqueDirs := make(map[string]*Directory)
	var dedupedDirs []*Directory
	for _, dir := range ds {
		if _, exists := uniqueDirs[dir.Path]; !exists {
			uniqueDirs[dir.Path] = dir
			dedupedDirs = append(dedupedDirs, dir)
		}
	}
	
	// Remove redundancies from files
	uniqueFiles := make(map[string]*File)
	var dedupedFiles []*File
	for _, file := range fs {
		if _, exists := uniqueFiles[file.Filepath]; !exists {
			uniqueFiles[file.Filepath] = file
			dedupedFiles = append(dedupedFiles, file)
		}
	}

	for _, f := range dedupedFiles {
		logrus.Debugf("Files added: %s", f.Filepath)
	}

	return &ValidationPlan{
		Directories: dedupedDirs,
		DirectoryWeights: []int{},
		Files: dedupedFiles,
		FileWeights: []int{},
		Metadata: map[string]interface{}{
			"registeredTraining": []string{},
			"validTraining": []string{},
		},
	}
}


// IsRegistered returns true if the file has been registered.
func (vp *ValidationPlan) IsRegistered(path string) bool {
	if rt, ok := vp.Metadata["registeredTraining"].([]string); ok {
		for _, t := range rt {
			if t == path {
				return true
			}
		}
	}
	return false
}

// IsValid returns true if the file is valid.
func (vp *ValidationPlan) IsValid(path string) bool {
	if vt, ok := vp.Metadata["validTraining"].([]string); ok {
		for _, t := range vt {
			if t == path {
				return true
			}
		}
	}
	return false
}


type swappable interface {
	File | Directory
}
func swap[E swappable](a *E, b *E) {
	temp := *a
	*a = *b
	*b = temp
}

// Init initializes the weights of the directories and files.
func (vp *ValidationPlan) Init() {
	for i := range vp.Directories {
		vp.DirectoryWeights = append(vp.DirectoryWeights, vp.getDirectoryWeight(vp.Directories[i]))
	}

	for i := range vp.Files {
		vp.FileWeights = append(vp.FileWeights, vp.getFileWeight(vp.Files[i]) * vp.getDirectoryWeight(vp.Files[i].Directory))
	}

	// bubble sort; lol
	for i := range vp.Directories {
		for j := range vp.Directories {
			if vp.DirectoryWeights[i] > vp.DirectoryWeights[j] {
				swap(vp.Directories[i], vp.Directories[j])
			}
		}
	}

	// bubble sort; lol
	for i := range vp.Files {
		for j := range vp.Files {
			if vp.FileWeights[i] > vp.FileWeights[j] {
				swap(vp.Files[i], vp.Files[j])
			}
		}
	}
}

func (vp *ValidationPlan) getDirectoryWeight(d *Directory) int {
	if strings.Contains(d.Path, "evaluations") {
		return 100
	}
	if strings.Contains(d.Path, "meta") {
		return 110
	}
	if strings.Contains(d.Path, "flash-cards") {
		return 90
	}
	if strings.Contains(d.Path, "mappings") {
		return 90
	}
	if strings.Contains(d.Path, "mental-models") {
		return 90
	}
	if strings.Contains(d.Path, "training") {
		return 10
	}
	return 100
}

func  (vp *ValidationPlan) getFileWeight(f *File) int {
	return 100
}
