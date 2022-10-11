package app

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/platform"
	"github.com/opensourceways/xihe-training-center/domain/synclock"
	"github.com/opensourceways/xihe-training-center/domain/syncrepo"
	"github.com/opensourceways/xihe-training-center/domain/training"
)

type TrainingCreateCmd struct {
	domain.UserTraining
}

func (cmd *TrainingCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating training")

	b := cmd.User != nil &&
		cmd.ProjectRepoId != "" &&
		cmd.ProjectName != nil &&
		cmd.Name != nil &&
		cmd.CodeDir != nil &&
		cmd.BootFile != nil

	if !b {
		return err
	}

	c := &cmd.Compute
	if c.Flavor == nil || c.Type == nil || c.Version == nil {
		return err
	}

	f := func(kv []domain.KeyValue) error {
		for i := range kv {
			if kv[i].Key == nil {
				return err
			}
		}

		return nil
	}

	if f(cmd.Hypeparameters) != nil {
		return err
	}

	if f(cmd.Env) != nil {
		return err
	}

	for i := range cmd.Inputs {
		v := &cmd.Inputs[i]
		if v.Key == nil || cmd.checkInput(&v.Value) != nil {
			return err
		}
	}

	return nil
}

func (cmd *TrainingCreateCmd) checkInput(i *domain.ResourceInput) error {
	if i.User == nil || i.Type == nil || i.RepoId == "" {
		return errors.New("invalide input")
	}

	return nil
}

type JobInfoDTO struct {
	JobId     string `json:"job_id"`
	LogDir    string `json:"log_dir"`
	OutputDir string `json:"output_dir"`
}

type JobDetailDTO struct {
	Status   string `json:"status"`
	Duration int    `json:"duration"`
}

type TrainingService interface {
	Create(cmd *TrainingCreateCmd) (JobInfoDTO, error)
	Delete(jobId string) error
	Terminate(jobId string) error
	Get(jobId string) (dto JobDetailDTO, err error)
	GetLogDownloadURL(jobId string) (string, error)
}

func NewTrainingService(
	ts training.Training,
	h syncrepo.SyncRepo,
	lock synclock.RepoSyncLock,
	p platform.Platform,
	log *logrus.Entry,
) TrainingService {
	return trainingService{
		ts: ts,
		ss: newSyncService(h, lock, p, log),
	}
}

type trainingService struct {
	ts training.Training
	ss *syncService
}

func (s trainingService) Create(cmd *TrainingCreateCmd) (dto JobInfoDTO, err error) {
	err = s.ss.syncProject(cmd.User, cmd.ProjectName, cmd.ProjectRepoId)
	if err != nil {
		return
	}

	for i := range cmd.Inputs {
		if err = s.ss.checkResourceReady(&cmd.Inputs[i].Value); err != nil {
			return
		}
	}

	v, err := s.ts.Create(&cmd.UserTraining)
	if err == nil {
		dto.JobId = v.JobId
		dto.LogDir = v.LogDir
		dto.OutputDir = v.OutputDir
	}

	return
}

func (s trainingService) Delete(jobId string) error {
	return s.ts.Delete(jobId)
}

func (s trainingService) Terminate(jobId string) error {
	return s.ts.Terminate(jobId)
}

func (s trainingService) Get(jobId string) (dto JobDetailDTO, err error) {
	v, err := s.ts.Get(jobId)
	if err != nil {
		return
	}

	dto.Status = v.Status.TrainingStatus()
	dto.Duration = v.Duration

	return
}

func (s trainingService) GetLogDownloadURL(jobId string) (string, error) {
	return s.ts.GetLogDownloadURL(jobId)
}
