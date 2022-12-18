package librarian

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// directoryBlacklist which directories are ignored?
var directoryBlacklist = []string{
	".git",
}

// fileExtensionWhitelist which file extension are to be included for processing?
var fileExtensionWhitelist = []string{
	".md",
	".yaml",
}

func shouldSkipDirectory(path string) bool {
	for _, dir := range directoryBlacklist {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

func getDirectory(memo map[string]*Directory, path string) *Directory {
	directory, ok := memo[path]

	if !ok {
		directory = &Directory{
			Path: path,
		}
		memo[path] = directory
	}

	return directory
}

func getParentDirectory(root string, path string) string {
	lastIndex := strings.LastIndex(path, "/")

	if lastIndex == -1 {
		return root
	}

	return path[0:lastIndex]
}

// Parse parses a specified file path and returns a librarian.Directory.
func Parse(rootDirectory string) (*Directory, error) {
	memo := map[string]*Directory{}

	err := filepath.WalkDir(rootDirectory, func(path string, d fs.DirEntry, err error) error {
		isDirectory := d.IsDir()
		isSubDirectory := isDirectory && path != rootDirectory

		parentDirPath := getParentDirectory(rootDirectory, path)
		parentDirectory := memo[parentDirPath]
		

		switch isDirectory {
		case true: // It is a directory
			if shouldSkipDirectory(path) {
				return filepath.SkipDir
			}

			directory := getDirectory(memo, path)

			if isSubDirectory {
				logrus.Debugf("Adding subdirectory %s to parent directory %s", path, parentDirectory.Path)
				parentDirectory.Directories = append(parentDirectory.Directories, directory)
			}

		case false: // It is a file

		f := File{
			Filepath: path,
		}

		logrus.Debugf("Adding file %s to directory %s", path, parentDirectory.Path)
		parentDirectory.Files = append(parentDirectory.Files, f)

			// for _, ext := range fileExtensionWhitelist {
			// 	if !isDir && !strings.HasSuffix(path, ext) {
			// 		logrus.Warnf("Suffix %s was not found on path %s", ext, path)
			// 		return nil
			// 	}
			// }
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return err
	})

	return memo[rootDirectory], err
}
