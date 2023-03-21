package librarian

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
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

func getDirectory(memo map[string]*pkg.Directory, path string) *pkg.Directory {
	directory, ok := memo[path]

	if !ok {
		directory = &pkg.Directory{
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
func Parse(rootDirectory string) (*pkg.Directory, error) {
	memo := map[string]*pkg.Directory{}
	memo[rootDirectory] = &pkg.Directory{
		Path: rootDirectory,
	}

	err := filepath.WalkDir(rootDirectory, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		isDirectory := d.IsDir()
		isSubDirectory := isDirectory && path != rootDirectory

		parentDirPath := getParentDirectory(rootDirectory, path)
		parentDirectory := memo[parentDirPath]

		if parentDirectory == nil {
			memo[parentDirPath] = &pkg.Directory{}
			parentDirectory = memo[parentDirPath]
		}
		

		switch isDirectory {
		case true: // It is a directory
			if shouldSkipDirectory(path) {
				return filepath.SkipDir
			}

			directory := getDirectory(memo, path)

			if isSubDirectory {
				logrus.WithFields(logrus.Fields{
					"parent": parentDirectory.Path,
					"subdir": path,
				}).Debugf("adding dir")
				parentDirectory.Directories = append(parentDirectory.Directories, directory)
			}

		case false: // It is a file

		f := pkg.File{
			Directory: parentDirectory,
			Filepath: path,
		}

		logrus.WithFields(logrus.Fields{
			"parent": parentDirectory.Path,
			"file": path,
		}).Debugf("adding file")
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
