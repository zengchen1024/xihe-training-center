package training

import (
	"fmt"
	"strings"

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

func NewTraining(cfg *HuaweiCloud) (dt.Training, error) {
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
		cli:       cli,
		obsBucket: cfg.OBSBucket,
	}, nil
}

type trainingImpl struct {
	cli       *golangsdk.ServiceClient
	obsBucket string
}

func (impl trainingImpl) Create(t *domain.UserTraining) (string, error) {
	desc := ""
	if t.Desc != nil {
		desc = t.Desc.TrainingDesc()
	}

	obs := fmt.Sprintf(
		"%s/projects/%s/%s",
		impl.obsBucket, t.User.Account(), t.ProjectId,
	)

	opt := modelarts.JobCreateOption{
		Kind: "job",
		Metadata: modelarts.MetadataOption{
			Name: t.Name.TrainingName(),
			Desc: desc,
		},
		Algorithm: modelarts.AlgorithmOption{
			CodeDir:  fmt.Sprintf("%s/%s", obs, t.CodeDir.Directory()),
			BootFile: fmt.Sprintf("%s/%s/%s", obs, t.CodeDir.Directory(), t.BootFile),
			Engine: modelarts.EngineOption{
				EngineName:    t.Compute.Type.ComputeType(),
				EngineVersion: t.Compute.Version.ComputeVersion(),
			},
		},
		Spec: modelarts.SpecOption{
			Resource: modelarts.ResourceOption{
				FlavorId:  t.Compute.Flavor.ComputeFlavor(),
				NodeCount: 1,
			},
			LogExportPath: modelarts.LogExportPathOption{
				OBSURL: fmt.Sprintf("%s/%s", obs, t.LogDir.Directory()),
			},
		},
	}

	if len(t.Inputs) > 0 {
		opt.Algorithm.Inputs = genInputOutputOption(obs, t.Inputs)
	}

	if len(t.Outputs) > 0 {
		opt.Algorithm.Outputs = genInputOutputOption(obs, t.Outputs)
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
	}

	return modelarts.CreateJob(impl.cli, opt)
}

func genInputOutputOption(obs string, kv []domain.KeyValue) []modelarts.InputOutputOption {
	r := make([]modelarts.InputOutputOption, len(kv))
	for i, v := range kv {
		r[i] = modelarts.InputOutputOption{
			Name: v.Key.CustomizedKey(),
			Remote: modelarts.RemoteOption{
				OBS: modelarts.OBSOption{
					OBSURL: fmt.Sprintf("%s/%s", obs, v.Value.CustomizedValue()),
				},
			},
		}
	}

	return r
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
