package repository

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	awstest "github.com/cloudskiff/driftctl/test/aws"

	"github.com/stretchr/testify/mock"

	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

func Test_cloudformationRepository_ListAllStacks(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeCloudformation)
		want    []*cloudformation.Stack
		wantErr error
	}{
		{
			name: "list multiple stacks",
			mocks: func(client *awstest.MockFakeCloudformation) {
				client.On("DescribeStacksPages",
					&cloudformation.DescribeStacksInput{},
					mock.MatchedBy(func(callback func(res *cloudformation.DescribeStacksOutput, lastPage bool) bool) bool {
						callback(&cloudformation.DescribeStacksOutput{
							Stacks: []*cloudformation.Stack{
								{StackId: aws.String("stack1")},
								{StackId: aws.String("stack2")},
								{StackId: aws.String("stack3")},
							},
						}, false)
						callback(&cloudformation.DescribeStacksOutput{
							Stacks: []*cloudformation.Stack{
								{StackId: aws.String("stack4")},
								{StackId: aws.String("stack5")},
								{StackId: aws.String("stack6")},
							},
						}, true)
						return true
					})).Return(nil).Once()
			},
			want: []*cloudformation.Stack{
				{StackId: aws.String("stack1")},
				{StackId: aws.String("stack2")},
				{StackId: aws.String("stack3")},
				{StackId: aws.String("stack4")},
				{StackId: aws.String("stack5")},
				{StackId: aws.String("stack6")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := cache.New(1)
			client := awstest.MockFakeCloudformation{}
			tt.mocks(&client)
			r := &cloudformationRepository{
				client: &client,
				cache:  store,
			}
			got, err := r.ListAllStacks()
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				// Check that results were cached
				cachedData, err := r.ListAllStacks()
				assert.NoError(t, err)
				assert.Equal(t, got, cachedData)
				assert.IsType(t, []*cloudformation.Stack{}, store.Get("cloudformationListAllStacks"))
			}

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %s -> %s", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
		})
	}
}
