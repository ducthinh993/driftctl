package resource

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type Resource interface {
	TerraformId() string
	TerraformType() string
	CtyValue() *cty.Value
}

var refactoredResources = []string{
	"aws_ami",
	"aws_cloudfront_distribution",
	"aws_db_instance",
	"aws_db_subnet_group",
	"aws_default_route_table",
	"aws_default_security_group",
	"aws_default_subnet",
	"aws_default_vpc",
	"aws_dynamodb_table",
	"aws_ebs_snapshot",
	"aws_ebs_volume",
	"aws_ecr_repository",
	"aws_eip",
	"aws_eip_association",
	"aws_iam_access_key",
	"aws_iam_policy",
	"aws_iam_policy_attachment",
	"aws_iam_role",
	"aws_iam_role_policy",
	"aws_iam_role_policy_attachment",
	"aws_iam_user",
	"aws_iam_user_policy",
	"aws_iam_user_policy_attachment",
	"aws_instance",
	"aws_internet_gateway",
	"aws_key_pair",
	"aws_kms_alias",
	"aws_kms_key",
	// "aws_lambda_event_source_mapping",
	// "aws_lambda_function",
	"aws_nat_gateway",
	"aws_route",
	"aws_route53_health_check",
	"aws_route53_record",
	"aws_route53_zone",
	"aws_route_table",
	"aws_route_table_association",
	"aws_s3_bucket",
	"aws_s3_bucket_analytics_configuration",
	"aws_s3_bucket_inventory",
	"aws_s3_bucket_metric",
	"aws_s3_bucket_notification",
	"aws_s3_bucket_policy",
	// "aws_security_group",
	// "aws_security_group_rule",
	"aws_sns_topic",
	"aws_sns_topic_policy",
	"aws_sns_topic_subscription",
	// "aws_sqs_queue",
	// "aws_sqs_queue_policy",
	// "aws_subnet",
	// "aws_vpc",

	"github_branch_protection",
	"github_membership",
	"github_repository",
	"github_team",
	"github_team_membership",
}

func IsRefactoredResource(typ string) bool {
	for _, refactoredResource := range refactoredResources {
		if typ == refactoredResource {
			return true
		}
	}
	return false
}

type AbstractResource struct {
	Id    string
	Type  string
	Attrs *Attributes
}

func (a *AbstractResource) TerraformId() string {
	return a.Id
}

func (a *AbstractResource) TerraformType() string {
	return a.Type
}

func (a *AbstractResource) CtyValue() *cty.Value {
	return nil
}

type ResourceFactory interface {
	CreateResource(data interface{}, ty string) (*cty.Value, error)
	CreateAbstractResource(ty, id string, data map[string]interface{}) *AbstractResource
}

type SerializableResource struct {
	Resource
}

type SerializedResource struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

func (u SerializedResource) TerraformId() string {
	return u.Id
}

func (u SerializedResource) TerraformType() string {
	return u.Type
}

func (u SerializedResource) CtyValue() *cty.Value {
	return &cty.NilVal
}

func (s *SerializableResource) UnmarshalJSON(bytes []byte) error {
	var res SerializedResource

	if err := json.Unmarshal(bytes, &res); err != nil {
		return err
	}
	s.Resource = res
	return nil
}

func (s SerializableResource) MarshalJSON() ([]byte, error) {
	return json.Marshal(SerializedResource{Id: s.TerraformId(), Type: s.TerraformType()})
}

type NormalizedResource interface {
	NormalizeForState() (Resource, error)
	NormalizeForProvider() (Resource, error)
}

func IsSameResource(rRs, lRs Resource) bool {
	return rRs.TerraformType() == lRs.TerraformType() && rRs.TerraformId() == lRs.TerraformId()
}

func Sort(res []Resource) []Resource {
	sort.SliceStable(res, func(i, j int) bool {
		if res[i].TerraformType() != res[j].TerraformType() {
			return res[i].TerraformType() < res[j].TerraformType()
		}
		return res[i].TerraformId() < res[j].TerraformId()
	})
	return res
}

func ToResourceAttributes(val *cty.Value) *Attributes {
	if val == nil {
		return nil
	}

	bytes, _ := ctyjson.Marshal(*val, val.Type())
	var attrs Attributes
	err := json.Unmarshal(bytes, &attrs)
	if err != nil {
		panic(err)
	}

	return &attrs
}

type Attributes map[string]interface{}

func (a *Attributes) Get(path string) (interface{}, bool) {
	val, exist := (*a)[path]
	return val, exist
}

func (a *Attributes) SafeDelete(path []string) {
	for i, key := range path {
		if i == len(path)-1 {
			delete(*a, key)
			return
		}

		v, exists := (*a)[key]
		if !exists {
			return
		}
		m, ok := v.(Attributes)
		if !ok {
			return
		}
		*a = m
	}
}

func (a *Attributes) SafeSet(path []string, value interface{}) error {
	for i, key := range path {
		if i == len(path)-1 {
			(*a)[key] = value
			return nil
		}

		v, exists := (*a)[key]
		if !exists {
			(*a)[key] = map[string]interface{}{}
			v = (*a)[key]
		}

		m, ok := v.(Attributes)
		if !ok {
			return errors.Errorf("Path %s cannot be set: %s is not a nested struct", strings.Join(path, "."), key)
		}
		*a = m
	}
	return errors.New("Error setting value") // should not happen ?
}

func (a *Attributes) DeleteIfDefault(path string) {
	val, exist := a.Get(path)
	ty := reflect.TypeOf(val)
	if exist && val == reflect.Zero(ty).Interface() {
		a.SafeDelete([]string{path})
	}
}

func concatenatePath(path, next string) string {
	if path == "" {
		return next
	}
	return strings.Join([]string{path, next}, ".")
}

func (a *Attributes) SanitizeDefaults() {
	original := reflect.ValueOf(*a)
	copy := reflect.New(original.Type()).Elem()
	a.sanitize("", original, copy)
	*a = copy.Interface().(Attributes)
}

func (a *Attributes) sanitize(path string, original, copy reflect.Value) bool {
	switch original.Kind() {
	case reflect.Ptr:
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return false
		}
		copy.Set(reflect.New(originalValue.Type()))
		a.sanitize(path, originalValue, copy.Elem())
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return false
		}
		if originalValue.Kind() == reflect.Slice || originalValue.Kind() == reflect.Map {
			if originalValue.Len() == 0 {
				return false
			}
		}
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		a.sanitize(path, originalValue, copyValue)
		copy.Set(copyValue)

	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			field := original.Field(i)
			a.sanitize(concatenatePath(path, field.String()), field, copy.Field(i))
		}
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			a.sanitize(concatenatePath(path, strconv.Itoa(i)), original.Index(i), copy.Index(i))
		}
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			created := a.sanitize(concatenatePath(path, key.String()), originalValue, copyValue)
			if created {
				copy.SetMapIndex(key, copyValue)
			}
		}
	default:
		copy.Set(original)
	}
	return true
}
