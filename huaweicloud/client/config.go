package client

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chnsz/golangsdk"
)

type Config struct {
	AccessKey           string
	SecretKey           string
	CACertFile          string
	ClientCertFile      string
	ClientKeyFile       string
	DomainID            string
	DomainName          string
	IdentityEndpoint    string
	Insecure            bool
	Region              string
	TenantID            string
	TenantName          string
	Token               string
	SecurityToken       string
	AssumeRoleAgency    string
	AssumeRoleDomain    string
	Cloud               string
	MaxRetries          int
	RegionClient        bool
	EnterpriseProjectID string
	SharedConfigFile    string
	Profile             string

	// metadata security key expires at
	SecurityKeyExpiresAt time.Time

	HwClient     *golangsdk.ProviderClient
	DomainClient *golangsdk.ProviderClient

	// the custom endpoints used to override the default endpoint URL
	Endpoints map[string]string

	// RegionProjectIDMap is a map which stores the region-projectId pairs,
	// and region name will be the key and projectID will be the value in this map.
	RegionProjectIDMap map[string]string

	// RPLock is used to make the accessing of RegionProjectIDMap serial,
	// prevent sending duplicate query requests
	RPLock *sync.Mutex

	// SecurityKeyLock is used to make the accessing of SecurityKeyExpiresAt serial,
	// prevent sending duplicate query metadata api
	SecurityKeyLock *sync.Mutex

	// Legacy
	Username         string
	UserID           string
	Password         string
	AgencyName       string
	AgencyDomainName string
	DelegatedProject string
}

func (c *Config) LoadAndValidate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries should be a positive value")
	}

	if err := buildClient(c); err != nil {
		return err
	}

	if c.Region == "" {
		return fmt.Errorf("region should be provided")
	}

	return nil
}

// NewServiceClient create a ServiceClient which was assembled from ServiceCatalog.
// If you want to add new ServiceClient, please make sure the catalog was already in allServiceCatalog.
// the endpoint likes https://{Name}.{Region}.myhuaweicloud.com/{Version}/{project_id}/{ResourceBase}
func (c *Config) NewServiceClient(srv string, sc ServiceCatalog) (*golangsdk.ServiceClient, error) {
	client := c.HwClient

	if endpoint, ok := c.Endpoints[srv]; ok {
		return c.newServiceClientByEndpoint(client, endpoint, sc)
	}

	return nil, errors.New("can't new service client")
}

// newServiceClientByEndpoint returns a ServiceClient which the endpoint was initialized by customer
// the format of customer endpoint likes https://{Name}.{Region}.xxxx.com
func (c *Config) newServiceClientByEndpoint(
	client *golangsdk.ProviderClient,
	endpoint string, catalog ServiceCatalog,
) (*golangsdk.ServiceClient, error) {
	e := strings.TrimSuffix(endpoint, "/")
	e += "/"

	sc := &golangsdk.ServiceClient{
		ProviderClient: client,
		Endpoint:       e,
	}

	sc.ResourceBase = sc.Endpoint
	if catalog.Version != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.Version + "/"
	}
	if !catalog.WithOutProjectID {
		sc.ResourceBase = sc.ResourceBase + client.ProjectID + "/"
	}
	if catalog.ResourceBase != "" {
		sc.ResourceBase = sc.ResourceBase + catalog.ResourceBase + "/"
	}

	return sc, nil
}
