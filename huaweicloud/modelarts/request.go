package modelarts

import "github.com/chnsz/golangsdk"

type JobCreateOption struct {
	Kind      string          `json:"kind" required:"true"`
	Metadata  MetadataOption  `json:"metadata" required:"true"`
	Algorithm AlgorithmOption `json:"algorithm"`
	Spec      SpecOption      `json:"spec"`
}

func (opts *JobCreateOption) toMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

type MetadataOption struct {
	Name string `json:"name" required:"true"`
	Desc string `json:"description"`
}

type AlgorithmOption struct {
	CodeDir      string              `json:"code_dir"`
	BootFile     string              `json:"boot_file"`
	Engine       EngineOption        `json:"engine"`
	Parameters   []ParameterOption   `json:"parameters"`
	Environments map[string]string   `json:"environments"`
	Inputs       []InputOutputOption `json:"inputs"`
	Outputs      []InputOutputOption `json:"outputs"`
}

type EngineOption struct {
	EngineName    string `json:"engine_name"`
	EngineVersion string `json:"engine_version"`
}

type ParameterOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InputOutputOption struct {
	Name   string       `json:"name" required:"true"`
	Remote RemoteOption `json:"remote" required:"true"`
}

type RemoteOption struct {
	OBS OBSOption `json:"obs" required:"true"`
}

type OBSOption struct {
	OBSURL string `json:"obs_url" required:"true"`
}

type SpecOption struct {
	Resource      ResourceOption      `json:"resource"`
	LogExportPath LogExportPathOption `json:"log_export_path"`
}

type ResourceOption struct {
	FlavorId  string `json:"flavor_id" Required:"true"`
	NodeCount int    `json:"node_count,omitempty"`
}

type LogExportPathOption struct {
	OBSURL string `json:"obs_url,omitempty"`
}

func CreateJob(client *golangsdk.ServiceClient, opts JobCreateOption) (string, error) {
	reqBody, err := opts.toMap()
	if err != nil {
		return "", err
	}

	r := golangsdk.Result{}
	_, r.Err = client.Post(
		createURL(client), reqBody, &r.Body,
		&golangsdk.RequestOpts{OkCodes: []int{201}},
	)

	job := new(Job)
	if err := r.ExtractInto(job); err != nil {
		return "", err
	}

	return job.Metadata.Id, nil
}

func DeleteJob(client *golangsdk.ServiceClient, jobId string) error {
	v, err := client.Delete(
		jobURL(client, jobId),
		&golangsdk.RequestOpts{OkCodes: []int{202}},
	)

	if err != nil && v.StatusCode == 404 {
		err = nil
	}

	return err
}

func TerminateJob(client *golangsdk.ServiceClient, jobId string) error {
	_, err := client.Post(
		actionURL(client, jobId),
		map[string]interface{}{
			"action_type": "terminate",
		},
		nil, &golangsdk.RequestOpts{OkCodes: []int{201, 202}},
	)

	return err
}

func GetJob(client *golangsdk.ServiceClient, jobId string) (j Job, err error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(
		jobURL(client, jobId), &r.Body,
		&golangsdk.RequestOpts{OkCodes: []int{200}},
	)

	err = r.ExtractInto(&j)

	return j, err
}

func GetLogDownloadURL(client *golangsdk.ServiceClient, jobId string) (string, error) {
	r := golangsdk.Result{}
	_, r.Err = client.Get(
		logURL(client, jobId), &r.Body,
		&golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/octet-stream",
			},
			OkCodes: []int{200},
		},
	)

	var v struct {
		OBSURL string `json:"obs_url"`
	}
	err := r.ExtractInto(&v)

	return v.OBSURL, err
}
