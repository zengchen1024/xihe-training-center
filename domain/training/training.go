package training

import (
	"path/filepath"

	"github.com/opensourceways/xihe-training-center/domain"
)

type ProjectInfo struct {
	Name        domain.ProjectName
	Owner       domain.Account
	RepoId      string
	RepoURL     string
	StartCommit string
}

func (p *ProjectInfo) ToRepoPath() string {
	return filepath.Join(
		p.Owner.Account(),
		domain.ResourceTypeProject.ResourceType(), p.RepoId,
	)
}

type Training interface {
	Create(*domain.UserTraining) (domain.JobInfo, error)
	Delete(string) error
	Terminate(string) error
	GetLogDownloadURL(string) (string, error)
	GetDetail(string) (domain.JobDetail, error)

	// GetLogFilePath return the obs path of log
	GetLogFilePath(logDir string) (string, error)

	// GenOutput generates the zip file of output dir and
	// return the obs path of that file.
	GenOutput(outputDir string) (string, error)

	// GenAim generates the zip file of aim dir
	// and return the obs path of that file.
	GenAim(aimDir string) (string, error)

	SyncProject(*ProjectInfo) (lastCommit string, err error)
	GetRepoSyncedCommit(*domain.ResourceRef) (c string, err error)
}
