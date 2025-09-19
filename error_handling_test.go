// coverage_improvement_test.go - Tests to push coverage to 95%+
package xmuslogger

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// Test error handling paths that might be uncovered
func TestErrorHandlingPaths(t *testing.T) {
	t.Run("RemoteWriter_FlushError", func(t *testing.T) {
		mockRemote := &mockRemoteWriter{
			flushError: errors.New("flush failed"),
		}
		logger := New().Remote(mockRemote)

		err := logger.Flush()
		if err == nil {
			t.Error("Expected flush error, got nil")
		}
		if err.Error() != "flush failed" {
			t.Errorf("Expected 'flush failed', got %v", err)
		}
	})

	t.Run("RemoteWriter_CloseError", func(t *testing.T) {
		mockRemote := &mockRemoteWriter{
			closeError: errors.New("close failed"),
		}
		logger := New().Remote(mockRemote)

		err := logger.Close()
		if err == nil {
			t.Error("Expected close error, got nil")
		}
	})

	t.Run("NoRemoteWriter_FlushClose", func(t *testing.T) {
		logger := New() // No remote writer

		if err := logger.Flush(); err != nil {
			t.Errorf("Flush with no remote writer should return nil, got %v", err)
		}

		if err := logger.Close(); err != nil {
			t.Errorf("Close with no remote writer should return nil, got %v", err)
		}
	})
}

// Test edge cases that might be uncovered
func TestEdgeCasesForCoverage(t *testing.T) {
	t.Run("NilEvent_AllMethods", func(t *testing.T) {
		logger := New().Level(ErrorLevel) // High level to get nil events

		// All these should be safe with nil event
		event := logger.Debug() // Should return nil
		if event != nil {
			t.Error("Expected nil event for disabled level")
		}

		// These should all be no-ops and not panic
		event.Str("key", "value")
		event.Int("num", 42)
		event.Bool("flag", true)
		event.Err(errors.New("test"))
		event.Msg("message")
		event.Msgf("formatted %s", "message")
		event.Send()
	})

	t.Run("EmptyContext", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New().Output(&buf)

		// Context with no fields
		emptyCtx := logger.With().Logger()
		emptyCtx.Info().Msg("empty context")

		if buf.Len() == 0 {
			t.Error("Expected output from empty context logger")
		}
	})

	t.Run("MultipleOutputWriters", func(t *testing.T) {
		var buf1, buf2 bytes.Buffer
		logger := New()

		// Manually set multiple writers (if supported)
		logger.mu.Lock()
		logger.writers = []io.Writer{&buf1, &buf2}
		logger.mu.Unlock()

		logger.Info().Msg("multi-writer test")

		if buf1.Len() == 0 || buf2.Len() == 0 {
			t.Error("Expected output to both writers")
		}
	})
}

// Test serializer edge cases
func TestSerializerEdgeCases(t *testing.T) {
	t.Run("AppendBytes_NoTrailingComma", func(t *testing.T) {
		buf := []byte(`"existing":"field"`)
		result := appendBytes(buf, []byte(`"new":"value"`))
		expected := `"existing":"field""new":"value"`

		if string(result) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result))
		}
	})

	t.Run("AppendBytes_WithTrailingComma", func(t *testing.T) {
		buf := []byte(`"existing":"field",`)
		result := appendBytes(buf, []byte(`"new":"value"`))
		expected := `"existing":"field","new":"value"`

		if string(result) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result))
		}
	})

	t.Run("WrapJSON_EmptyBuffer", func(t *testing.T) {
		var buf []byte
		result := wrapJSON(buf)
		expected := "{}"

		if string(result) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result))
		}
	})

	t.Run("WrapJSON_NoTrailingComma", func(t *testing.T) {
		buf := []byte(`"key":"value"`)
		result := wrapJSON(buf)
		expected := `{"key":"value"}`

		if string(result) != expected {
			t.Errorf("Expected %s, got %s", expected, string(result))
		}
	})
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("UpdateEnabledLevels", func(t *testing.T) {
		logger := New()

		// Test all levels
		levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}

		for _, level := range levels {
			logger.level = level
			logger.updateEnabledLevels()

			// Verify enabled array is correct
			for i, enabled := range logger.enabled {
				expected := Level(i) >= level
				if enabled != expected {
					t.Errorf("Level %d with threshold %d: expected %v, got %v", i, level, expected, enabled)
				}
			}
		}
	})
}

// Test global functions for complete coverage
func TestGlobalFunctionsCoverage(t *testing.T) {
	// Capture current default
	oldDefault := Default()
	defer SetDefault(oldDefault)

	// Test all global functions
	var buf bytes.Buffer
	testLogger := New().Output(&buf)
	SetDefault(testLogger)

	// Test all global print functions
	Print("test print")
	Printf("test printf %s", "arg")
	Println("test println")

	// Test configuration functions
	SetOutput(&buf)
	SetFlags(0)
	SetPrefix("TEST: ")

	// Verify output was generated
	if buf.Len() == 0 {
		t.Error("Expected output from global functions")
	}

	// Test Default() returns our test logger
	if Default() != testLogger {
		t.Error("Default() should return the logger we set")
	}
}

// Test findMessageStart edge cases

// Test HTTP option functions edge cases
func TestHTTPOptionEdgeCases(t *testing.T) {
	t.Run("WithHTTPAuth_SpecialCharacters", func(t *testing.T) {
		specialToken := "token-with-!@#$%^&*()_+"
		writer := NewHTTPRemoteWriter("https://example.com", WithHTTPAuth(specialToken))

		expected := "Bearer " + specialToken
		if writer.headers["Authorization"] != expected {
			t.Errorf("Expected Authorization %q, got %q", expected, writer.headers["Authorization"])
		}
	})

	t.Run("WithHTTPHeaders_OverwriteMultiple", func(t *testing.T) {
		headers1 := map[string]string{"A": "1", "B": "2"}
		headers2 := map[string]string{"B": "3", "C": "4"}
		headers3 := map[string]string{"A": "5", "D": "6"}

		writer := NewHTTPRemoteWriter("https://example.com",
			WithHTTPHeaders(headers1),
			WithHTTPHeaders(headers2),
			WithHTTPHeaders(headers3),
		)

		expected := map[string]string{
			"A": "5", // Overwritten by headers3
			"B": "3", // From headers2
			"C": "4", // From headers2
			"D": "6", // From headers3
		}

		for k, v := range expected {
			if writer.headers[k] != v {
				t.Errorf("Header %s: expected %s, got %s", k, v, writer.headers[k])
			}
		}
	})
}

// Test event pool edge cases
func TestEventPoolEdgeCases(t *testing.T) {
	t.Run("EventPool_LargeBuffer", func(t *testing.T) {
		e := getEvent()

		// Create a buffer larger than the pool threshold
		e.buf = make([]byte, 1<<17) // Larger than 1<<16

		// This should not panic and should not pool the oversized buffer
		putEvent(e)

		// Get a new event - should be a fresh one, not the oversized one
		e2 := getEvent()
		if cap(e2.buf) > 1<<16 {
			t.Error("Event pool returned oversized buffer")
		}
	})

	t.Run("EventPool_NormalSize", func(t *testing.T) {
		e := getEvent()
		originalCap := cap(e.buf)

		// Use normal sized buffer
		e.buf = e.buf[:10]
		putEvent(e)

		// Should be pooled and reused
		e2 := getEvent()
		if cap(e2.buf) < originalCap {
			t.Error("Event pool didn't reuse normal-sized buffer")
		}
	})
}
