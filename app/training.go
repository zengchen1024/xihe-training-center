package app

import (
	"errors"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/training"
)

type TrainingCreateCmd struct {
	domain.UserTraining
}

func (cmd *TrainingCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating training")

	b := cmd.User != nil &&
		cmd.ProjectId != "" &&
		cmd.Name != nil &&
		cmd.CodeDir != nil &&
		cmd.BootFile != nil &&
		cmd.LogDir != nil

	if !b {
		return err
	}

	c := &cmd.Compute
	b = c.Flavor != nil && c.Type != nil && c.Version != nil
	if !b {
		return err
	}

	f := func(kv []domain.KeyValue) error {
		for i := range kv {
			if kv[i].Key == nil || kv[i].Value == nil {
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

	if f(cmd.Inputs) != nil {
		return err
	}

	if f(cmd.Outputs) != nil {
		return err
	}

	return nil
}

type TrainingDTO struct {
	Status   string `json:"status"`
	Duration int    `json:"duration"`
}

type TrainingService interface {
	Create(cmd *TrainingCreateCmd) (string, error)
	Delete(jobId string) error
	Terminate(jobId string) error
	Get(jobId string) (dto TrainingDTO, err error)
	GetLogURL(jobId string) (string, error)
}

func NewTrainingService(ts training.Training) TrainingService {
	return trainingService{ts}
}

type trainingService struct {
	ts training.Training
}

func (s trainingService) Create(cmd *TrainingCreateCmd) (string, error) {
	return s.ts.Create(&cmd.UserTraining)
}

func (s trainingService) Delete(jobId string) error {
	return s.ts.Delete(jobId)
}

func (s trainingService) Terminate(jobId string) error {
	return s.ts.Terminate(jobId)
}

func (s trainingService) Get(jobId string) (dto TrainingDTO, err error) {
	v, err := s.ts.Get(jobId)
	if err != nil {
		return
	}

	dto.Status = v.Status.TrainingStatus()
	dto.Duration = v.Duration.TrainingDuration()

	return
}

func (s trainingService) GetLogURL(jobId string) (string, error) {
	return s.ts.GetLogURL(jobId)
}
