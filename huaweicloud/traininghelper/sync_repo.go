package traininghelper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	libutils "github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/syncrepo"
	"github.com/opensourceways/xihe-training-center/utils"
)

func NewTrainingHelper(cfg *Config) (*TrainingHelper, error) {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("new obs client failed, err:%s", err.Error())
	}

	_, err, _ = libutils.RunCmd(
		cfg.SyncConfig.OBSUtilPath, "config",
		"-i="+cfg.AccessKey, "-k="+cfg.SecretKey, "-e="+cfg.Endpoint,
	)
	if err != nil {
		return nil, fmt.Errorf("obsutil config failed, err:%s", err.Error())
	}

	if err := os.Mkdir(cfg.SyncConfig.WorkDir, 0755); err != nil {
		return nil, err
	}

	if err := os.Mkdir(cfg.UploadConfig.WorkDir, 0755); err != nil {
		return nil, err
	}

	return &TrainingHelper{
		obsClient: cli,
		bucket:    cfg.Bucket,
		config:    cfg.SyncConfig,
		uploadCfg: cfg.UploadConfig,
	}, nil
}

type TrainingHelper struct {
	obsClient *obs.ObsClient
	bucket    string
	config    SyncConfig
	uploadCfg UploadConfig
}

func (s *TrainingHelper) GetRepoSyncedCommit(i *domain.ResourceRef) (
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

func (s *TrainingHelper) getObject(path string) ([]byte, error) {
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

func (s *TrainingHelper) SyncProject(repo *syncrepo.ProjectInfo) (lastCommit string, err error) {
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
	)

	params := []string{
		cfg.SyncFileShell, tempDir,
		repo.RepoURL, repo.Name.ProjectName(),
		cfg.OBSUtilPath, s.bucket, obsRepoPath,
		repo.StartCommit,
	}

	v, err, _ := libutils.RunCmd(params...)
	if err != nil {
		err = fmt.Errorf(
			"run sync shell, err=%s, params=%v",
			err.Error(), params,
		)

		return
	}

	lastCommit = strings.TrimSuffix(string(v), "\n")

	return
}