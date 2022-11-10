package app

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/platform"
	"github.com/opensourceways/xihe-training-center/domain/synclock"
	"github.com/opensourceways/xihe-training-center/domain/training"
	"github.com/opensourceways/xihe-training-center/domain/watch"
)

type TrainingCreateCmd struct {
	ProjectId  string
	TrainingId string

	domain.UserTraining
}

func (cmd *TrainingCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating training")

	b := cmd.User != nil &&
		cmd.ProjectRepoId != "" &&
		cmd.ProjectName != nil &&
		cmd.Name != nil &&
		cmd.CodeDir != nil &&
		cmd.BootFile != nil &&
		cmd.ProjectId != "" &&
		cmd.TrainingId != ""

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

		if v.Key == nil || v.User == nil || v.RepoId == "" {
			return errors.New("invalide input")
		}
	}

	return nil
}

type JobInfoDTO struct {
	JobId     string `json:"job_id"`
	LogDir    string `json:"log_dir"`
	AimDir    string `json:"aim_dir"`
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
	GetLogDownloadURL(jobId string) (string, error)
	GenFileDownloadURL(obsfile string) (string, error)
}

func NewTrainingService(
	ts training.Training,
	pf platform.Platform,
	ws watch.WatchService,
	log *logrus.Entry,
	lock synclock.RepoSyncLock,
	maxTrainingNum int,
) TrainingService {
	t := &trainingService{
		ts:  ts,
		ws:  ws,
		log: log,
		ss:  newSyncService(ts, pf, log, lock),

		maxTrainingNum: maxTrainingNum,
	}

	ws.RegisterTrainingDone(t.callback)

	return t
}

type trainingService struct {
	ss  *syncService
	log *logrus.Entry
	ts  training.Training
	ws  watch.WatchService

	lock           sync.RWMutex
	currentNum     int
	maxTrainingNum int
}

func (s *trainingService) callback(*watch.TrainingInfo) {
	s.lock.Lock()
	s.currentNum--
	s.lock.Unlock()
}

func (s *trainingService) Create(cmd *TrainingCreateCmd) (dto JobInfoDTO, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.currentNum >= s.maxTrainingNum {
		err = errors.New("too many trainings")

		return
	}

	err = s.ss.syncProject(cmd.User, cmd.ProjectName, cmd.ProjectRepoId)
	if err != nil {
		s.log.Debugf(
			"sync project(%s) failed",
			cmd.User.Account(), cmd.ProjectId,
		)

		return
	}

	for i := range cmd.Inputs {
		dep := &cmd.Inputs[i].ResourceRef

		if err = s.ss.checkResourceReady(dep); err != nil {
			s.log.Debugf(
				"check dependent resource:%s failed, err:%s",
				dep.ToRepoPath(), err.Error(),
			)

			return
		}
	}

	v, err := s.ts.Create(&cmd.UserTraining)
	if err != nil {
		return
	}

	dto.JobId = v.JobId
	dto.LogDir = v.LogDir
	dto.AimDir = v.AimDir
	dto.OutputDir = v.OutputDir

	s.ws.WatchTraining(&watch.TrainingInfo{
		User:       cmd.User,
		ProjectId:  cmd.ProjectId,
		TrainingId: cmd.TrainingId,
		JobInfo:    v,
	})

	s.currentNum++

	return
}

func (s *trainingService) Delete(jobId string) error {
	return s.ts.Delete(jobId)
}

func (s *trainingService) Terminate(jobId string) error {
	return s.ts.Terminate(jobId)
}

func (s *trainingService) GetLogDownloadURL(jobId string) (string, error) {
	return s.ts.GetLogDownloadURL(jobId)
}

func (s *trainingService) GenFileDownloadURL(obsfile string) (string, error) {
	return s.ts.GenFileDownloadURL(obsfile)
}
