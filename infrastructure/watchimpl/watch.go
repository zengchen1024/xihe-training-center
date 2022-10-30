package watchimpl

import (
	"time"

	pt "github.com/opensourceways/xihe-grpc-protocol/training"
	"github.com/opensourceways/xihe-grpc-protocol/training/client"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-training-center/domain/training"
	"github.com/opensourceways/xihe-training-center/domain/watch"
)

type trainingData = pt.TrainingInfo

func NewWatcher(
	cfg *Config, ts training.Training,
	maxTrainingNum int, log *logrus.Entry,
) (*Watcher, error) {
	cli, err := client.NewClient(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		log:       log,
		cli:       cli,
		ts:        ts,
		interval:  time.Duration(cfg.Interval) * time.Second,
		stop:      make(chan struct{}),
		stopped:   make(chan struct{}),
		trainings: make(chan trainingInfo, maxTrainingNum+1),
	}, nil
}

type trainingInfo struct {
	watch.TrainingInfo

	// TODO if timeout, ignore this work and set status to timeout
	//timeout      int

	result trainingData

	done         bool
	success      bool
	logDone      bool
	aimDone      bool
	outputDone   bool
	notifyFailed bool
}

func (t *trainingInfo) toIndex() pt.TrainingIndex {
	return pt.TrainingIndex{
		Id:        t.TrainingId,
		User:      t.User.Account(),
		ProjectId: t.ProjectId,
	}
}

func (t *trainingInfo) isDone() bool {
	// TODO: check if timeout for this training. return true if timeout

	done := t.done && t.logDone

	if done && t.success {
		done = t.outputDone && t.aimDone
	}

	return done
}

type Watcher struct {
	log *logrus.Entry
	cli *client.Client
	ts  training.Training

	interval time.Duration

	stop      chan struct{}
	stopped   chan struct{}
	trainings chan trainingInfo

	callback func(*watch.TrainingInfo)
}

func (w *Watcher) WatchTraining(t *watch.TrainingInfo) {
	w.trainings <- trainingInfo{TrainingInfo: *t}
}

func (w *Watcher) RegisterTrainingDone(f func(*watch.TrainingInfo)) {
	w.callback = f
}

func (w *Watcher) Run() {
	if w.callback == nil {
		w.callback = func(*watch.TrainingInfo) {}
	}

	start := time.Now()

	for {
		select {
		case info := <-w.trainings:
			// use =="" stands for the case that the loop is done
			if info.User == nil {
				t := start.Add(w.interval)

				if n := time.Now(); t.After(n) {
					time.Sleep(t.Sub(n))
				}

				w.trainings <- trainingInfo{}

				start = time.Now()

			} else {
				changed := w.check(&info)
				w.log.Debugf("check training %s/%s", info.TrainingId, info.JobId)

				if info.isDone() {
					index := info.toIndex()

					if err := w.cli.SetTrainingInfo(&index, &info.result); err == nil {
						w.callback(&info.TrainingInfo)
					} else {
						w.log.Errorf("set training info failed, err:%s", err.Error())
						w.trainings <- info
					}

				} else {
					if changed {
						index := info.toIndex()
						if err := w.cli.SetTrainingInfo(&index, &info.result); err != nil {
							w.log.Errorf("set training info failed, err:%s", err.Error())
						}
					}

					w.trainings <- info
				}
			}

		case <-w.stop:
			close(w.stopped)

			return
		}
	}
}

func (w *Watcher) Exit() {
	close(w.stop)

	<-w.stopped

	w.cli.Disconnect()
}

func (w *Watcher) check(info *trainingInfo) (changed bool) {
	result := &info.result

	if !info.done {
		detail, err := w.ts.GetDetail(info.JobId)
		if err != nil || detail.Status.TrainingStatus() == result.Status {
			return
		}

		result.Status = detail.Status.TrainingStatus()
		result.Duration = detail.Duration

		changed = true

		if !detail.Status.IsDone() {
			return
		}

		info.done = true
		info.success = detail.Status.IsSuccess()
	}

	if !info.logDone {
		if v, err := w.ts.GetLogFilePath(info.LogDir); err != nil {
			w.log.Errorf("generate log failed, err:%s", err.Error())
		} else {
			result.LogPath = v
			info.logDone = true
			changed = true
		}
	}

	if !info.success {
		return
	}

	if !info.outputDone {
		if v, err := w.ts.GenOutput(info.OutputDir); err != nil {
			w.log.Errorf("generate output failed, err:%s", err.Error())
		} else {
			result.OutputZipPath = v
			info.outputDone = true
			changed = true
		}
	}

	if !info.aimDone {
		if v, err := w.ts.GenAim(info.AimDir); err != nil {
			w.log.Errorf("generate aim failed, err:%s", err.Error())
		} else {
			result.AimZipPath = v
			info.aimDone = true
			changed = true
		}
	}

	return
}
