package trainingimpl

import (
	"errors"
	"path/filepath"
)

type configSetDefault interface {
	setDefault()
}

type configValidate interface {
	validate() error
}

type Config struct {
	OBS           OBSConfig           `json:"obs"         required:"true"`
	Train         TrainingConfig      `json:"train"       required:"true"`
	Modelarts     ModelartsConfig     `json:"modelarts"   required:"true"`
	SyncAndUpload SyncAndUploadConfig `json:"sync"        required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.OBS,
		&cfg.Train,
		&cfg.Modelarts,
		&cfg.SyncAndUpload,
	}
}

func (cfg *Config) Validate() error {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configValidate); ok {
			if err := v.validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configSetDefault); ok {
			v.setDefault()
		}
	}
}

type ModelartsConfig struct {
	AccessKey   string `json:"access_key" required:"true"`
	SecretKey   string `json:"secret_key" required:"true"`
	Region      string `json:"region" required:"true"`
	ProjectName string `json:"project_name" required:"true"`
	ProjectId   string `json:"project_id" required:"true"`

	// modelarts endpoint
	Endpoint string `json:"endpoint" required:"true"`
}

type TrainingConfig struct {
	LogDir    string `json:"log_dir"`
	AimKey    string `json:"aim_key"`
	AimDir    string `json:"aim_dir"`
	OutputKey string `json:"output_key"`
	OutputDir string `json:"output_dir"`
}

func (cfg *TrainingConfig) setDefault() {
	cfg.LogDir = "train-log"
	cfg.AimKey = "aim_repo"
	cfg.AimDir = "tain-aim"
	cfg.OutputKey = "output_path"
	cfg.OutputDir = "train-output"
}

type OBSConfig struct {
	AccessKey string `json:"access_key"    required:"true"`
	SecretKey string `json:"secret_key"    required:"true"`
	Endpoint  string `json:"endpoint"      required:"true"`
	Bucket    string `json:"bucket"        required:"true"`
}

type SyncAndUploadConfig struct {
	RepoPath      string `json:"repo_path"                required:"true"`
	CommitFile    string `json:"commit_file"              required:"true"`
	OBSUtilPath   string `json:"obsutil_path"             required:"true"`
	SyncWorkDir   string `json:"sync_work_dir"            required:"true"`
	SyncFileShell string `json:"sync_file_shell"          required:"true"`

	UploadWorkDir     string `json:"upload_work_dir"      required:"true"`
	UploadFolderShell string `json:"upload_folder_shell"  required:"true"`
}

func (c *SyncAndUploadConfig) validate() error {
	if !filepath.IsAbs(c.OBSUtilPath) {
		return errors.New("obsutil_path must be an absolute path")
	}

	if !filepath.IsAbs(c.SyncWorkDir) {
		return errors.New("sync_work_dir must be an absolute path")
	}

	if !filepath.IsAbs(c.SyncFileShell) {
		return errors.New("sync_file_shell must be an absolute path")
	}

	if filepath.IsAbs(c.RepoPath) {
		return errors.New("repo_path can't start with /")
	}

	if !filepath.IsAbs(c.UploadWorkDir) {
		return errors.New("upload_work_dir must be an absolute path")
	}

	if !filepath.IsAbs(c.UploadFolderShell) {
		return errors.New("upload_folder_shell must be an absolute path")
	}

	return nil
}
