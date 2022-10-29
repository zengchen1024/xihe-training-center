package domain

import (
	"path/filepath"
)

type UserTraining struct {
	User Account

	TrainingConfig
}

func (t *UserTraining) ToPath() string {
	return filepath.Join(
		t.User.Account(),
		ResourceTypeProject.ResourceType(), t.ProjectRepoId,
	)
}

type TrainingConfig struct {
	ProjectName   ProjectName
	ProjectRepoId string

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
	Key CustomizedKey
	ResourceRef
}

type ResourceRef struct {
	User   Account
	Type   ResourceType
	RepoId string
	File   string
}

func (r *ResourceRef) ToPath() string {
	s := filepath.Join(
		r.User.Account(), r.Type.ResourceType(),
		r.RepoId, r.File,
	)

	if r.File == "" {
		// The input is the directory. Appending "/" to make sure
		// the path is a directory for object storage service.
		return s + "/"
	}

	return s
}

type JobDetail struct {
	Status   TrainingStatus
	Duration int
}

type JobInfo struct {
	JobId     string
	LogDir    string
	AimDir    string
	OutputDir string
}
