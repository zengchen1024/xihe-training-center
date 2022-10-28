package watch

import (
	"errors"
	"sync"
	"time"

	"github.com/opensourceways/xihe-grpc-protocol/training/client"
)

type ErrorTooManyTrainings struct {
	error
}

func NewWatchService(cfg *Config) WatchService {
	n := cfg.MaxTrainingNum + 1
	return &watcher{
		maxTrainingNum:  n,
		intervalPerLoop: time.Duration(cfg.Interval) * time.Second,
		trainings:       make(chan TrainingInfo, n),
	}
}

type watcher struct {
	cli *client.Client
	ts  TrainingService

	lock           sync.RWMutex
	currentNum     int
	maxTrainingNum int

	intervalPerLoop time.Duration

	stop      chan struct{}
	stopped   chan struct{}
	trainings chan TrainingInfo
}

func (w *watcher) AddTraining(t *TrainingInfo) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.currentNum+1 >= w.maxTrainingNum {
		return ErrorTooManyTrainings{
			errors.New("too many trainings"),
		}
	}

	w.currentNum++

	w.trainings <- *t

	return nil
}

func (w *watcher) Run() {
	start := time.Now()

	for {
		select {
		case info := <-w.trainings:
			// TrainingInfo of which all fields are empty
			// stands for the tag specifying that the loop is done
			if info.User == "" {
				t := start.Add(w.intervalPerLoop)

				if n := time.Now(); t.After(n) {
					time.Sleep(t.Sub(n))
				}

				w.trainings <- TrainingInfo{}

				start = time.Now()

			} else {
				changed, done := w.do(&info)
				if !done {
					// TODO: check if timeout for this training. ignore it if timeout
				}

				if !done {
					w.trainings <- info
				}

				if changed || info.notifyFailed {
					//  send info.result, if send failed, info.notifyFailed = true
				}
			}

		case <-w.stop:
			close(w.stopped)

			return
		}
	}
}

func (w *watcher) Exit() {
	close(w.stop)

	<-w.stopped
}

type Data struct {
	Status   string
	Duration int

	Log    string
	Output string
	Aim    string
}

func (w *watcher) do(info *TrainingInfo) (changed, done bool) {
	data, changed := w.check(info)
	if !changed {
		return
	}

	done = w.ts.IsDone(data.Status) && data.Log != ""

	if done && w.ts.IsSucess(data.Status) {
		done = data.Output != "" && data.Aim != ""
	}

	if !done {
		info.result = data
	}

	return
}

func (w *watcher) check(info *TrainingInfo) (data Data, changed bool) {
	result := &info.result
	success := false

	if !w.ts.IsDone(result.Status) {
		detail, err := w.ts.GetDetail(info.JobId)
		if err != nil || detail.Status == result.Status {
			return
		}

		data.Status = detail.Status
		data.Duration = detail.Duration

		changed = true

		if !w.ts.IsDone(detail.Status) {
			return
		}

		success = w.ts.IsSucess(detail.Status)
	} else {
		data.Status = result.Status
		data.Duration = result.Duration

		success = w.ts.IsSucess(result.Status)
	}

	if result.Log == "" {
		if v, err := w.ts.GetLogFilePath(info.LogDir); err != nil {
			// log it
		} else {
			data.Log = v
			changed = true
		}
	}

	if !success {
		return
	}

	if result.Output == "" {
		if v, err := w.ts.GenOutput(info.OutputDir); err != nil {
			//log
		} else {
			data.Output = v
			changed = true
		}
	}

	if result.Aim == "" {
		if v, err := w.ts.GenAim(info.AimDir); err != nil {
			// log
		} else {
			data.Aim = v
			changed = true
		}
	}

	return
}
