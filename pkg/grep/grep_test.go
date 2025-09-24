package grep

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	cfg "grep_utility/pkg/config"
)

// Helper function to create config with default values
func createConfig() *cfg.GrepConfig {
	return &cfg.GrepConfig{
		Pattern:    "",
		FilePath:   "",
		A:          0,
		B:          0,
		C:          0,
		Count:      false,
		IgnoreCase: false,
		Invert:     false,
		Fixed:      false,
		LineNumber: false,
	}
}

// Helper function to capture stdout output
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	fn()
	w.Close()
	os.Stdout = old
	out := <-outC
	return out
}

func TestRunGrepBasicMatch(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"

	input := "hello\nno match\nother line"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> no match <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepCountFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"
	config.Count = true

	input := "hello\nno match\nanother no\nother line"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "2\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepIgnoreCaseFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "NO"
	config.IgnoreCase = true

	input := "hello\nno match\nNO MATCH"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> no match <-\n-> NO MATCH <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepInvertFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"
	config.Invert = true

	input := "hello\nno match\nother line"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> hello <-\n-> other line <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepLineNumberFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"
	config.LineNumber = true

	input := "hello\nno match\nother line\nno again"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> 2 no match <-\n-> 4 no again <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepAfterContextFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.A = 2

	input := "line1\nThis is a match\nafter1\nafter2\nafter3\nline6"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> This is a match <-\nafter1\nafter2\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepBeforeContextFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.B = 2

	input := "before1\nbefore2\nThis is a match\nafter1\nline5"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "before1\nbefore2\n-> This is a match <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepContextFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.C = 1

	input := "line1\nbefore\nThis is a match\nafter\nline5"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "before\n-> This is a match <-\nafter\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepContextWithLineNumbers(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.C = 1
	config.LineNumber = true

	input := "line1\nbefore\nThis is a match\nafter\nline5"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "2 before\n-> 3 This is a match <-\n4 after\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepMultipleMatches(t *testing.T) {
	config := createConfig()
	config.Pattern = "test"
	config.A = 1

	input := "line1\ntest1\nafter1\ntest2\nafter2\nline6"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> test1 <-\nafter1\n-> test2 <-\nafter2\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepInvalidRegex(t *testing.T) {
	config := createConfig()
	config.Pattern = "[invalid"

	input := "test line"
	reader := strings.NewReader(input)

	err := RunGrep(config, reader)
	if err == nil {
		t.Error("Expected error for invalid regex, got nil")
	}
}

func TestRunGrepEmptyInput(t *testing.T) {
	config := createConfig()
	config.Pattern = "test"

	input := ""
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	if output != "" {
		t.Errorf("Expected empty output, got '%s'", output)
	}
}

func TestRunGrepNoMatches(t *testing.T) {
	config := createConfig()
	config.Pattern = "nomatch"
	config.Count = true

	input := "line1\nline2\nline3"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "0\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepFixedString(t *testing.T) {
	config := createConfig()
	config.Pattern = "n.*o"
	config.Fixed = true

	input := "hello\nn.*o literal\nno match"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> n.*o literal <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepFixedWithIgnoreCase(t *testing.T) {
	config := createConfig()
	config.Pattern = "Test"
	config.Fixed = true
	config.IgnoreCase = true

	input := "test\nTEST\nTest\nother"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> test <-\n-> TEST <-\n-> Test <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestRunGrepRegexPattern(t *testing.T) {
	config := createConfig()
	config.Pattern = "test[0-9]+"

	input := "test\ntest123\ntest456\ntestABC"
	reader := strings.NewReader(input)

	output := captureStdout(func() {
		err := RunGrep(config, reader)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}
	})

	expected := "-> test123 <-\n-> test456 <-\n"
	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}

func TestIsMatchFunction(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		config   *cfg.GrepConfig
		expected bool
	}{
		{
			name: "Basic match",
			line: "hello world",
			config: &cfg.GrepConfig{
				Pattern: "world",
				Fixed:   true,
			},
			expected: true,
		},
		{
			name: "Case insensitive match",
			line: "Hello World",
			config: &cfg.GrepConfig{
				Pattern:    "world",
				Fixed:      true,
				IgnoreCase: true,
			},
			expected: true,
		},
		{
			name: "Inverted match",
			line: "hello world",
			config: &cfg.GrepConfig{
				Pattern: "xyz",
				Fixed:   true,
				Invert:  true,
			},
			expected: true,
		},
		{
			name: "Regex match",
			line: "test123",
			config: &cfg.GrepConfig{
				Pattern: "test[0-9]+",
				Fixed:   false,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var regex *regexp.Regexp
			var err error

			if !tt.config.Fixed {
				pattern := tt.config.Pattern
				if tt.config.IgnoreCase {
					pattern = "(?i)" + pattern
				}
				regex, err = regexp.Compile(pattern)
				if err != nil {
					t.Fatalf("Failed to compile regex: %v", err)
				}
			}

			result := isMatch(tt.line, tt.config, regex)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDataConstructor(t *testing.T) {
	config := cfg.GrepConfig{
		B: 5,
	}
	data := newData(config)

	if data.Config.B != 5 {
		t.Errorf("Expected B=5, got %d", data.Config.B)
	}
	if data.matchCount != 0 {
		t.Errorf("Expected matchCount=0, got %d", data.matchCount)
	}
	if data.lineNumber != 0 {
		t.Errorf("Expected lineNumber=0, got %d", data.lineNumber)
	}
	if data.lastPrinted != -1 {
		t.Errorf("Expected lastPrinted=-1, got %d", data.lastPrinted)
	}
}
