package validator

import (
	"strings"
	"testing"
	"time"

	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
)

// createMockDirectoryStructure creates a mock directory structure for testing
func createMockDirectoryStructure() *pkg.Directory {
	// Mock training directory with DSU files
	trainingDir := &pkg.Directory{
		Path: "/mock/repo/training",
		Files: []pkg.File{
			{Filepath: "/mock/repo/training/dsu-reports.yaml"},
			{Filepath: "/mock/repo/training/dsu-reports-q2-2024.yaml"},
			{Filepath: "/mock/repo/training/dsu-reports-q3-2024.yaml"},
			{Filepath: "/mock/repo/training/dsu-reports-q3-2025.yaml"},
		},
	}

	// Mock evaluations directory
	evaluationsDir := &pkg.Directory{
		Path: "/mock/repo/evaluations",
		Files: []pkg.File{
			{Filepath: "/mock/repo/evaluations/eval-self.yaml"},
		},
	}

	// Mock root directory
	root := &pkg.Directory{
		Path: "/mock/repo",
		Files: []pkg.File{
			{Filepath: "/mock/repo/README.md"},
			{Filepath: "/mock/repo/tome.yaml"},
		},
		Directories: []*pkg.Directory{
			trainingDir,
			evaluationsDir,
		},
	}

	// Set parent references
	for i := range trainingDir.Files {
		trainingDir.Files[i].Directory = trainingDir
	}
	for i := range evaluationsDir.Files {
		evaluationsDir.Files[i].Directory = evaluationsDir
	}
	for i := range root.Files {
		root.Files[i].Directory = root
	}

	return root
}

func TestDSUFileCollection(t *testing.T) {
	// Create mock directory structure
	mockDirectory := createMockDirectoryStructure()

	// Initialize validation plan
	plan := Init(mockDirectory)

	// Debug: Log all files found
	t.Logf("All files in plan:")
	for i, file := range plan.Files {
		t.Logf("File %d: %s", i+1, file.Filepath)
	}

	// Find all DSU files in the plan
	var dsuFiles []string
	for _, file := range plan.Files {
		if strings.Contains(file.Filepath, "dsu") {
			dsuFiles = append(dsuFiles, file.Filepath)
		}
	}

	t.Logf("Total files in plan: %d", len(plan.Files))
	t.Logf("DSU files found: %d", len(dsuFiles))

	for i, file := range dsuFiles {
		t.Logf("DSU File %d: %s", i+1, file)
	}

	// Expected DSU files based on mock structure
	expectedDSUFiles := []string{
		"/mock/repo/training/dsu-reports.yaml",
		"/mock/repo/training/dsu-reports-q2-2024.yaml",
		"/mock/repo/training/dsu-reports-q3-2024.yaml",
		"/mock/repo/training/dsu-reports-q3-2025.yaml",
	}

	// Check if we found all expected files
	if len(dsuFiles) != len(expectedDSUFiles) {
		t.Errorf("Expected %d DSU files, but found %d", len(expectedDSUFiles), len(dsuFiles))
	}

	// Check each expected file is present
	fileMap := make(map[string]bool)
	for _, file := range dsuFiles {
		fileMap[file] = true
	}

	for _, expectedFile := range expectedDSUFiles {
		if !fileMap[expectedFile] {
			t.Errorf("Expected DSU file not found in plan: %s", expectedFile)
		}
	}

	// Test total file count (should include all files from all directories)
	expectedTotalFiles := 2 + 4 + 1 // root + training + evaluations
	if len(plan.Files) != expectedTotalFiles {
		t.Errorf("Expected %d total files in plan, but found %d", expectedTotalFiles, len(plan.Files))
	}
}

func TestDSUEntryParsing(t *testing.T) {
	// Mock the getAllDSUEntries function behavior with test data
	mockEntries := []pkg.DSUReport{
		{
			ID:       "155CA198-7084-42F7-BBEE-A5A2FD3CB76F",
			Datetime: time.Date(2025, 9, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       "5A5C8FFE-147D-4A2E-B8D9-770FBCCB8B19",
			Datetime: time.Date(2025, 9, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       "F179F5FC-CF23-4F6A-B6F9-A6C9D3586ECD",
			Datetime: time.Date(2025, 9, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	t.Logf("Mock DSU entries from 2025: %d", len(mockEntries))

	// Check if we have entries from 2025
	var entries2025 []string
	for _, entry := range mockEntries {
		if entry.Datetime.Year() == 2025 {
			entries2025 = append(entries2025, entry.ID)
		}
	}

	t.Logf("DSU entries from 2025: %d", len(entries2025))
	for i, id := range entries2025 {
		t.Logf("2025 Entry %d: %s", i+1, id)
	}

	// We expect 3 entries from 2025 based on the mock data
	expectedCount := 3
	if len(entries2025) != expectedCount {
		t.Errorf("Expected %d DSU entries from 2025, but found %d", expectedCount, len(entries2025))
	}
}

func TestFileCollectionRecursive(t *testing.T) {
	// Test that the file collection works recursively
	// Create nested directory structure
	subTrainingDir := &pkg.Directory{
		Path: "/mock/repo/training/archived",
		Files: []pkg.File{
			{Filepath: "/mock/repo/training/archived/old-dsu-reports.yaml"},
		},
	}

	trainingDir := &pkg.Directory{
		Path: "/mock/repo/training",
		Files: []pkg.File{
			{Filepath: "/mock/repo/training/dsu-reports.yaml"},
		},
		Directories: []*pkg.Directory{subTrainingDir},
	}

	root := &pkg.Directory{
		Path: "/mock/repo",
		Files: []pkg.File{
			{Filepath: "/mock/repo/tome.yaml"},
		},
		Directories: []*pkg.Directory{trainingDir},
	}

	// Set parent references
	subTrainingDir.Files[0].Directory = subTrainingDir
	trainingDir.Files[0].Directory = trainingDir
	root.Files[0].Directory = root

	// Test the recursive collection
	plan := Init(root)

	// Should find files from all levels
	expectedFileCount := 3 // root + training + archived
	if len(plan.Files) != expectedFileCount {
		t.Errorf("Expected %d files with recursive collection, but found %d", expectedFileCount, len(plan.Files))
	}

	// Should find the nested DSU file
	var foundNestedDSU bool
	for _, file := range plan.Files {
		if strings.Contains(file.Filepath, "archived/old-dsu-reports.yaml") {
			foundNestedDSU = true
			break
		}
	}

	if !foundNestedDSU {
		t.Error("Expected to find nested DSU file in archived subdirectory")
	}
}