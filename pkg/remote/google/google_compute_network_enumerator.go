package google

import (
	"strings"

	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/remote/google/repository"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/google"
)

type GoogleComputeNetworkEnumerator struct {
	repository repository.AssetRepository
	factory    resource.ResourceFactory
}

func NewGoogleComputeNetworkEnumerator(repo repository.AssetRepository, factory resource.ResourceFactory) *GoogleComputeNetworkEnumerator {
	return &GoogleComputeNetworkEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *GoogleComputeNetworkEnumerator) SupportedType() resource.ResourceType {
	return google.GoogleComputeNetworkResourceType
}

func (e *GoogleComputeNetworkEnumerator) Enumerate() ([]*resource.Resource, error) {
	resources, err := e.repository.SearchAllNetworks()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, len(resources))

	for _, res := range resources {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				strings.Replace(res.Name, "//compute.googleapis.com/", "", 1),
				map[string]interface{}{
					"display_name": res.DisplayName,
				},
			),
		)
	}

	return results, err
}
