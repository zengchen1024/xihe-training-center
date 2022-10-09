package training

import (
	"github.com/opensourceways/xihe-training-center/domain"
)

type TrainingInfo struct {
	Id        string
	LogDir    string
	OutputDir string
}

type Training interface {
	Create(*domain.UserTraining) (TrainingInfo, error)
	Delete(string) error
	Get(string) (domain.TrainingDetail, error)
	Terminate(string) error
	GetLogURL(string) (string, error)
}
