package trainingimpl

type Config struct {
	AccessKey   string `json:"access_key" required:"true"`
	SecretKey   string `json:"secret_key" required:"true"`
	Region      string `json:"region" required:"true"`
	ProjectName string `json:"project_name" required:"true"`
	ProjectId   string `json:"project_id" required:"true"`

	// modelarts endpoint
	Endpoint string `json:"endpoint" required:"true"`

	OBSBucket   string `json:"obs_bucket" required:"true"`
	OBSRepoPath string `json:"obs_repo_path" required:"true"`

	TrainingConfig
}

func (cfg *Config) SetDefault() {
	cfg.TrainingConfig.setDefault()
}

type TrainingConfig struct {
	TrainOutputKey string `json:"train_output_key"`
	TrainOutputDir string `json:"train_output_dir"`
	TrainLogDir    string `json:"train_log_dir"`
}

func (cfg *TrainingConfig) setDefault() {
	cfg.TrainOutputKey = "output_url"
	cfg.TrainOutputDir = "train-output"
	cfg.TrainLogDir = "train-log"
}
