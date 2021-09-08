package terraform

import (
	"testing"

	"github.com/hashicorp/terraform/addrs"
	"github.com/stretchr/testify/assert"
)

func Test_ReadLockFile(t *testing.T) {
	cases := []struct {
		test     string
		filepath string
		assert   func(*Locks, error)
	}{
		{
			test:     "should read valid lock file",
			filepath: "testdata/lockfile_valid.hcl",
			assert: func(locks *Locks, err error) {
				provider := locks.Provider(&addrs.Provider{
					Type:      "aws",
					Namespace: "hashicorp",
					Hostname:  "registry.terraform.io",
				})

				assert.Len(t, locks.providers, 10)
				assert.Equal(t, "3.47.0", provider.Version)
				assert.Nil(t, err)
			},
		},
		{
			test:     "should fail to read with invalid address",
			filepath: "testdata/lockfile_valid.hcl",
			assert: func(locks *Locks, err error) {
				provider := locks.Provider(&addrs.Provider{})

				assert.Len(t, locks.providers, 10)
				assert.Nil(t, provider)
				assert.Nil(t, err)
			},
		},
		{
			test:     "should read empty file without error",
			filepath: "testdata/lockfile_empty.hcl",
			assert: func(locks *Locks, err error) {
				provider := locks.Provider(&addrs.Provider{})

				assert.Len(t, locks.providers, 0)
				assert.Nil(t, provider)
				assert.Nil(t, err)
			},
		},
		{
			test:     "should return error for invalid lock file",
			filepath: "testdata/lockfile_invalid.hcl",
			assert: func(locks *Locks, err error) {
				provider := locks.Provider(&addrs.Provider{})

				assert.Len(t, locks.providers, 1)
				assert.Nil(t, provider)
				assert.EqualError(t, err, "Missing required argument: The argument \"version\" is required, but no definition was found.")
			},
		},
		{
			test:     "should return error for invalid provider block",
			filepath: "testdata/lockfile_invalid_provider.hcl",
			assert: func(locks *Locks, err error) {
				provider := locks.Provider(&addrs.Provider{})

				assert.Len(t, locks.providers, 1)
				assert.Nil(t, provider)
				assert.EqualError(t, err, "the provider source address for this provider lock must be written as \"registry.terraform.io/hashicorp/invalid\", the fully-qualified and normalized form")
			},
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(tt *testing.T) {
			locks, diags := LoadLocksFromFile(c.filepath)
			c.assert(locks, diags)
		})
	}
}
