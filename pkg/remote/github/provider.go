package github

import (
	"os"

	"github.com/cloudskiff/driftctl/pkg/output"
	"github.com/hashicorp/terraform/addrs"

	"github.com/cloudskiff/driftctl/pkg/remote/terraform"
	tf "github.com/cloudskiff/driftctl/pkg/terraform"
)

type GithubTerraformProvider struct {
	*terraform.TerraformProvider
	name    string
	version string
	address *addrs.Provider
}

type githubConfig struct {
	Token        string
	Owner        string `cty:"owner"`
	Organization string
}

func NewGithubTerraformProvider(version string, progress output.Progress, configDir string) (*GithubTerraformProvider, error) {
	if version == "" {
		version = "4.4.0"
	}
	p := &GithubTerraformProvider{
		version: version,
		name:    "github",
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
		Name:         p.name,
		DefaultAlias: p.GetConfig().getDefaultOwner(),
		GetProviderConfig: func(owner string) interface{} {
			return githubConfig{
				Owner: p.GetConfig().getDefaultOwner(),
			}
		},
		Addr: p.Address(),
	}, progress)
	if err != nil {
		return nil, err
	}
	p.TerraformProvider = tfProvider
	return p, err
}

func (c githubConfig) getDefaultOwner() string {
	if c.Organization != "" {
		return c.Organization
	}
	return c.Owner
}

func (p GithubTerraformProvider) GetConfig() githubConfig {
	return githubConfig{
		Token:        os.Getenv("GITHUB_TOKEN"),
		Owner:        os.Getenv("GITHUB_OWNER"),
		Organization: os.Getenv("GITHUB_ORGANIZATION"),
	}
}

func (p *GithubTerraformProvider) Name() string {
	return p.name
}

func (p *GithubTerraformProvider) Version() string {
	return p.version
}

func (p *GithubTerraformProvider) Address() *addrs.Provider {
	return p.address
}
