package aws

import "github.com/cloudskiff/driftctl/pkg/resource"

func InitResourcesMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	initAwsAmiMetaData(resourceSchemaRepository)
	initAwsCloudfrontDistributionMetaData(resourceSchemaRepository)
	initAwsDbInstanceMetaData(resourceSchemaRepository)
	initAwsDbSubnetGroupMetaData(resourceSchemaRepository)
	initAwsDefaultSecurityGroupMetaData(resourceSchemaRepository)
	initAwsDefaultSubnetMetaData(resourceSchemaRepository)
	initAwsDefaultRouteTableMetadata(resourceSchemaRepository)
	initAwsDynamodbTableMetaData(resourceSchemaRepository)
	initAwsEbsSnapshotMetaData(resourceSchemaRepository)
	initAwsInstanceMetaData(resourceSchemaRepository)
	initAwsEbsVolumeMetaData(resourceSchemaRepository)
	initAwsEipMetaData(resourceSchemaRepository)
	initAwsS3BucketMetaData(resourceSchemaRepository)
	initAwsS3BucketPolicyMetaData(resourceSchemaRepository)
	initAwsS3BucketInventoryMetadata(resourceSchemaRepository)
	initAwsS3BucketMetricMetadata(resourceSchemaRepository)
	initAwsS3BucketNotificationMetadata(resourceSchemaRepository)
	initAwsS3BucketAnalyticsConfigurationMetaData(resourceSchemaRepository)
	initAwsEcrRepositoryMetaData(resourceSchemaRepository)
	initAwsRouteMetaData(resourceSchemaRepository)
	initAwsRouteTableAssociationMetaData(resourceSchemaRepository)
	initAwsRoute53RecordMetaData(resourceSchemaRepository)
	initAwsRoute53ZoneMetaData(resourceSchemaRepository)
	initAwsRoute53HealthCheckMetaData(resourceSchemaRepository)
	initSnsTopicSubscriptionMetaData(resourceSchemaRepository)
	initSnsTopicPolicyMetaData(resourceSchemaRepository)
	initSnsTopicMetaData(resourceSchemaRepository)
	initAwsIAMAccessKeyMetaData(resourceSchemaRepository)
	initAwsIAMPolicyMetaData(resourceSchemaRepository)
	initAwsIAMPolicyAttachmentMetaData(resourceSchemaRepository)
	initAwsIAMRoleMetaData(resourceSchemaRepository)
	initAwsIAMRolePolicyMetaData(resourceSchemaRepository)
	initAwsIamRolePolicyAttachmentMetaData(resourceSchemaRepository)
	initAwsIamUserPolicyAttachmentMetaData(resourceSchemaRepository)
	initAwsIAMUserMetaData(resourceSchemaRepository)
	initAwsIAMUserPolicyMetaData(resourceSchemaRepository)
	initAwsKeyPairMetaData(resourceSchemaRepository)
	initAwsKmsKeyMetaData(resourceSchemaRepository)
	initAwsKmsAliasMetaData(resourceSchemaRepository)
	initAwsLambdaFunctionMetaData(resourceSchemaRepository)
	initAwsLambdaEventSourceMappingMetaData(resourceSchemaRepository)
	initAwsSubnetMetaData(resourceSchemaRepository)
	initAwsSQSQueuePolicyMetaData(resourceSchemaRepository)
	initAwsSecurityGroupRuleMetaData(resourceSchemaRepository)
	initAwsSecurityGroupMetaData(resourceSchemaRepository)
	initAwsAppAutoscalingTargetMetaData(resourceSchemaRepository)
	initAwsAppAutoscalingPolicyMetaData(resourceSchemaRepository)
	initAwsAppAutoscalingScheduledActionMetaData(resourceSchemaRepository)
}
