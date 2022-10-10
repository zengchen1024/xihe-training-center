package syncrepoimpl

import (
	"errors"
	"path/filepath"
)

type Config struct {
	AccessKey string `json:"access_key"    required:"true"`
	SecretKey string `json:"secret_key"    required:"true"`
	Endpoint  string `json:"endpoint"      required:"true"`
	Bucket    string `json:"bucket"        required:"true"`

	SyncConfig
}

func (c *Config) Validate() error {
	return c.SyncConfig.validate()
}

type SyncConfig struct {
	WorkDir       string `json:"work_dir"        required:"true"`
	RepoPath      string `json:"repo_path"       required:"true"`
	CommitFile    string `json:"commit_file"     required:"true"`
	OBSUtilPath   string `json:"obsutil_path"    required:"true"`
	SyncFileShell string `json:"sync_file_shell" required:"true"`
}

func (c *SyncConfig) validate() error {
	if !filepath.IsAbs(c.OBSUtilPath) {
		return errors.New("obsutil_path must be an absolute path")
	}

	if !filepath.IsAbs(c.WorkDir) {
		return errors.New("work_dir must be an absolute path")
	}

	if !filepath.IsAbs(c.SyncFileShell) {
		return errors.New("sync_file_shell must be an absolute path")
	}

	if filepath.IsAbs(c.RepoPath) {
		return errors.New("repo_path can't start with /")
	}

	return nil
}
