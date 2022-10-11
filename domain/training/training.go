package training

import (
	"github.com/opensourceways/xihe-training-center/domain"
)

type JobInfo struct {
	JobId     string
	LogDir    string
	OutputDir string
}

type Training interface {
	Create(*domain.UserTraining) (JobInfo, error)
	Delete(string) error
	Get(string) (domain.JobDetail, error)
	Terminate(string) error
	GetLogDownloadURL(string) (string, error)
}
