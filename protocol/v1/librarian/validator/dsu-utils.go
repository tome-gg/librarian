package validator

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
	"gopkg.in/yaml.v2"
)

// FindMissingEvaluations returns DSU entries that don't have corresponding self evaluations
// If limitToLast3 is true, returns only the 3 most recent entries in ascending order
func FindMissingEvaluations(plan *pkg.ValidationPlan, limitToLast3 bool) ([]pkg.DSUReport, error) {
	// First, collect all DSU entries
	dsuEntries, err := getAllDSUEntries(plan)
	if err != nil {
		return nil, err
	}

	// Then, collect all evaluation IDs
	evaluationIDs, err := getAllEvaluationIDs(plan)
	if err != nil {
		return nil, err
	}

	// Find DSU entries without evaluations
	var missingEvaluations []pkg.DSUReport
	for _, dsu := range dsuEntries {
		hasEvaluation := false
		for _, evalID := range evaluationIDs {
			if evalID == dsu.ID {
				hasEvaluation = true
				break
			}
		}
		if !hasEvaluation {
			missingEvaluations = append(missingEvaluations, dsu)
		}
	}

	// Sort in ascending order by date
	sort.Slice(missingEvaluations, func(i, j int) bool {
		return missingEvaluations[i].Datetime.Before(missingEvaluations[j].Datetime)
	})

	// Limit to last 3 if requested
	if limitToLast3 && len(missingEvaluations) > 3 {
		missingEvaluations = missingEvaluations[len(missingEvaluations)-3:]
	}

	return missingEvaluations, nil
}

// GetDSUByUUID retrieves a DSU entry by its UUID
func GetDSUByUUID(plan *pkg.ValidationPlan, uuid string) (*pkg.DSUReport, error) {
	dsuEntries, err := getAllDSUEntries(plan)
	if err != nil {
		return nil, err
	}

	for _, dsu := range dsuEntries {
		if dsu.ID == uuid {
			return &dsu, nil
		}
	}

	return nil, fmt.Errorf("DSU entry with UUID %s not found", uuid)
}

// getAllDSUEntries collects all DSU entries from training files
func getAllDSUEntries(plan *pkg.ValidationPlan) ([]pkg.DSUReport, error) {
	var allEntries []pkg.DSUReport

	for _, file := range plan.Files {
		if !strings.Contains(file.Filepath, "dsu") || !strings.Contains(file.Filepath, "training") {
			continue
		}

		fileBytes, err := os.ReadFile(file.Filepath)
		if err != nil {
			continue // Skip files we can't read
		}

		var result = pkg.TrainingDefinition[pkg.DSUReport]{}
		err = yaml.Unmarshal(fileBytes, &result)
		if err != nil {
			continue // Skip invalid YAML files
		}

		if result.Tomegg.Type != "training" || result.Meta.Format.Type != "dsu" {
			continue // Skip non-DSU training files
		}

		allEntries = append(allEntries, result.Content...)
	}

	return allEntries, nil
}

// GetLatestDSU retrieves the most recent DSU entry by date
func GetLatestDSU(plan *pkg.ValidationPlan) (*pkg.DSUReport, error) {
	dsuEntries, err := getAllDSUEntries(plan)
	if err != nil {
		return nil, err
	}

	if len(dsuEntries) == 0 {
		return nil, fmt.Errorf("no DSU entries found")
	}

	// Find the entry with the latest date
	var latestEntry *pkg.DSUReport = &dsuEntries[0]
	for i := 1; i < len(dsuEntries); i++ {
		if dsuEntries[i].Datetime.After(latestEntry.Datetime) {
			latestEntry = &dsuEntries[i]
		}
	}

	return latestEntry, nil
}

// getAllEvaluationIDs collects all evaluation IDs from evaluation files
func getAllEvaluationIDs(plan *pkg.ValidationPlan) ([]string, error) {
	var allIDs []string

	for _, file := range plan.Files {
		if !strings.Contains(file.Filepath, "evaluations") {
			continue
		}

		fileBytes, err := os.ReadFile(file.Filepath)
		if err != nil {
			continue // Skip files we can't read
		}

		result := pkg.EvaluationDefinition[pkg.StandardMeasurement]{}
		err = yaml.Unmarshal(fileBytes, &result)
		if err != nil {
			continue // Skip invalid YAML files
		}

		if result.Tomegg.Type != "evaluations" {
			continue // Skip non-evaluation files
		}

		for _, evaluation := range result.Evaluations {
			allIDs = append(allIDs, evaluation.ID)
		}
	}

	return allIDs, nil
}