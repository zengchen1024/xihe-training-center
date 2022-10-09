package trainingimpl

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"

	"github.com/opensourceways/xihe-training-center/domain"
	dt "github.com/opensourceways/xihe-training-center/domain/training"
	"github.com/opensourceways/xihe-training-center/huaweicloud/client"
	"github.com/opensourceways/xihe-training-center/huaweicloud/modelarts"
)

var statusMap = map[string]domain.TrainingStatus{
	"creating":    domain.NewStatusCreating(),
	"pending":     domain.NewStatusPending(),
	"running":     domain.NewStatusRunning(),
	"failed":      domain.NewStatusFailed(),
	"completed":   domain.NewStatusCompleted(),
	"terminating": domain.NewStatusTerminating(),
	"terminated":  domain.NewStatusTerminated(),
	"abnormal":    domain.NewStatusAbnormal(),
}

func NewTraining(cfg *Config) (dt.Training, error) {
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
		obsRepoPath: filepath.Join(cfg.OBSBucket, cfg.OBSRepoPath),
		config:      cfg.TrainingConfig,
	}, nil
}

type trainingImpl struct {
	cli         *golangsdk.ServiceClient
	obsRepoPath string

	config TrainingConfig
}

func (impl trainingImpl) Create(t *domain.UserTraining) (info dt.TrainingInfo, err error) {
	desc := ""
	if t.Desc != nil {
		desc = t.Desc.TrainingDesc()
	}

	obs := impl.obsFilePath(t.ToPath())
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	outputDir := filepath.Join(obs, impl.config.TrainOutputDir, timestamp) + "/"
	logDir := filepath.Join(obs, impl.config.TrainLogDir, timestamp) + "/"

	opt := modelarts.JobCreateOption{
		Kind: "job",
		Metadata: modelarts.MetadataOption{
			Name: t.Name.TrainingName(),
			Desc: desc,
		},
		Algorithm: modelarts.AlgorithmOption{
			CodeDir:  filepath.Join(obs, t.CodeDir.Directory()) + "/",
			BootFile: filepath.Join(obs, t.CodeDir.Directory(), t.BootFile.FilePath()),
			Engine: modelarts.EngineOption{
				EngineName:    t.Compute.Type.ComputeType(),
				EngineVersion: t.Compute.Version.ComputeVersion(),
			},
			Outputs: []modelarts.InputOutputOption{
				{
					Name: impl.config.TrainOutputKey,
					Remote: modelarts.RemoteOption{
						OBS: modelarts.OBSOption{
							OBSURL: outputDir,
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
				OBSURL: logDir,
			},
		},
	}

	if len(t.Inputs) > 0 {
		opt.Algorithm.Inputs = impl.genInputOption(t.Inputs)
	}

	if n := len(t.Hypeparameters); n > 0 {
		p := make([]modelarts.ParameterOption, n)
		for i, v := range t.Hypeparameters {
			p[i] = modelarts.ParameterOption{
				Name:  v.Key.CustomizedKey(),
				Value: v.Value.CustomizedValue(),
			}
		}

		opt.Algorithm.Parameters = p
	}

	if n := len(t.Env); n > 0 {
		m := make(map[string]string)
		for _, v := range t.Env {
			m[v.Key.CustomizedKey()] = v.Value.CustomizedValue()
		}

		opt.Algorithm.Environments = m
	}

	info.Id, err = modelarts.CreateJob(impl.cli, opt)
	if err == nil {
		info.OutputDir = outputDir
		info.LogDir = logDir
	}

	return
}

func (impl trainingImpl) genInputOption(kv []domain.Input) []modelarts.InputOutputOption {
	r := make([]modelarts.InputOutputOption, len(kv))

	for i, v := range kv {
		r[i] = modelarts.InputOutputOption{
			Name: v.Key.CustomizedKey(),
			Remote: modelarts.RemoteOption{
				OBS: modelarts.OBSOption{
					OBSURL: impl.obsFilePath(v.Value.ToPath()),
				},
			},
		}
	}

	return r
}

func (impl trainingImpl) obsFilePath(p string) string {
	return filepath.Join(impl.obsRepoPath, p)
}

func (impl trainingImpl) Delete(jobId string) error {
	return modelarts.DeleteJob(impl.cli, jobId)
}

func (impl trainingImpl) Get(jobId string) (r domain.TrainingDetail, err error) {
	v, err := modelarts.GetJob(impl.cli, jobId)
	if err != nil {
		return
	}

	if status, ok := statusMap[strings.ToLower(v.Status.Phase)]; ok {
		r.Status = status
	} else {
		r.Status = domain.NewStatusAbnormal()
	}

	r.Duration, err = domain.NewTrainingDuration(v.Status.Duration)

	return
}

func (impl trainingImpl) Terminate(jobId string) error {
	return modelarts.TerminateJob(impl.cli, jobId)
}

func (impl trainingImpl) GetLogURL(jobId string) (string, error) {
	return modelarts.GetLogURL(impl.cli, jobId)
}
