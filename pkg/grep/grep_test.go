package grep

import (
	"bufio"
	"fmt"
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

// Helper function to capture output
func captureOutput(config *cfg.GrepConfig, input string) (string, error) {
	reader := strings.NewReader(input)

	// Temporary redirect stdout to capture output
	// Create a custom version that writes to buffer instead of stdout
	data := newData(*config)
	scanner := bufio.NewScanner(reader)

	pattern := data.Config.Pattern
	if data.Config.IgnoreCase && !data.Config.Fixed {
		pattern = "(?i)" + pattern
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	var output []string

	for scanner.Scan() {
		data.lineNumber++
		line := scanner.Text()

		if isMatch(line, config, regex) {
			data.matchCount++

			if data.Config.Count {
				continue
			}

			// Handle before context
			if data.Config.B > 0 || data.Config.C > 0 {
				i := max(data.Config.B, data.Config.C)
				if i >= len(data.beforeBuffer) {
					i = 0
				}
				for i < len(data.beforeBuffer) {
					if data.Config.LineNumber {
						lineNum := data.lineNumber - len(data.beforeBuffer) + i
						output = append(output, fmt.Sprintf("%d %s", lineNum, data.beforeBuffer[i]))
					} else {
						output = append(output, data.beforeBuffer[i])
					}
					i++
				}
				clear(data.beforeBuffer)
			}

			// Print matched line
			if data.Config.LineNumber {
				output = append(output, fmt.Sprintf("-> %d %s <-", data.lineNumber, line))
			} else {
				output = append(output, "-> "+line+" <-")
			}

			data.lastPrinted = data.lineNumber
			data.linesToPrintAfterMatch = max(data.Config.A, data.Config.C)
			continue
		}

		// Handle after context
		if data.linesToPrintAfterMatch > 0 {
			if data.lineNumber > data.lastPrinted {
				if data.Config.LineNumber {
					output = append(output, fmt.Sprintf("%d %s", data.lineNumber, line))
				} else {
					output = append(output, line)
				}
				data.lastPrinted = data.lineNumber
				data.linesToPrintAfterMatch--
			} else {
				data.linesToPrintAfterMatch--
			}
		}

		// Update before buffer
		if data.Config.B > 0 || data.Config.C > 0 {
			data.beforeBuffer = append(data.beforeBuffer, line)
			maxBuffer := max(data.Config.B, data.Config.C)
			if len(data.beforeBuffer) > maxBuffer {
				data.beforeBuffer = data.beforeBuffer[1:]
			}
		}
	}

	if data.Config.Count {
		return fmt.Sprintf("%d", data.matchCount), nil
	}

	return strings.Join(output, "\n"), nil
}

func TestBasicMatch(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"

	input := "hello\nno match\nother line"
	expected := "-> no match <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCountFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"
	config.Count = true

	input := "hello\nno match\nanother no\nother line"
	expected := "2"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestIgnoreCaseFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "NO"
	config.IgnoreCase = true

	input := "hello\nno match\nNO MATCH"
	expected := "-> no match <-\n-> NO MATCH <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestInvertFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"
	config.Invert = true

	input := "hello\nno match\nother line"
	expected := "-> hello <-\n-> other line <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFixedFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "n.*o"
	config.Fixed = true

	input := "hello\nn.*o literal\nno match"
	expected := "-> n.*o literal <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestLineNumberFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "no"
	config.LineNumber = true

	input := "hello\nno match\nother line\nno again"
	expected := "-> 2 no match <-\n-> 4 no again <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAfterContextFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.A = 2

	input := "line1\nThis is a match\nafter1\nafter2\nafter3\nline6"
	expected := "-> This is a match <-\nafter1\nafter2"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestBeforeContextFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.B = 2

	input := "before1\nbefore2\nThis is a match\nafter1\nline5"
	expected := "before1\nbefore2\n-> This is a match <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestContextFlag(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.C = 1

	input := "line1\nbefore\nThis is a match\nafter\nline5"
	expected := "before\n-> This is a match <-\nafter"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCombinedFlags(t *testing.T) {
	config := createConfig()
	config.Pattern = "MATCH"
	config.IgnoreCase = true
	config.LineNumber = true
	config.C = 1

	input := "line1\nbefore\nThis is a match\nafter\nline5"
	expected := "2 before\n-> 3 This is a match <-\n4 after"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultipleMatches(t *testing.T) {
	config := createConfig()
	config.Pattern = "test"
	config.A = 1

	input := "line1\ntest1\nafter1\ntest2\nafter2\nline6"
	expected := "-> test1 <-\nafter1\n-> test2 <-\nafter2"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEdgeCaseEmptyInput(t *testing.T) {
	config := createConfig()
	config.Pattern = "test"

	input := ""
	expected := ""

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEdgeCaseNoMatches(t *testing.T) {
	config := createConfig()
	config.Pattern = "nomatch"
	config.Count = true

	input := "line1\nline2\nline3"
	expected := "0"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEdgeCaseContextAtBeginning(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.B = 5

	input := "match here\nafter1\nafter2"
	expected := "-> match here <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestEdgeCaseContextAtEnd(t *testing.T) {
	config := createConfig()
	config.Pattern = "match"
	config.A = 5

	input := "before1\nbefore2\nmatch here"
	expected := "-> match here <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestRegexPattern(t *testing.T) {
	config := createConfig()
	config.Pattern = "test[0-9]+"

	input := "test\ntest123\ntest456\ntestABC"
	expected := "-> test123 <-\n-> test456 <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFixedWithIgnoreCase(t *testing.T) {
	config := createConfig()
	config.Pattern = "Test"
	config.Fixed = true
	config.IgnoreCase = true

	input := "test\nTEST\nTest\nother"
	expected := "-> test <-\n-> TEST <-\n-> Test <-"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestInvertWithCount(t *testing.T) {
	config := createConfig()
	config.Pattern = "skip"
	config.Invert = true
	config.Count = true

	input := "line1\nskip this\nline2\nline3"
	expected := "3"

	result, err := captureOutput(config, input)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
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
