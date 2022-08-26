package training

import (
	"github.com/opensourceways/xihe-training-center/domain"
)

type Training interface {
	Create(*domain.UserTraining) (string, error)
	Delete(string) error
	Get(string) (domain.TrainingDetail, error)
	Terminate(string) error
	GetLogURL(string) (string, error)
}
