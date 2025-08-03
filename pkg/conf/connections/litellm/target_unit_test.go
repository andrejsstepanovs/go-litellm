package litellm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/pkg/conf/connections/litellm"
)

func Test_Target_Validate_Unit(t *testing.T) {
	validTarget := &litellm.Target{
		Timeout:          5 * time.Second,
		RetryInterval:    1 * time.Second,
		RetryMaxAttempts: 3,
		RetryBackoffRate: 2.0,
		MaxRetry:         5,
	}

	tests := []struct {
		name        string
		target      *litellm.Target
		expectError bool
		errContains string
	}{
		{
			name:        "happy path - all fields valid",
			target:      validTarget,
			expectError: false,
		},
		{
			name:        "nil receiver",
			target:      nil,
			expectError: true,
			errContains: "target is required",
		},
		{
			name: "zero Timeout",
			target: &litellm.Target{
				Timeout:          0,
				RetryInterval:    1 * time.Second,
				RetryMaxAttempts: 3,
				RetryBackoffRate: 2.0,
				MaxRetry:         5,
			},
			expectError: true,
			errContains: "Timeout",
		},
		{
			name: "zero RetryInterval",
			target: &litellm.Target{
				Timeout:          5 * time.Second,
				RetryInterval:    0,
				RetryMaxAttempts: 3,
				RetryBackoffRate: 2.0,
				MaxRetry:         5,
			},
			expectError: true,
			errContains: "RetryInterval",
		},
		{
			name: "zero RetryMaxAttempts",
			target: &litellm.Target{
				Timeout:          5 * time.Second,
				RetryInterval:    1 * time.Second,
				RetryMaxAttempts: 0,
				RetryBackoffRate: 2.0,
				MaxRetry:         5,
			},
			expectError: true,
			errContains: "RetryMaxAttempts",
		},
		{
			name: "zero RetryBackoffRate",
			target: &litellm.Target{
				Timeout:          5 * time.Second,
				RetryInterval:    1 * time.Second,
				RetryMaxAttempts: 3,
				RetryBackoffRate: 0,
				MaxRetry:         5,
			},
			expectError: true,
			errContains: "RetryBackoffRate",
		},
		{
			name: "zero MaxRetry",
			target: &litellm.Target{
				Timeout:          5 * time.Second,
				RetryInterval:    1 * time.Second,
				RetryMaxAttempts: 3,
				RetryBackoffRate: 2.0,
				MaxRetry:         0,
			},
			expectError: true,
			errContains: "MaxRetry",
		},
		{
			name: "all zero values",
			target: &litellm.Target{
				Timeout:          0,
				RetryInterval:    0,
				RetryMaxAttempts: 0,
				RetryBackoffRate: 0,
				MaxRetry:         0,
			},
			expectError: true,
			errContains: "validation error",
		},
		{
			name: "boundary values - minimal valid",
			target: &litellm.Target{
				Timeout:          1,
				RetryInterval:    1,
				RetryMaxAttempts: 1,
				RetryBackoffRate: 0.0001,
				MaxRetry:         1,
			},
			expectError: false,
		},
		{
			name: "invalid combination - RetryInterval zero but retries set",
			target: &litellm.Target{
				Timeout:          5 * time.Second,
				RetryInterval:    0,
				RetryMaxAttempts: 5,
				RetryBackoffRate: 2.0,
				MaxRetry:         5,
			},
			expectError: true,
			errContains: "RetryInterval",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := func() error {
				if tc.target == nil {
					return (*litellm.Target)(nil).Validate()
				}
				return tc.target.Validate()
			}()
			if tc.expectError {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
