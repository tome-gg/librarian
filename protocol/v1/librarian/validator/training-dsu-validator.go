package validator

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tome-gg/librarian/protocol/v1/librarian/pkg"
	"gopkg.in/yaml.v2"
)

type dailyStandUpValidator struct {
	log *logrus.Entry
	plan *pkg.ValidationPlan
}

type trainingValidator interface {
	IsRegistered(path string) bool
	IsValid(path string) bool
}

// File implements Validator
func (m *dailyStandUpValidator) File(dir *pkg.File) error {

	if strings.Contains(dir.Filepath, "dsu") == false {
		return nil
	}

	m.log.
	WithField("file", dir.Filepath).
	Debugf("DSU evaluator - processing file")

	m.log.Debugf("Files found: %+v", dir.Filepath)

	fileBytes, err := os.ReadFile(dir.Filepath)

	if err !=nil {
		return err
	}

	var result = pkg.TrainingDefinition[pkg.DSUReport]{}
	err = yaml.Unmarshal(fileBytes, &result)
	
	if err != nil {
		return err
	}

	if result.Tomegg.Type != "training" {
		return nil
	}

	if result.Tomegg.Version != "0.1.0" {
		return ErrUnsupportedVersion
	}

	if result.Meta.Format.Type != "dsu" {
		return ErrUnsupportedFormat
	}

	if result.Meta.Format.Version != "0.1.0" {
		return ErrUnsupportedVersion
	}

	expectedTomeggDef := fmt.Sprintf("https://protocol.tome.gg/%s/%s", result.Tomegg.Type, result.Tomegg.Version)
	if result.Tomegg.Definition != expectedTomeggDef {
		return ErrMismatchedTomeggDefinition(expectedTomeggDef, result.Tomegg.Definition)
	}

	expectedFormatDef := fmt.Sprintf("https://protocol.tome.gg/formats/%s/%s", result.Meta.Format.Type, result.Meta.Format.Version)
	if result.Meta.Format.Definition != expectedFormatDef {
		return ErrMismatchedFormatDefinition(result.Meta.Format.Type, expectedFormatDef, result.Meta.Format.Definition)
	}

	if len(result.Content) == 0 {
		m.log.Warnf("empty training set")
	}

	for _, e := range result.Content {
		m.plan.Metadata["registeredTraining"] = append(m.plan.Metadata["registeredTraining"].([]string), e.ID)
		m.log.WithField("training", e.ID).Debugf("registered training")

		err := m.validateDSUEntry(e)
		if err != nil {
			return err
		}
		m.plan.Metadata["validTraining"] = append(m.plan.Metadata["validTraining"].([]string), e.ID)
	}

	m.log.
		WithField("validator", "training").
		WithField("type", "dsu").
		Infof("ok")

	return nil
}


func (m *dailyStandUpValidator) validateDSUEntry(e pkg.DSUReport) error {
	if strings.TrimSpace(e.DoingToday) == "" {
		return ErrRequiredField(e.ID, "doing_today")
	}
	if strings.TrimSpace(e.DoneYesterday) == "" {
		return ErrRequiredField(e.ID, "done_yesterday")
	}
	if strings.TrimSpace(e.DatetimeRaw) == "" {
		return ErrRequiredField(e.ID, "datetime")
	}
	if strings.TrimSpace(e.ID) == "" {
		return ErrRequiredField(e.ID, "id")
	}
	return nil
}

// Validate defines the process for validating a certain directory.
func (m *dailyStandUpValidator) Directory(dir *pkg.Directory) error {
	if strings.Contains(dir.Path, "training") == false {
		return nil
	}

	return nil
}

// NewDSUValidator ...
func NewDSUValidator(plan *pkg.ValidationPlan) Validator {
	return &dailyStandUpValidator{
		log: logrus.WithFields(logrus.Fields{
			"validator": "training",
			"type": "dsu",
		}),
		plan: plan,
	}
}
