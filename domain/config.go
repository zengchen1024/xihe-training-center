package domain

var config Config

func Init(cfg *Config) {
	config = *cfg
}

type Config struct {
	MaxTrainingNameLength int `json:"max_training_name_length"`
	MinTrainingNameLength int `json:"min_training_name_length"`
	MaxTrainingDescLength int `json:"max_training_desc_length"`
}

func (r *Config) Setdefault() {
	if r.MaxTrainingNameLength == 0 {
		r.MaxTrainingNameLength = 50
	}

	if r.MinTrainingNameLength == 0 {
		r.MinTrainingNameLength = 5
	}

	if r.MaxTrainingDescLength == 0 {
		r.MaxTrainingDescLength = 100
	}
}
