// json_escaping_test.go
package xmuslogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestJSONEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Newline",
			input:    "line1\nline2",
			expected: "line1\\nline2",
		},
		{
			name:     "Tab",
			input:    "col1\tcol2",
			expected: "col1\\tcol2",
		},
		{
			name:     "Carriage Return",
			input:    "line1\rline2",
			expected: "line1\\rline2",
		},
		{
			name:     "Double Quote",
			input:    `say "hello"`,
			expected: `say \"hello\"`,
		},
		{
			name:     "Backslash",
			input:    `path\to\file`,
			expected: `path\\to\\file`,
		},
		{
			name:     "Mixed Special Characters",
			input:    "Special chars: \n\t\r\"'\\",
			expected: "Special chars: \\n\\t\\r\\\"'\\\\",
		},
		{
			name:     "Unicode",
			input:    "Unicode: 你好世界 🌍",
			expected: "Unicode: 你好世界 🌍", // Unicode should be preserved
		},
		{
			name:     "Control Characters",
			input:    "Control: \x01\x02\x1F",
			expected: "Control: \\u0001\\u0002\\u001f",
		},
		{
			name:     "Normal String",
			input:    "normal string",
			expected: "normal string",
		},
		{
			name:     "Empty String",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New().Output(&buf)

			logger.Info().Str("test_field", tt.input).Msg("test message")

			output := buf.String()

			// Parse the JSON to ensure it's valid
			var logEntry map[string]interface{}
			if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
				t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
			}

			// Verify the field was properly parsed back
			if testField, ok := logEntry["test_field"]; !ok {
				t.Error("test_field not found in output")
			} else if testField != tt.input {
				t.Errorf("Expected test_field=%q, got %q", tt.input, testField)
			}

			// Also verify the raw JSON contains the expected escaped sequence
			if tt.expected != tt.input {
				if !strings.Contains(output, tt.expected) {
					t.Errorf("Expected output to contain escaped sequence %q, but got: %s", tt.expected, output)
				}
			}
		})
	}
}

func TestComplexSpecialCharacterMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	complexMessage := `Complex message with:
	- Newlines
	- Tabs:	here
	- Quotes: "double" and 'single'
	- Backslashes: C:\path\to\file
	- Unicode: 你好 🌍
	- Control chars: ` + "\x01\x02"

	logger.Error().
		Str("complex", complexMessage).
		Str("quotes", `"quoted value"`).
		Str("backslashes", `C:\Windows\System32`).
		Msg("Complex special character test")

	output := buf.String()

	// Should be valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse complex JSON output: %v\nOutput: %s", err, output)
	}

	// Verify all fields are correctly parsed
	if complex, ok := logEntry["complex"]; !ok {
		t.Error("complex field not found")
	} else if complex != complexMessage {
		t.Errorf("Complex field not preserved correctly")
	}

	if quotes, ok := logEntry["quotes"]; !ok {
		t.Error("quotes field not found")
	} else if quotes != `"quoted value"` {
		t.Errorf("Quotes field not preserved correctly, got: %v", quotes)
	}

	if backslashes, ok := logEntry["backslashes"]; !ok {
		t.Error("backslashes field not found")
	} else if backslashes != `C:\Windows\System32` {
		t.Errorf("Backslashes field not preserved correctly, got: %v", backslashes)
	}
}

func TestErrorMessageWithSpecialChars(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	specialError := "Database error:\nConnection failed on host \"localhost\"\nPath: C:\\data\\db"
	testErr := fmt.Errorf(specialError)

	logger.Error().Err(testErr).Msg("Error with special characters")

	output := buf.String()

	// Should be valid JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse error JSON output: %v\nOutput: %s", err, output)
	}

	// Verify error field is correctly parsed
	if errorField, ok := logEntry["error"]; !ok {
		t.Error("error field not found")
	} else if errorField != specialError {
		t.Errorf("Error field not preserved correctly.\nExpected: %q\nGot: %q", specialError, errorField)
	}
}

// Benchmark to ensure escaping doesn't hurt performance too much
func BenchmarkJSONEscaping(b *testing.B) {
	logger := New().Output(&bytes.Buffer{})
	specialMessage := "Message with\nnewlines\tand\rspecial\"chars\\and\x01control"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("special", specialMessage).Msg("benchmark")
	}
}

// Test the serializer functions directly
func TestSerializerEscaping(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		contains []string // Substrings that should be in the output
	}{
		{
			key:      "newline",
			value:    "line1\nline2",
			contains: []string{`"newline"`, `"line1\nline2"`},
		},
		{
			key:      "quotes",
			value:    `say "hello"`,
			contains: []string{`"quotes"`, `"say \"hello\""`},
		},
		{
			key:      "backslash",
			value:    `path\file`,
			contains: []string{`"backslash"`, `"path\\file"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			var buf []byte
			buf = appendString(buf, tt.key, tt.value)
			result := string(buf)

			// Verify all expected substrings are present
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, but got: %s", expected, result)
				}
			}

			// Verify it produces valid JSON when wrapped
			jsonResult := string(wrapJSON(buf[:len(buf)-1])) // Remove trailing comma
			var parsed map[string]interface{}
			if err := json.Unmarshal([]byte(jsonResult), &parsed); err != nil {
				t.Errorf("Produced invalid JSON: %v\nJSON: %s", err, jsonResult)
			}
		})
	}
}
