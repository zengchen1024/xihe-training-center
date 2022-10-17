package syncrepo

import "github.com/opensourceways/xihe-training-center/domain"

type ProjectInfo struct {
	Name        domain.ProjectName
	Owner       domain.Account
	RepoId      string
	RepoURL     string
	StartCommit string
}

type SyncRepo interface {
	SyncProject(*ProjectInfo) (lastCommit string, err error)
	GetRepoSyncedCommit(*domain.ResourceRef) (c string, err error)
}
