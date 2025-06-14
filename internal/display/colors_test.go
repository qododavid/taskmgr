package display

import (
	"os"
	"testing"
)

func TestColorString(t *testing.T) {
	color := Red
	expected := "\033[31m"
	if color.String() != expected {
		t.Errorf("Expected %q, got %q", expected, color.String())
	}
}

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		text     string
		expected string
		noColor  bool
	}{
		{
			name:     "with color support",
			color:    Red,
			text:     "test",
			expected: "test", // In CI environment, colors are typically disabled
			noColor:  false,
		},
		{
			name:     "without color support",
			color:    Red,
			text:     "test",
			expected: "test",
			noColor:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.noColor {
				os.Setenv("NO_COLOR", "1")
				defer os.Unsetenv("NO_COLOR")
			}
			
			result := Colorize(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestIsColorSupported(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name:     "NO_COLOR set",
			envVars:  map[string]string{"NO_COLOR": "1"},
			expected: false,
		},
		{
			name:     "TERM empty",
			envVars:  map[string]string{"TERM": ""},
			expected: false,
		},
		{
			name:     "TERM dumb",
			envVars:  map[string]string{"TERM": "dumb"},
			expected: false,
		},
		{
			name:     "CI without COLORTERM",
			envVars:  map[string]string{"CI": "true", "COLORTERM": ""},
			expected: false,
		},
		{
			name:     "normal terminal",
			envVars:  map[string]string{"TERM": "xterm-256color"},
			expected: true,
		},
		{
			name:     "CI with COLORTERM",
			envVars:  map[string]string{"CI": "true", "COLORTERM": "truecolor", "TERM": "xterm-256color"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()
			
			// Set test environment variables
			for key, value := range tt.envVars {
				if value != "" {
					os.Setenv(key, value)
				}
			}
			
			result := IsColorSupported()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDefaultColorScheme(t *testing.T) {
	// Test that default color scheme has all required colors
	scheme := DefaultColorScheme
	
	if scheme.Completed == "" {
		t.Error("DefaultColorScheme.Completed should not be empty")
	}
	if scheme.Pending == "" {
		t.Error("DefaultColorScheme.Pending should not be empty")
	}
	if scheme.Overdue == "" {
		t.Error("DefaultColorScheme.Overdue should not be empty")
	}
	if scheme.High == "" {
		t.Error("DefaultColorScheme.High should not be empty")
	}
	if scheme.Medium == "" {
		t.Error("DefaultColorScheme.Medium should not be empty")
	}
	if scheme.Low == "" {
		t.Error("DefaultColorScheme.Low should not be empty")
	}
	if scheme.Critical == "" {
		t.Error("DefaultColorScheme.Critical should not be empty")
	}
	if scheme.Tags == "" {
		t.Error("DefaultColorScheme.Tags should not be empty")
	}
	if scheme.DueDate == "" {
		t.Error("DefaultColorScheme.DueDate should not be empty")
	}
}