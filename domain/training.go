package domain

import "path/filepath"

type UserTraining struct {
	User Account

	Training
}

func (t *UserTraining) ToPath() string {
	return filepath.Join(
		t.User.Account(),
		ResourceTypeProject.ResourceType(), t.ProjectId,
	)
}

type Training struct {
	ProjectId   string
	ProjectName ProjectName

	Name TrainingName
	Desc TrainingDesc

	CodeDir  Directory
	BootFile FilePath

	Hypeparameters []KeyValue
	Env            []KeyValue
	Inputs         []Input

	Compute Compute
}

type Compute struct {
	Type    ComputeType
	Version ComputeVersion
	Flavor  ComputeFlavor
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type Input struct {
	Key   CustomizedKey
	Value ResourceInput
}

type ResourceInput struct {
	User   Account
	Type   ResourceType
	RepoId string
	File   string
}

func (r *ResourceInput) ToPath() string {
	s := filepath.Join(
		r.User.Account(), r.Type.ResourceType(),
		r.RepoId, r.File,
	)

	if r.File == "" {
		return s + "/"
	}

	return s
}

type TrainingDetail struct {
	Status   TrainingStatus
	Duration TrainingDuration
}
