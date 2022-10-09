package syncrepoimpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	libutils "github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/syncrepo"
	"github.com/opensourceways/xihe-training-center/utils"
)

func NewSyncRepo(cfg *Config) (syncrepo.SyncRepo, error) {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("new obs client failed, err:%s", err.Error())
	}

	_, err, _ = libutils.RunCmd(
		cfg.OBSUtilPath, "config",
		"-i="+cfg.AccessKey, "-k="+cfg.SecretKey, "-e="+cfg.Endpoint,
	)
	if err != nil {
		return nil, fmt.Errorf("obsutil config failed, err:%s", err.Error())
	}

	return &syncRepoImpl{
		obsClient: cli,
		bucket:    cfg.Bucket,
		config:    cfg.SyncConfig,
	}, nil
}

type syncRepoImpl struct {
	obsClient *obs.ObsClient
	bucket    string
	config    SyncConfig
}

func (s *syncRepoImpl) GetRepoSyncedCommit(i *domain.ResourceInput) (
	c string, err error,
) {
	p := filepath.Join(s.config.RepoPath, i.ToPath(), s.config.CommitFile)

	err = utils.Retry(func() error {
		v, err := s.getObject(p)
		if err == nil && len(v) > 0 {
			c = string(v)
		}

		return err
	})

	return
}

func (s *syncRepoImpl) getObject(path string) ([]byte, error) {
	input := &obs.GetObjectInput{}
	input.Bucket = s.bucket
	input.Key = path

	output, err := s.obsClient.GetObject(input)
	if err != nil {
		v, ok := err.(obs.ObsError)
		if ok && v.BaseModel.StatusCode == 404 {
			return nil, nil
		}

		return nil, err
	}

	v, err := ioutil.ReadAll(output.Body)

	output.Body.Close()

	return v, err
}

func (s *syncRepoImpl) SyncProject(repo *syncrepo.ProjectInfo) (lastCommit string, err error) {
	cfg := &s.config

	tempDir, err := ioutil.TempDir(cfg.WorkDir, "sync")
	if err != nil {
		return
	}

	defer os.RemoveAll(tempDir)

	obsRepoPath := filepath.Join(
		cfg.RepoPath,
		repo.Owner.Account(),
		domain.ResourceTypeProject.ResourceType(), repo.RepoId,
	) + "/"

	v, err, _ := libutils.RunCmd(
		cfg.SyncFileShell, tempDir,
		repo.RepoURL, repo.Name.ProjectName(),
		repo.StartCommit, cfg.OBSUtilPath, obsRepoPath,
	)
	if err != nil {
		return
	}

	lastCommit = string(v)

	return
}
