package azurerm

import (
	"os"

	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/cloudskiff/driftctl/pkg/remote/azurerm/common"
	"github.com/cloudskiff/driftctl/pkg/remote/terraform"
	tf "github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/hashicorp/terraform/addrs"
)

type AzureTerraformProvider struct {
	*terraform.TerraformProvider
	name    string
	version string
	address *addrs.Provider
}

func NewAzureTerraformProvider(version string, progress output.Progress, configDir string) (*AzureTerraformProvider, error) {
	// Just pass your version and name
	p := &AzureTerraformProvider{
		version: version,
		name:    tf.AZURE,
		address: &addrs.Provider{
			Hostname:  "registry.terraform.io",
			Namespace: "hashicorp",
			Type:      "azurerm",
		},
	}
	// Use TerraformProviderInstaller to retrieve the provider if needed
	installer := tf.NewProviderInstaller(tf.ProviderConfig{
		Key:       p.name,
		Version:   version,
		ConfigDir: configDir,
	})

	tfProvider, err := terraform.NewTerraformProvider(installer, terraform.TerraformProviderConfig{
		Name: p.name,
		Addr: p.Address(),
	}, progress)
	if err != nil {
		return nil, err
	}
	p.TerraformProvider = tfProvider
	return p, err
}

func (p *AzureTerraformProvider) GetConfig() common.AzureProviderConfig {
	return common.AzureProviderConfig{
		SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
		TenantID:       os.Getenv("ARM_TENANT_ID"),
		ClientID:       os.Getenv("ARM_CLIENT_ID"),
		ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
	}
}

func (p *AzureTerraformProvider) Name() string {
	return p.name
}

func (p *AzureTerraformProvider) Version() string {
	return p.version
}

func (p *AzureTerraformProvider) Address() *addrs.Provider {
	return p.address
}
