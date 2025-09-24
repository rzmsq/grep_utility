package config

import (
	"flag"
	"os"
	"testing"
)

// TestNewFromFlags tests the flag parsing functionality
func TestNewFromFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *GrepConfig
	}{
		{
			name: "Default flags",
			args: []string{"grep", "pattern"},
			expected: &GrepConfig{
				Pattern:    "pattern",
				FilePath:   "",
				A:          0,
				B:          0,
				C:          0,
				Count:      false,
				IgnoreCase: false,
				Invert:     false,
				Fixed:      false,
				LineNumber: false,
			},
		},
		{
			name: "All flags set",
			args: []string{"grep", "-A", "2", "-B", "1", "-C", "3", "-c", "-i", "-v", "-F", "-n", "-f", "test.txt", "pattern"},
			expected: &GrepConfig{
				Pattern:    "pattern",
				FilePath:   "test.txt",
				A:          2,
				B:          1,
				C:          3,
				Count:      true,
				IgnoreCase: true,
				Invert:     true,
				Fixed:      true,
				LineNumber: true,
			},
		},
		{
			name: "No pattern provided",
			args: []string{"grep", "-c"},
			expected: &GrepConfig{
				Pattern:    "",
				FilePath:   "",
				A:          0,
				B:          0,
				C:          0,
				Count:      true,
				IgnoreCase: false,
				Invert:     false,
				Fixed:      false,
				LineNumber: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag package
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set os.Args
			oldArgs := os.Args
			os.Args = tt.args
			defer func() { os.Args = oldArgs }()

			config := NewFromFlags()

			if config.Pattern != tt.expected.Pattern {
				t.Errorf("Pattern: expected '%s', got '%s'", tt.expected.Pattern, config.Pattern)
			}
			if config.FilePath != tt.expected.FilePath {
				t.Errorf("FilePath: expected '%s', got '%s'", tt.expected.FilePath, config.FilePath)
			}
			if config.A != tt.expected.A {
				t.Errorf("A: expected %d, got %d", tt.expected.A, config.A)
			}
			if config.B != tt.expected.B {
				t.Errorf("B: expected %d, got %d", tt.expected.B, config.B)
			}
			if config.C != tt.expected.C {
				t.Errorf("C: expected %d, got %d", tt.expected.C, config.C)
			}
			if config.Count != tt.expected.Count {
				t.Errorf("Count: expected %t, got %t", tt.expected.Count, config.Count)
			}
			if config.IgnoreCase != tt.expected.IgnoreCase {
				t.Errorf("IgnoreCase: expected %t, got %t", tt.expected.IgnoreCase, config.IgnoreCase)
			}
			if config.Invert != tt.expected.Invert {
				t.Errorf("Invert: expected %t, got %t", tt.expected.Invert, config.Invert)
			}
			if config.Fixed != tt.expected.Fixed {
				t.Errorf("Fixed: expected %t, got %t", tt.expected.Fixed, config.Fixed)
			}
			if config.LineNumber != tt.expected.LineNumber {
				t.Errorf("LineNumber: expected %t, got %t", tt.expected.LineNumber, config.LineNumber)
			}
		})
	}
}

// TestGrepConfigStruct tests the GrepConfig struct
func TestGrepConfigStruct(t *testing.T) {
	config := &GrepConfig{
		Pattern:    "test",
		FilePath:   "file.txt",
		A:          1,
		B:          2,
		C:          3,
		Count:      true,
		IgnoreCase: true,
		Invert:     false,
		Fixed:      true,
		LineNumber: true,
	}

	if config.Pattern != "test" {
		t.Errorf("Expected pattern 'test', got '%s'", config.Pattern)
	}
	if config.A != 1 {
		t.Errorf("Expected A = 1, got %d", config.A)
	}
	if !config.Count {
		t.Error("Expected Count = true")
	}
}
