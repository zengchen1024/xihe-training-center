package training

type HuaweiCloud struct {
	Account   string `json:"account" required:"true"`
	AccountId string `json:"account_id" required:"true"`
	AccessKey string `json:"access_key" required:"true"`
	SecretKey string `json:"secret_key" required:"true"`
	Endpoint  string `json:"endpoint" required:"true"`
	OBSBucket string `json:"obs_bucket" required:"true"`
}
