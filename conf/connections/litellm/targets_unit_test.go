package litellm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/andrejsstepanovs/go-litellm/conf/connections/litellm"
)

func Test_Get_Unit(t *testing.T) {
	// Create distinct Target values for each field
	systemTarget := litellm.Target{Timeout: 1 * time.Second}
	llmTarget := litellm.Target{Timeout: 2 * time.Second}
	mcpTarget := litellm.Target{Timeout: 3 * time.Second}

	targets := litellm.Targets{
		System: systemTarget,
		LLM:    llmTarget,
		MCP:    mcpTarget,
	}

	type testCase struct {
		name     string
		input    litellm.TargetName
		expected litellm.Target
	}

	testCases := []testCase{
		{
			name:     "Get System Target",
			input:    litellm.CLIENT_SYSTEM,
			expected: systemTarget,
		},
		{
			name:     "Get LLM Target",
			input:    litellm.CLIENT_LLM,
			expected: llmTarget,
		},
		{
			name:     "Get MCP Target",
			input:    litellm.CLIENT_MCP,
			expected: mcpTarget,
		},
		{
			name:     "Unknown TargetName returns zero Target",
			input:    litellm.TargetName("unknown"),
			expected: litellm.Target{},
		},
		{
			name:     "Empty TargetName returns zero Target",
			input:    litellm.TargetName(""),
			expected: litellm.Target{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := targets.Get(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_NewTargets_Unit(t *testing.T) {
	type testCase struct {
		name           string
		viperSetup     func()
		expectedSystem time.Duration
		expectedLLM    time.Duration
		expectedMCP    time.Duration
	}

	validSystem := 5 * time.Second
	validLLM := 10 * time.Second
	validMCP := 15 * time.Second

	testCases := []testCase{
		{
			name: "all present and valid",
			viperSetup: func() {
				viper.Reset()
				viper.Set("litellm.targets.system.timeout", validSystem.String())
				viper.Set("litellm.targets.llm.timeout", validLLM.String())
				viper.Set("litellm.targets.mcp.timeout", validMCP.String())
			},
			expectedSystem: validSystem,
			expectedLLM:    validLLM,
			expectedMCP:    validMCP,
		},
		{
			name: "system missing",
			viperSetup: func() {
				viper.Reset()
				// viper.Set("litellm.targets.system.timeout", ...) // missing
				viper.Set("litellm.targets.llm.timeout", validLLM.String())
				viper.Set("litellm.targets.mcp.timeout", validMCP.String())
			},
			expectedSystem: 0,
			expectedLLM:    validLLM,
			expectedMCP:    validMCP,
		},
		{
			name: "llm missing",
			viperSetup: func() {
				viper.Reset()
				viper.Set("litellm.targets.system.timeout", validSystem.String())
				// viper.Set("litellm.targets.llm.timeout", ...) // missing
				viper.Set("litellm.targets.mcp.timeout", validMCP.String())
			},
			expectedSystem: validSystem,
			expectedLLM:    0,
			expectedMCP:    validMCP,
		},
		{
			name: "mcp missing",
			viperSetup: func() {
				viper.Reset()
				viper.Set("litellm.targets.system.timeout", validSystem.String())
				viper.Set("litellm.targets.llm.timeout", validLLM.String())
				// viper.Set("litellm.targets.mcp.timeout", ...) // missing
			},
			expectedSystem: validSystem,
			expectedLLM:    validLLM,
			expectedMCP:    0,
		},
		{
			name: "all missing",
			viperSetup: func() {
				viper.Reset()
				// none set
			},
			expectedSystem: 0,
			expectedLLM:    0,
			expectedMCP:    0,
		},
		{
			name: "invalid values",
			viperSetup: func() {
				viper.Reset()
				viper.Set("litellm.targets.system.timeout", "notaduration")
				// viper.Set("litellm.targets.llm.timeout", 12345) // not a string
				// viper.Set("litellm.targets.mcp.timeout", "")    // empty string
			},
			expectedSystem: 0,
			expectedLLM:    0,
			expectedMCP:    0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.viperSetup()
			targets := litellm.NewTargets()
			assert.Equal(t, tc.expectedSystem, targets.System.Timeout, "system timeout")
			assert.Equal(t, tc.expectedLLM, targets.LLM.Timeout, "llm timeout")
			assert.Equal(t, tc.expectedMCP, targets.MCP.Timeout, "mcp timeout")
		})
	}
}

func Test_Validate_Targets_Unit(t *testing.T) {
	// Helper: valid and invalid Target
	validTarget := litellm.Target{
		Timeout:          1 * time.Second,
		RetryInterval:    1 * time.Second,
		RetryMaxAttempts: 1,
		RetryBackoffRate: 1.0,
		MaxRetry:         1,
	}
	invalidTarget := litellm.Target{} // Assumes zero values are invalid per Target.Validate()

	type testCase struct {
		name        string
		targets     litellm.Targets
		expectError bool
		expectMsgs  []string // substrings expected in error message
	}

	testCases := []testCase{
		{
			name: "all valid",
			targets: litellm.Targets{
				System: validTarget,
				LLM:    validTarget,
				MCP:    validTarget,
			},
			expectError: false,
		},
		{
			name: "system invalid",
			targets: litellm.Targets{
				System: invalidTarget,
				LLM:    validTarget,
				MCP:    validTarget,
			},
			expectError: true,
			expectMsgs:  []string{"system target"},
		},
		{
			name: "llm invalid",
			targets: litellm.Targets{
				System: validTarget,
				LLM:    invalidTarget,
				MCP:    validTarget,
			},
			expectError: true,
			expectMsgs:  []string{"llm target"},
		},
		{
			name: "mcp invalid",
			targets: litellm.Targets{
				System: validTarget,
				LLM:    validTarget,
				MCP:    invalidTarget,
			},
			expectError: true,
			expectMsgs:  []string{"mcp target"},
		},
		{
			name: "system and llm invalid",
			targets: litellm.Targets{
				System: invalidTarget,
				LLM:    invalidTarget,
				MCP:    validTarget,
			},
			expectError: true,
			expectMsgs:  []string{"system target", "llm target"},
		},
		{
			name: "all invalid",
			targets: litellm.Targets{
				System: invalidTarget,
				LLM:    invalidTarget,
				MCP:    invalidTarget,
			},
			expectError: true,
			expectMsgs:  []string{"system target", "llm target", "mcp target"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.targets.Validate()
			if tc.expectError {
				assert.Error(t, err)
				for _, msg := range tc.expectMsgs {
					assert.ErrorContains(t, err, msg, fmt.Sprintf("expected error to contain %q", msg))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
