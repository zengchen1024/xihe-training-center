package watch

type Config struct {
	// MaxTrainingNum specifies the max num of training
	// which the training center can support
	MaxTrainingNum int `json:"max_training_num"`

	// Interval specifies the interval of second between two loops
	// that check all trainings in a loop.
	Interval int `json:"interval"`
}
