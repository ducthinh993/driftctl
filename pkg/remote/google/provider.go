package google

import (
	"os"

	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/google/config"
	"github.com/cloudskiff/driftctl/pkg/remote/terraform"
	tf "github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/hashicorp/terraform/addrs"
)

type GCPTerraformProvider struct {
	*terraform.TerraformProvider
	name    string
	version string
	address *addrs.Provider
}

func NewGCPTerraformProvider(version string, progress output.Progress, configDir string) (*GCPTerraformProvider, error) {
	if version == "" {
		version = "3.78.0"
	}
	p := &GCPTerraformProvider{
		version: version,
		name:    tf.GOOGLE,
		address: &addrs.Provider{
			Hostname:  "registry.terraform.io",
			Namespace: "hashicorp",
			Type:      "google",
		},
	}
	installer := tf.NewProviderInstaller(tf.ProviderConfig{
		Key:       p.name,
		Version:   version,
		ConfigDir: configDir,
	})
	tfProvider, err := terraform.NewTerraformProvider(installer, terraform.TerraformProviderConfig{
		Name: p.name,
		GetProviderConfig: func(alias string) interface{} {
			return p.GetConfig()
		},
		Addr: p.Address(),
	}, progress)

	if err != nil {
		return nil, err
	}

	p.TerraformProvider = tfProvider

	return p, err
}

func (p *GCPTerraformProvider) Name() string {
	return p.name
}

func (p *GCPTerraformProvider) Version() string {
	return p.version
}

func (p *GCPTerraformProvider) GetConfig() config.GCPTerraformConfig {
	return config.GCPTerraformConfig{
		Project: os.Getenv("CLOUDSDK_CORE_PROJECT"),
		Region:  os.Getenv("CLOUDSDK_COMPUTE_REGION"),
		Zone:    os.Getenv("CLOUDSDK_COMPUTE_ZONE"),
	}
}

func (p *GCPTerraformProvider) Address() *addrs.Provider {
	return p.address
}
