package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single value",
			input:    "http://localhost:3000",
			expected: []string{"http://localhost:3000"},
		},
		{
			name:     "multiple values",
			input:    "http://localhost:3000,http://localhost:5173",
			expected: []string{"http://localhost:3000", "http://localhost:5173"},
		},
		{
			name:     "values with spaces",
			input:    "  http://localhost:3000  ,  http://localhost:5173  ",
			expected: []string{"http://localhost:3000", "http://localhost:5173"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only commas",
			input:    ",,",
			expected: []string{},
		},
		{
			name:     "mixed empty and values",
			input:    "http://localhost:3000,,http://localhost:5173",
			expected: []string{"http://localhost:3000", "http://localhost:5173"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommaSeparated(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvOrViper(t *testing.T) {
	viper.Reset()

	tests := []struct {
		name       string
		envKey     string
		viperKey   string
		defaultVal string
		setupEnv   func()
		setupViper func()
		cleanup    func()
		expected   string
	}{
		{
			name:       "env takes precedence",
			envKey:     "TEST_ENV_KEY",
			viperKey:   "test.key",
			defaultVal: "default",
			setupEnv: func() {
				os.Setenv("TEST_ENV_KEY", "from_env")
			},
			setupViper: func() {
				viper.Set("test.key", "from_viper")
			},
			cleanup: func() {
				os.Unsetenv("TEST_ENV_KEY")
				viper.Reset()
			},
			expected: "from_env",
		},
		{
			name:       "viper when no env",
			envKey:     "TEST_ENV_KEY2",
			viperKey:   "test.key2",
			defaultVal: "default",
			setupEnv:   func() {},
			setupViper: func() {
				viper.Set("test.key2", "from_viper")
			},
			cleanup: func() {
				viper.Reset()
			},
			expected: "from_viper",
		},
		{
			name:       "default when nothing set",
			envKey:     "TEST_ENV_KEY3",
			viperKey:   "test.key3",
			defaultVal: "default_value",
			setupEnv:   func() {},
			setupViper: func() {},
			cleanup:    func() {},
			expected:   "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			tt.setupViper()
			defer tt.cleanup()

			result := getEnvOrViper(tt.envKey, tt.viperKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvOrViperInt(t *testing.T) {
	viper.Reset()

	tests := []struct {
		name       string
		envKey     string
		viperKey   string
		defaultVal int
		setupEnv   func()
		setupViper func()
		cleanup    func()
		expected   int
	}{
		{
			name:       "env int takes precedence",
			envKey:     "TEST_INT_ENV",
			viperKey:   "test.int",
			defaultVal: 10,
			setupEnv: func() {
				os.Setenv("TEST_INT_ENV", "42")
			},
			setupViper: func() {
				viper.Set("test.int", 99)
			},
			cleanup: func() {
				os.Unsetenv("TEST_INT_ENV")
				viper.Reset()
			},
			expected: 42,
		},
		{
			name:       "viper int when no env",
			envKey:     "TEST_INT_ENV2",
			viperKey:   "test.int2",
			defaultVal: 10,
			setupEnv:   func() {},
			setupViper: func() {
				viper.Set("test.int2", 55)
			},
			cleanup: func() {
				viper.Reset()
			},
			expected: 55,
		},
		{
			name:       "default int when nothing set",
			envKey:     "TEST_INT_ENV3",
			viperKey:   "test.int3",
			defaultVal: 100,
			setupEnv:   func() {},
			setupViper: func() {},
			cleanup:    func() {},
			expected:   100,
		},
		{
			name:       "invalid env int falls through to viper",
			envKey:     "TEST_INT_ENV4",
			viperKey:   "test.int4",
			defaultVal: 10,
			setupEnv: func() {
				os.Setenv("TEST_INT_ENV4", "not_a_number")
			},
			setupViper: func() {
				viper.Set("test.int4", 77)
			},
			cleanup: func() {
				os.Unsetenv("TEST_INT_ENV4")
				viper.Reset()
			},
			expected: 77,
		},
		{
			name:       "invalid env int falls through to default",
			envKey:     "TEST_INT_ENV5",
			viperKey:   "test.int5",
			defaultVal: 25,
			setupEnv: func() {
				os.Setenv("TEST_INT_ENV5", "invalid")
			},
			setupViper: func() {},
			cleanup: func() {
				os.Unsetenv("TEST_INT_ENV5")
			},
			expected: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			tt.setupViper()
			defer tt.cleanup()

			result := getEnvOrViperInt(tt.envKey, tt.viperKey, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}
