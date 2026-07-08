package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/client"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  client.Config
		wantErr bool
	}{
		{
			name:    "valid config without extra headers",
			config:  client.Config{APIKey: "sk-1234", Temperature: 0.7},
			wantErr: false,
		},
		{
			name: "valid config with extra headers",
			config: client.Config{
				APIKey:      "sk-1234",
				Temperature: 0.7,
				ExtraHeaders: map[string]string{
					"X-App-Name": "YourAppName",
					"X-User-Id":  "user_123",
				},
			},
			wantErr: false,
		},
		{
			name: "empty header key",
			config: client.Config{
				APIKey:      "sk-1234",
				Temperature: 0.7,
				ExtraHeaders: map[string]string{
					"": "user_123",
				},
			},
			wantErr: true,
		},
		{
			name: "whitespace-only header key",
			config: client.Config{
				APIKey:      "sk-1234",
				Temperature: 0.7,
				ExtraHeaders: map[string]string{
					"   ": "user_123",
				},
			},
			wantErr: true,
		},
		{
			name: "empty header value",
			config: client.Config{
				APIKey:      "sk-1234",
				Temperature: 0.7,
				ExtraHeaders: map[string]string{
					"X-User-Id": "",
				},
			},
			wantErr: true,
		},
		{
			name: "whitespace-only header value",
			config: client.Config{
				APIKey:      "sk-1234",
				Temperature: 0.7,
				ExtraHeaders: map[string]string{
					"X-User-Id": "   ",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
