package pkg

// EvaluationDefinition defines the training definition for a daily stand up
type EvaluationDefinition[E any] struct {
	Tomegg struct {
		Type       string `yaml:"type"`
		Version    string `yaml:"version"`
		Definition string `yaml:"definition"`
	} `yaml:"tomegg"`

	Meta struct {
		Dimensions []struct {
			Alias      string `yaml:"alias"`
			Name       string `yaml:"name"`
			Version    string `yaml:"version"`
			Definition string `yaml:"definition"`
		} `yaml:"dimensions"`
	} `yaml:"meta"`

	Evaluations []EvaluationRecord[E] `yaml:"evaluations"`
}

// EvaluationRecord ...
type EvaluationRecord[E any] struct {
	ID           string                `yaml:"id"`
	Measurements []StandardMeasurement `yaml:"measurements"`
}

// StandardMeasurement ...
type StandardMeasurement struct {
	Dimension string `yaml:"dimension"`
	Score     *int   `yaml:"score"`
	Remarks   string `yaml:"remarks"`
	Wins      string `yaml:"wins"`
	Mistakes  string `yaml:"mistakes"`
	Meta 			string `yaml:"meta"`
}
