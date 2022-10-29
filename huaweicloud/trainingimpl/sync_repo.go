package trainingimpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	libutils "github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/training"
	"github.com/opensourceways/xihe-training-center/utils"
)

func newHelper(cfg *Config) (*helper, error) {
	obsCfg := &cfg.OBS
	cli, err := obs.New(obsCfg.AccessKey, obsCfg.SecretKey, obsCfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("new obs client failed, err:%s", err.Error())
	}

	suc := &cfg.SyncAndUpload

	_, err, _ = libutils.RunCmd(
		suc.OBSUtilPath, "config",
		"-i="+obsCfg.AccessKey, "-k="+obsCfg.SecretKey, "-e="+obsCfg.Endpoint,
	)
	if err != nil {
		return nil, fmt.Errorf("obsutil config failed, err:%s", err.Error())
	}

	if err := os.Mkdir(suc.WorkDir, 0755); err != nil {
		return nil, err
	}

	if err := os.Mkdir(suc.UploadWorkDir, 0755); err != nil {
		return nil, err
	}

	return &helper{
		obsClient: cli,
		bucket:    obsCfg.Bucket,
		suc:       *suc,
	}, nil
}

type helper struct {
	obsClient *obs.ObsClient
	bucket    string
	suc       SyncAndUploadConfig
}

func (s *helper) GetRepoSyncedCommit(i *domain.ResourceRef) (
	c string, err error,
) {
	p := filepath.Join(s.suc.RepoPath, i.ToPath(), s.suc.CommitFile)

	err = utils.Retry(func() error {
		v, err := s.getObject(p)
		if err == nil && len(v) > 0 {
			c = string(v)
		}

		return err
	})

	return
}

func (s *helper) getObject(path string) ([]byte, error) {
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

func (s *helper) SyncProject(repo *training.ProjectInfo) (lastCommit string, err error) {
	cfg := &s.suc

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
