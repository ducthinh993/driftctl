package azurerm

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitResourcesMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	initAzureVirtualNetworkMetaData(resourceSchemaRepository)
	initAzureRouteTableMetaData(resourceSchemaRepository)
	initAzureRouteMetaData(resourceSchemaRepository)
	initAzureResourceGroupMetadata(resourceSchemaRepository)
	initAzureContainerRegistryMetadata(resourceSchemaRepository)
	initAzureFirewallMetadata(resourceSchemaRepository)
	initAzurePostgresqlServerMetadata(resourceSchemaRepository)
	initAzurePublicIPMetadata(resourceSchemaRepository)
	initAzurePostgresqlDatabaseMetadata(resourceSchemaRepository)
	initAzureNetworkSecurityGroupMetadata(resourceSchemaRepository)
	initAzureLoadBalancerMetadata(resourceSchemaRepository)
	initAzurePrivateDNSZoneMetaData(resourceSchemaRepository)
	initAzureImageMetaData(resourceSchemaRepository)
}
