package training

type HuaweiCloud struct {
	AccessKey   string `json:"access_key" required:"true"`
	SecretKey   string `json:"secret_key" required:"true"`
	Region      string `json:"region" required:"true"`
	Endpoint    string `json:"endpoint" required:"true"`
	ProjectName string `json:"project_name" required:"true"`
	ProjectId   string `json:"project_id" required:"true"`
	OBSBucket   string `json:"obs_bucket" required:"true"`
}
