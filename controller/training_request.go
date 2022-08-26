package controller

import (
	"github.com/opensourceways/xihe-training-center/app"
	"github.com/opensourceways/xihe-training-center/domain"
)

type trainingLogResp struct {
	LogURL string `json:"job_id"`
}

type trainingCreateResp struct {
	JobId string `json:"job_id"`
}

type trainingCreateRequest struct {
	User      string `json:"user"`
	ProjectId string `json:"project_id"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	CodeDir  string `json:"code_dir"`
	BootFile string `json:"boot_file"`
	LogDir   string `json:"log_dir"`

	Hypeparameters []keyValue `json:"hyperparameter"`
	Env            []keyValue `json:"evn"`
	Inputs         []keyValue `json:"inputs"`
	Outputs        []keyValue `json:"outputs"`

	Compute struct {
		Type    string `json:"type"`
		Version string `json:"version"`
		Flavor  string `json:"flavor"`
	} `json:"compute"`
}

type keyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (kv *keyValue) toKeyValue() (r domain.KeyValue, err error) {
	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}

func (req *trainingCreateRequest) toCmd() (cmd app.TrainingCreateCmd, err error) {
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

	c := &req.Compute

	if cmd.Compute.Type, err = domain.NewComputeType(c.Type); err != nil {
		return
	}

	if cmd.Compute.Version, err = domain.NewComputeVersion(c.Version); err != nil {
		return
	}

	if cmd.Compute.Flavor, err = domain.NewComputeFlavor(c.Flavor); err != nil {
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

func (req *trainingCreateRequest) toKeyValue(kv []keyValue) (r []domain.KeyValue, err error) {
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
