package app

import (
	"errors"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/platform"
	"github.com/opensourceways/xihe-training-center/domain/synclock"
	"github.com/opensourceways/xihe-training-center/domain/syncrepo"
	"github.com/opensourceways/xihe-training-center/utils"
	"github.com/sirupsen/logrus"
)

func newSyncService(
	h syncrepo.SyncRepo,
	lock synclock.RepoSyncLock,
	p platform.Platform,
	log *logrus.Entry,
) *syncService {
	return &syncService{
		h:    h,
		p:    p,
		lock: lock,
		log:  log,
	}
}

type syncService struct {
	log  *logrus.Entry
	h    syncrepo.SyncRepo
	lock synclock.RepoSyncLock
	p    platform.Platform
}

func (s *syncService) checkResourceReady(i *domain.ResourceRef) error {
	c, err := s.h.GetRepoSyncedCommit(i)
	if err != nil {
		return err
	}

	if c == "" {
		return errors.New("not ready")
	}

	lastCommit, err := s.p.GetLastCommit(i.RepoId)
	if err != nil {
		return err
	}

	if string(c) != lastCommit {
		return errors.New("not ready")
	}

	return nil
}

func (s *syncService) syncProject(
	owner domain.Account, repoName domain.ProjectName, repoId string,
) (syncErr error) {
	c, err := s.lock.Find(owner, domain.ResourceTypeProject, repoId)
	if err != nil {
		if !synclock.IsRepoSyncLockNotExist(err) {
			return err
		}

		c.Owner = owner
		c.RepoId = repoId
		c.RepoType = domain.ResourceTypeProject
	}

	if c.Status != nil && !c.Status.IsDone() {
		return errors.New("can't sync")
	}

	lastCommit, err := s.p.GetLastCommit(repoId)
	if err != nil {
		return err
	}

	if c.LastCommit == lastCommit {
		return nil
	}

	// try lock
	c.Status = domain.RepoSyncStatusRunning
	c, err = s.lock.Save(&c)
	if err != nil {
		return err
	}

	// do sync
	repoURL := s.p.GetCloneURL(owner.Account(), repoName.ProjectName())
	lastCommit, syncErr = s.h.SyncProject(&syncrepo.ProjectInfo{
		Name:        repoName,
		Owner:       owner,
		RepoId:      repoId,
		RepoURL:     repoURL,
		StartCommit: c.LastCommit,
	})

	if syncErr == nil {
		c.LastCommit = lastCommit
	}
	c.Status = domain.RepoSyncStatusDone

	// unlock
	err = utils.Retry(func() error {
		_, err := s.lock.Save(&c)
		if err != nil {
			s.log.Errorf(
				"unlock sync repo failed, err:%s",
				err.Error(),
			)
		}

		return err
	})
	if err != nil {
		s.log.Errorf(
			"dead lock happened for repo: %s:%s",
			owner.Account(), repoId,
		)
	}

	return syncErr
}
