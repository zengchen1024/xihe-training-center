package controller

import (
	"github.com/opensourceways/xihe-training-center/app"
	"github.com/opensourceways/xihe-training-center/domain"
)

type TrainingLogResp struct {
	LogURL string `json:"log_url"`
}

type TrainingCreateResp struct {
	JobId string `json:"job_id"`
}

type TrainingCreateRequest struct {
	User      string `json:"user"`
	ProjectId string `json:"project_id"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	CodeDir  string `json:"code_dir"`
	BootFile string `json:"boot_file"`
	LogDir   string `json:"log_dir"`

	Hypeparameters []KeyValue `json:"hyperparameter"`
	Env            []KeyValue `json:"evn"`
	Inputs         []KeyValue `json:"inputs"`
	Outputs        []KeyValue `json:"outputs"`

	Compute Compute `json:"compute"`
}

type Compute struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Flavor  string `json:"flavor"`
}

func (c *Compute) toCompute() (r domain.Compute, err error) {
	if r.Type, err = domain.NewComputeType(c.Type); err != nil {
		return
	}

	if r.Version, err = domain.NewComputeVersion(c.Version); err != nil {
		return
	}

	if r.Flavor, err = domain.NewComputeFlavor(c.Flavor); err != nil {
		return
	}

	return
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (kv *KeyValue) toKeyValue() (r domain.KeyValue, err error) {
	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}

func (req *TrainingCreateRequest) toCmd() (cmd app.TrainingCreateCmd, err error) {
	if cmd.User, err = domain.NewAccount(req.User); err != nil {
		return
	}

	cmd.ProjectId = req.ProjectId

	if cmd.Name, err = domain.NewTrainingName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewTrainingDesc(req.Desc); err != nil {
		return
	}

	if cmd.CodeDir, err = domain.NewDirectory(req.CodeDir); err != nil {
		return
	}

	if cmd.BootFile, err = domain.NewFilePath(req.BootFile); err != nil {
		return
	}

	if cmd.LogDir, err = domain.NewDirectory(req.LogDir); err != nil {
		return
	}

	if cmd.Compute, err = req.Compute.toCompute(); err != nil {
		return
	}

	if cmd.Hypeparameters, err = req.toKeyValue(req.Hypeparameters); err != nil {
		return
	}

	if cmd.Env, err = req.toKeyValue(req.Env); err != nil {
		return
	}

	if cmd.Inputs, err = req.toKeyValue(req.Inputs); err != nil {
		return
	}

	if cmd.Outputs, err = req.toKeyValue(req.Outputs); err != nil {
		return
	}

	return
}

func (req *TrainingCreateRequest) toKeyValue(kv []KeyValue) (r []domain.KeyValue, err error) {
	n := len(kv)
	if n == 0 {
		return nil, nil
	}

	r = make([]domain.KeyValue, n)
	for i := range kv {
		if r[i], err = kv[i].toKeyValue(); err != nil {
			return
		}
	}

	return
}
