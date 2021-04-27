package aws

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitResourcesMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	initAwsAmiMetaData(resourceSchemaRepository)
	initAwsCloudfrontDistributionMetaData(resourceSchemaRepository)
	initAwsDbInstanceMetaData(resourceSchemaRepository)
	initAwsDbSubnetGroupMetaData(resourceSchemaRepository)
	initAwsDefaultSecurityGroupMetaData(resourceSchemaRepository)
	initAwsDefaultSubnetMetaData(resourceSchemaRepository)
	initAwsDynamodbTableMetaData(resourceSchemaRepository)
	initAwsEbsSnapshotMetaData(resourceSchemaRepository)
}
