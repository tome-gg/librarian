package validator

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
	"gopkg.in/yaml.v2"
)

type evaluationValidator struct {
	log *logrus.Entry
	plan *pkg.ValidationPlan
}

var registeredDimensions []string

// File implements Validator
func (m *evaluationValidator) File(dir *pkg.File) error {
	if strings.Contains(dir.Filepath, "evaluations") == false {
		return nil
	}

	fileBytes, err := os.ReadFile(dir.Filepath)
	if err != nil {
		return err
	}

	result := pkg.EvaluationDefinition[pkg.StandardMeasurement]{}
	err = yaml.Unmarshal(fileBytes, &result)

	if err != nil {
		return err
	}

	if result.Tomegg.Type != "evaluations" {
		return nil
	}

	if result.Tomegg.Version != "0.1.0" {
		return ErrUnsupportedVersion
	}

	expectedTomeggDef := fmt.Sprintf("https://protocol.tome.gg/%s/%s", result.Tomegg.Type, result.Tomegg.Version)
	if result.Tomegg.Definition != expectedTomeggDef {
		err := ErrMismatchedTomeggDefinition(expectedTomeggDef, result.Tomegg.Definition)
		logrus.Error(err)
		return err
	}

	if len(result.Meta.Dimensions) == 0 {
		return ErrNoDimension
	}

	for _, dimension := range result.Meta.Dimensions {
		expectedDimensionDef := fmt.Sprintf("https://protocol.tome.gg/dimensions/%s/%s", dimension.Name, dimension.Version)
		if dimension.Definition != expectedDimensionDef {
			err := ErrMismatchedDimensionDefinition(dimension.Name, expectedDimensionDef, dimension.Definition)
			logrus.Error(err)
			return err
		}
		registeredDimensions = append(registeredDimensions, dimension.Name, dimension.Alias)
	}

	if len(result.Evaluations) == 0 {
		logrus.WithField("file", dir.Filepath).Warnf("empty evaluations set")
	}

	for _, records := range result.Evaluations {
		err := m.validateEvaluationRecord(records)
		if err != nil {
			return err
		}
	}

	m.log.Infof("ok")

	return nil
}

func (m *evaluationValidator) validateEvaluationRecord(records pkg.EvaluationRecord[pkg.StandardMeasurement]) error {
	if records.ID == "" {
		return ErrRequiredField(records.ID, "id")
	}

	if m.plan.IsRegistered(records.ID) == false {
		return ErrTrainingNotFound(records.ID)
	}

	if m.plan.IsValid(records.ID) == false {
		return ErrTrainingNotFound(records.ID)
	}

	if len(records.Measurements) == 0 {
		m.log.WithFields(logrus.Fields{
			"id": records.ID,
		}).Errorf(ErrNoMeasurements.Error())
		return ErrNoMeasurements
	}

	for _, measure := range records.Measurements {

		if strings.TrimSpace(measure.Dimension) == "" {
			return ErrRequiredField(records.ID, "dimension")
		}

		if measure.Score == nil {
			return ErrRequiredField(records.ID, "score")
		}

		isRegistered := false
		for _, allowed := range registeredDimensions {
			if allowed == measure.Dimension {
				isRegistered = true
				break
			}
		}

		if isRegistered == false {
			m.log.WithField("dimension", measure.Dimension).Warn(ErrUnregisteredDimension)
		}
	}

	return nil
}

// Validate defines the process for validating a certain directory.
func (m *evaluationValidator) Directory(dir *pkg.Directory) error {
	if strings.Contains(dir.Path, "evaluations") == false {
		return nil
	}

	return nil
}

// NewEvaluationValidator ...
func NewEvaluationValidator(plan *pkg.ValidationPlan) Validator {
	return &evaluationValidator{
		log: logrus.WithFields(logrus.Fields{
			"validator": "evaluation",
		}),
		plan: plan,
	}
}
