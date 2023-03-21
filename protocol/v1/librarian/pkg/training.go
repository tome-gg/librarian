package pkg

import "time"

// TrainingDefinition defines the training definition for a daily stand up
type TrainingDefinition[E any] struct {
	Tomegg struct {
			Type        string `yaml:"type"`
			Version     string `yaml:"version"`
			Definition  string `yaml:"definition"`
	} `yaml:"tomegg"`

	Meta struct {
			Format struct {
					Type        string `yaml:"type"`
					Version     string `yaml:"version"`
					Definition  string `yaml:"definition"`
			} `yaml:"format"`

			Tags []string `yaml:"tags"`
	} `yaml:"meta"`

	Content []E `yaml:"content"`
}

// TrainingBase ...
type TrainingBase struct {
	ID string
	Date time.Time
	Remarks string
}