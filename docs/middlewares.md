# Middlewares

Main goal of middleware is to make IaC and remote resources set to match. For this we filter some default resources, mutate or remove some fields in resources or even create and delete new resources.

```go
type AwsDefaultRoute struct{}

func NewAwsDefaultRoute() AwsDefaultRoute {
	return AwsDefaultRoute{}
}

func (m AwsDefaultRoute) Execute(remoteResources, resourcesFromState *[]*resource.Resource) error {

	newRemoteResources := make([]*resource.Resource, 0)

	// ...

	*remoteResources = newRemoteResources
  
  return nil
}
```

In the above example, we define a middleware called `AwsDefaultRoute` that will modify remote resources. Middleware can access two `*[]*resource.Resource` : the one from IaC side and the other one from remote side and the goal is to rework theses slice to remove false positive drifts. Notice both `remoteResources` and `resourcesFromState` variables are pointers, which mean middlewares can perform mutations on resources before the comparison is made.

## Different kind of middlewares

1) Help driftctl match IaC resources and remote resources.
2) Filter noise from provider default resources
3) Resource transformation
4) Specific edge cases

## Examples

1) `aws_route_table_expander` explode inline route in dedicated resources
2) `aws_default_route` ignore route created by default when creating a table **if they are not managed in IaC**
3) `aws_iam_user_policy_attachment` and `aws_iam_role_policy_attachment` transformed to `aws_iam_policy_attachment`
4) `route53_records_id_reconcilier` that rework ID's to match terraform ones
