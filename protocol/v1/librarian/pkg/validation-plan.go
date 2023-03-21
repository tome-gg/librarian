package pkg

import (
	"strings"
)

type (

	// ValidationPlan ...
	ValidationPlan struct {
		Directories []*Directory
		DirectoryWeights []int

		Files []*File
		FileWeights []int
	}
)

// NewValidationPlan ...
func NewValidationPlan(ds []*Directory, fs []*File) *ValidationPlan {
	return &ValidationPlan{
		Directories: ds,
		DirectoryWeights: []int{},
		Files: fs,
		FileWeights: []int{},
	}
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
