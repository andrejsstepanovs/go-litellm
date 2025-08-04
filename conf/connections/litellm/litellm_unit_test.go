package litellm_test

import (
	"net/url"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
)

func Test_Validate_Unit(t *testing.T) {
	validURL, _ := url.Parse("https://example.com")
	invalidURL := url.URL{} // zero value

	// URL with empty host
	emptyHostURL, _ := url.Parse("http://")
	emptyHostURL.Host = ""

	// URL with empty scheme
	emptySchemeURL, _ := url.Parse("example.com")
	emptySchemeURL.Scheme = ""

	// Invalid targets (default Target{} will fail validation)
	invalidTargets := litellm.Targets{
		System: litellm.Target{}, // This will fail validation
		LLM:    litellm.Target{},
		MCP:    litellm.Target{},
	}

	type testCase struct {
		name        string
		conn        *litellm.Connection
		wantErr     bool
		errContains string
	}

	tests := []testCase{
		{
			name: "not enough without env vars",
			conn: &litellm.Connection{
				URL:     *validURL,
				Targets: litellm.NewTargets(),
			},
			wantErr: true,
		},
		{
			name:        "nil connection pointer",
			conn:        nil,
			wantErr:     true,
			errContains: "litellm is required",
		},
		{
			name: "zero-value URL",
			conn: &litellm.Connection{
				URL:     invalidURL,
				Targets: litellm.NewTargets(),
			},
			wantErr:     true,
			errContains: "validation error",
		},
		{
			name: "URL with empty host",
			conn: &litellm.Connection{
				URL:     *emptyHostURL,
				Targets: litellm.NewTargets(),
			},
			wantErr: true,
		},
		{
			name: "URL with empty scheme",
			conn: &litellm.Connection{
				URL:     *emptySchemeURL,
				Targets: litellm.NewTargets(),
			},
			wantErr: true,
		},
		{
			name: "empty URL string",
			conn: &litellm.Connection{
				URL:     url.URL{}, // Zero value URL has empty String()
				Targets: litellm.NewTargets(),
			},
			wantErr: true,
		},
		{
			name: "targets validation failure",
			conn: &litellm.Connection{
				URL:     *validURL,
				Targets: invalidTargets,
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.conn.Validate()
			if tc.wantErr {
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

func Test_New_Unit(t *testing.T) {
	type testCase struct {
		name        string
		viperValue  interface{}
		wantErr     bool
		errContains string
		wantURL     string
	}

	tests := []testCase{
		{
			name:       "valid litellm.url",
			viperValue: "https://example.com",
			wantErr:    false,
			wantURL:    "https://example.com",
		},
		{
			name:        "invalid litellm.url (malformed)",
			viperValue:  "http://[::1]:namedport",
			wantErr:     true,
			errContains: "parse",
		},
		{
			name:       "empty string litellm.url",
			viperValue: "",
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Reset Viper before each test
			viper.Reset()
			if tc.viperValue != nil {
				viper.Set("litellm.url", tc.viperValue)
			}

			conn, err := litellm.New()
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantURL, conn.URL.String())
			}
		})
	}
}
