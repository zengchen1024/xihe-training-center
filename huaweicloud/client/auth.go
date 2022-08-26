package client

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/chnsz/golangsdk"
	huaweisdk "github.com/chnsz/golangsdk/openstack"
)

func buildClient(c *Config) error {
	if c.AccessKey != "" && c.SecretKey != "" {
		return buildClientByAKSK(c)
	}

	return errors.New("can't build client")
}

func generateTLSConfig(c *Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: c.Insecure,
	}

	return config, nil
}

func genClient(c *Config, ao golangsdk.AuthOptionsProvider) (*golangsdk.ProviderClient, error) {
	client, err := huaweisdk.NewClient(ao.GetIdentityEndpoint())
	if err != nil {
		return nil, err
	}

	// Set UserAgent
	client.UserAgent.Prepend("huaweicloud-client")

	config, err := generateTLSConfig(c)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: config,
	}

	client.HTTPClient = http.Client{
		Transport: &LogRoundTripper{
			Rt:         transport,
			MaxRetries: c.MaxRetries,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if client.AKSKAuthOptions.AccessKey != "" {
				golangsdk.ReSign(req, golangsdk.SignOptions{
					AccessKey:  client.AKSKAuthOptions.AccessKey,
					SecretKey:  client.AKSKAuthOptions.SecretKey,
					RegionName: client.AKSKAuthOptions.Region,
				})
			}
			return nil
		},
	}

	if c.MaxRetries > 0 {
		client.MaxBackoffRetries = uint(c.MaxRetries)
		client.RetryBackoffFunc = retryBackoffFunc
	}

	// Validate authentication normally.
	err = huaweisdk.Authenticate(client, ao)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func genClients(c *Config, projectAuthOptions, domainAuthOptions golangsdk.AuthOptionsProvider) error {
	client, err := genClient(c, projectAuthOptions)
	if err != nil {
		return err
	}
	c.HwClient = client

	client, err = genClient(c, domainAuthOptions)
	if err == nil {
		c.DomainClient = client
	}
	return err
}

func buildClientByAKSK(c *Config) error {
	var projectAuthOptions, domainAuthOptions golangsdk.AKSKAuthOptions

	if c.AgencyDomainName != "" && c.AgencyName != "" {
		projectAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
			DelegatedProject: c.DelegatedProject,
		}

		domainAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID:         c.DomainID,
			Domain:           c.DomainName,
			AgencyName:       c.AgencyName,
			AgencyDomainName: c.AgencyDomainName,
		}
	} else {
		projectAuthOptions = golangsdk.AKSKAuthOptions{
			ProjectName: c.TenantName,
			ProjectId:   c.TenantID,
		}

		domainAuthOptions = golangsdk.AKSKAuthOptions{
			DomainID: c.DomainID,
			Domain:   c.DomainName,
		}
	}

	for _, ao := range []*golangsdk.AKSKAuthOptions{&projectAuthOptions, &domainAuthOptions} {
		ao.IdentityEndpoint = c.IdentityEndpoint
		ao.AccessKey = c.AccessKey
		ao.SecretKey = c.SecretKey
		if c.Region != "" {
			ao.Region = c.Region
		}
		if c.SecurityToken != "" {
			ao.SecurityToken = c.SecurityToken
			ao.WithUserCatalog = true
		}
	}
	return genClients(c, projectAuthOptions, domainAuthOptions)
}

func retryBackoffFunc(ctx context.Context, respErr *golangsdk.ErrUnexpectedResponseCode, e error, retries uint) error {
	minutes := int(math.Pow(2, float64(retries)))
	if minutes > 30 { // won't wait more than 30 minutes
		minutes = 30
	}

	log.Printf("[WARN] Received StatusTooManyRequests response code, try to sleep %d minutes", minutes)
	sleep := time.Duration(minutes) * time.Minute

	if ctx != nil {
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return e
		}
	} else {
		//lintignore:R018
		time.Sleep(sleep)
	}

	return nil
}
