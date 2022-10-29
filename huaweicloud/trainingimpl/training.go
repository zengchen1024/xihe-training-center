package trainingimpl

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"

	"github.com/opensourceways/xihe-training-center/domain"
	"github.com/opensourceways/xihe-training-center/domain/training"
	"github.com/opensourceways/xihe-training-center/huaweicloud/client"
	"github.com/opensourceways/xihe-training-center/huaweicloud/modelarts"
)

const obsPrefix = "obs://"

var statusMap = map[string]domain.TrainingStatus{
	"failed":      domain.TrainingStatusFailed,
	"pending":     domain.TrainingStatusRunning,
	"running":     domain.TrainingStatusRunning,
	"creating":    domain.TrainingStatusRunning,
	"abnormal":    domain.TrainingStatusFailed,
	"completed":   domain.TrainingStatusCompleted,
	"terminated":  domain.TrainingStatusTerminated,
	"terminating": domain.TrainingStatusTerminated,
}

func NewTraining(cfg *Config) (training.Training, error) {
	s := "modelarts"
	v := client.Config{
		AccessKey:  cfg.AccessKey,
		SecretKey:  cfg.SecretKey,
		TenantName: cfg.ProjectName,
		TenantID:   cfg.ProjectId,
		Region:     cfg.Region,
		Endpoints: map[string]string{
			s: cfg.Endpoint,
		},
		IdentityEndpoint: fmt.Sprintf("https://iam.%s.myhuaweicloud.com:443/v3", cfg.Region),
	}
	if err := v.LoadAndValidate(); err != nil {
		return nil, err
	}

	cli, err := v.NewServiceClient(s, client.ServiceCatalog{
		Version: "v2",
	})
	if err != nil {
		return nil, err
	}

	return trainingImpl{
		cli:         cli,
		config:      cfg.TrainingConfig,
		obsRepoPath: filepath.Join(cfg.OBSBucket, cfg.OBSRepoPath),
	}, nil
}

type trainingImpl struct {
	cli         *golangsdk.ServiceClient
	config      TrainingConfig
	obsRepoPath string
}

func (impl trainingImpl) genJobParameter(t *domain.UserTraining, opt *modelarts.JobCreateOption) {
	if n := len(t.Hypeparameters); n > 0 {
		p := make([]modelarts.ParameterOption, n)

		for i, v := range t.Hypeparameters {
			s := ""
			if v.Value != nil {
				s = v.Value.CustomizedValue()
			}

			p[i] = modelarts.ParameterOption{
				Name:  v.Key.CustomizedKey(),
				Value: s,
			}
		}

		opt.Algorithm.Parameters = p
	}

	if n := len(t.Env); n > 0 {
		m := make(map[string]string)

		for _, v := range t.Env {
			s := ""
			if v.Value != nil {
				s = v.Value.CustomizedValue()
			}

			m[v.Key.CustomizedKey()] = s
		}

		opt.Algorithm.Environments = m
	}
}

func (impl trainingImpl) Create(t *domain.UserTraining) (info domain.JobInfo, err error) {
	desc := ""
	if t.Desc != nil {
		desc = t.Desc.TrainingDesc()
	}

	cfg := &impl.config
	obs := filepath.Join(impl.obsRepoPath, t.ToPath())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	info.LogDir = filepath.Join(obs, cfg.TrainLogDir, timestamp) + "/"
	info.OutputDir = filepath.Join(obs, cfg.TrainOutputDir, timestamp) + "/"

	opt := modelarts.JobCreateOption{
		Kind: "job",
		Metadata: modelarts.MetadataOption{
			Name: t.Name.TrainingName(),
			Desc: desc,
		},
		Algorithm: modelarts.AlgorithmOption{
			CodeDir:  obsPrefix + filepath.Join(obs, t.CodeDir.Directory()) + "/",
			BootFile: obsPrefix + filepath.Join(obs, t.CodeDir.Directory(), t.BootFile.FilePath()),
			Engine: modelarts.EngineOption{
				EngineName:    t.Compute.Type.ComputeType(),
				EngineVersion: t.Compute.Version.ComputeVersion(),
			},
			Outputs: []modelarts.InputOutputOption{
				{
					Name: cfg.TrainOutputKey,
					Remote: modelarts.RemoteOption{
						OBS: modelarts.OBSOption{
							OBSURL: obsPrefix + info.OutputDir,
						},
					},
				},
			},
		},
		Spec: modelarts.SpecOption{
			Resource: modelarts.ResourceOption{
				FlavorId:  t.Compute.Flavor.ComputeFlavor(),
				NodeCount: 1,
			},
			LogExportPath: modelarts.LogExportPathOption{
				OBSURL: obsPrefix + info.LogDir,
			},
		},
	}

	if len(t.Inputs) > 0 {
		opt.Algorithm.Inputs = impl.genInputOption(t.Inputs)
	}

	impl.genJobParameter(t, &opt)

	info.JobId, err = modelarts.CreateJob(impl.cli, opt)

	return
}

func (impl trainingImpl) genInputOption(kv []domain.Input) []modelarts.InputOutputOption {
	r := make([]modelarts.InputOutputOption, len(kv))

	for i, v := range kv {
		r[i] = modelarts.InputOutputOption{
			Name: v.Key.CustomizedKey(),
			Remote: modelarts.RemoteOption{
				OBS: modelarts.OBSOption{
					// v.Value maybe a directory.
					OBSURL: obsPrefix + impl.obsRepoPath + "/" + v.ToPath(),
				},
			},
		}
	}

	return r
}

func (impl trainingImpl) Delete(jobId string) error {
	return modelarts.DeleteJob(impl.cli, jobId)
}

func (impl trainingImpl) GetDetail(jobId string) (r domain.JobDetail, err error) {
	v, err := modelarts.GetJob(impl.cli, jobId)
	if err != nil {
		return
	}

	if status, ok := statusMap[strings.ToLower(v.Status.Phase)]; ok {
		r.Status = status
	} else {
		r.Status = domain.TrainingStatusFailed
	}

	r.Duration = v.Status.Duration

	return
}

func (impl trainingImpl) Terminate(jobId string) error {
	return modelarts.TerminateJob(impl.cli, jobId)
}

func (impl trainingImpl) GetLogDownloadURL(jobId string) (string, error) {
	return modelarts.GetLogDownloadURL(impl.cli, jobId)
}

// GetLogFilePath return the obs path of log
func (impl trainingImpl) GetLogFilePath(logDir string) (string, error) {
	return "", nil
}

// GenOutput generates the zip file of output dir and
// return the obs path of that file.
func (impl trainingImpl) GenOutput(outputDir string) (string, error) {
	return "", nil
}

// GenAim generates the zip file of aim dir
// and return the obs path of that file.
func (impl trainingImpl) GenAim(aimDir string) (string, error) {
	return "", nil
}
