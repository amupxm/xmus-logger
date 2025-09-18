package xmuslogger

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"testing"
)

// Test to verify context isolation
func TestContextIsolation(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	// Create context logger
	ctxLogger := logger.With().Str("service", "test-service").Int("version", 2).Logger()

	// Log with original logger (should NOT have context)
	logger.Info().Msg("Original logger message")

	// Log with context logger (should HAVE context)
	ctxLogger.Info().Msg("Context logger message")

	// Log with original logger again (should STILL not have context)
	logger.Info().Msg("Original logger message 2")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Fatalf("Expected 3 log lines, got %d", len(lines))
	}

	// Parse each line
	var logs []map[string]interface{}
	for i, line := range lines {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Fatalf("Failed to parse log line %d: %v", i, err)
		}
		logs = append(logs, logEntry)
	}

	// Test line 1: Original logger (no context)
	if _, hasService := logs[0]["service"]; hasService {
		t.Error("Original logger should NOT have 'service' field")
	}
	if _, hasVersion := logs[0]["version"]; hasVersion {
		t.Error("Original logger should NOT have 'version' field")
	}
	if logs[0]["message"] != "Original logger message" {
		t.Errorf("Expected 'Original logger message', got %v", logs[0]["message"])
	}

	// Test line 2: Context logger (has context)
	if logs[1]["service"] != "test-service" {
		t.Errorf("Context logger should have service='test-service', got %v", logs[1]["service"])
	}
	if logs[1]["version"] != float64(2) { // JSON numbers are float64
		t.Errorf("Context logger should have version=2, got %v", logs[1]["version"])
	}
	if logs[1]["message"] != "Context logger message" {
		t.Errorf("Expected 'Context logger message', got %v", logs[1]["message"])
	}

	// Test line 3: Original logger again (still no context)
	if _, hasService := logs[2]["service"]; hasService {
		t.Error("Original logger (second use) should NOT have 'service' field")
	}
	if _, hasVersion := logs[2]["version"]; hasVersion {
		t.Error("Original logger (second use) should NOT have 'version' field")
	}
	if logs[2]["message"] != "Original logger message 2" {
		t.Errorf("Expected 'Original logger message 2', got %v", logs[2]["message"])
	}
}

// Test multiple context loggers from same parent
func TestMultipleContextLoggers(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	// Create multiple context loggers
	dbLogger := logger.With().Str("component", "database").Logger()
	httpLogger := logger.With().Str("component", "http").Int("port", 8080).Logger()
	cacheLogger := logger.With().Str("component", "cache").Bool("enabled", true).Logger()

	// Log with each
	logger.Info().Msg("Main application")
	dbLogger.Info().Msg("Database connected")
	httpLogger.Info().Msg("HTTP server started")
	cacheLogger.Info().Msg("Cache initialized")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Fatalf("Expected 4 log lines, got %d", len(lines))
	}

	// Parse all lines
	var logs []map[string]interface{}
	for i, line := range lines {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Fatalf("Failed to parse log line %d: %v", i, err)
		}
		logs = append(logs, logEntry)
	}

	// Test main logger (no context)
	if _, hasComponent := logs[0]["component"]; hasComponent {
		t.Error("Main logger should not have component field")
	}

	// Test DB logger
	if logs[1]["component"] != "database" {
		t.Errorf("DB logger should have component='database', got %v", logs[1]["component"])
	}
	if _, hasPort := logs[1]["port"]; hasPort {
		t.Error("DB logger should not have port field from HTTP logger")
	}

	// Test HTTP logger
	if logs[2]["component"] != "http" {
		t.Errorf("HTTP logger should have component='http', got %v", logs[2]["component"])
	}
	if logs[2]["port"] != float64(8080) {
		t.Errorf("HTTP logger should have port=8080, got %v", logs[2]["port"])
	}
	if _, hasEnabled := logs[2]["enabled"]; hasEnabled {
		t.Error("HTTP logger should not have enabled field from cache logger")
	}

	// Test Cache logger
	if logs[3]["component"] != "cache" {
		t.Errorf("Cache logger should have component='cache', got %v", logs[3]["component"])
	}
	if logs[3]["enabled"] != true {
		t.Errorf("Cache logger should have enabled=true, got %v", logs[3]["enabled"])
	}
	if _, hasPort := logs[3]["port"]; hasPort {
		t.Error("Cache logger should not have port field from HTTP logger")
	}
}

// Test nested context chaining
func TestNestedContextChaining(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	// Create nested context
	serviceLogger := logger.With().Str("service", "api").Logger()
	requestLogger := serviceLogger.With().Str("request_id", "req-123").Logger()
	userLogger := requestLogger.With().Str("user_id", "user-456").Logger()

	// Log at each level
	logger.Info().Msg("Application started")
	serviceLogger.Info().Msg("Service initialized")
	requestLogger.Info().Msg("Processing request")
	userLogger.Info().Msg("User authenticated")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 4 {
		t.Fatalf("Expected 4 log lines, got %d", len(lines))
	}

	// Parse all lines
	var logs []map[string]interface{}
	for i, line := range lines {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Fatalf("Failed to parse log line %d: %v", i, err)
		}
		logs = append(logs, logEntry)
	}

	// Test application level (no context)
	if _, hasService := logs[0]["service"]; hasService {
		t.Error("Application logger should not have service field")
	}

	// Test service level (only service)
	if logs[1]["service"] != "api" {
		t.Errorf("Service logger should have service='api', got %v", logs[1]["service"])
	}
	if _, hasRequestID := logs[1]["request_id"]; hasRequestID {
		t.Error("Service logger should not have request_id field")
	}

	// Test request level (service + request_id)
	if logs[2]["service"] != "api" {
		t.Errorf("Request logger should have service='api', got %v", logs[2]["service"])
	}
	if logs[2]["request_id"] != "req-123" {
		t.Errorf("Request logger should have request_id='req-123', got %v", logs[2]["request_id"])
	}
	if _, hasUserID := logs[2]["user_id"]; hasUserID {
		t.Error("Request logger should not have user_id field")
	}

	// Test user level (all fields)
	if logs[3]["service"] != "api" {
		t.Errorf("User logger should have service='api', got %v", logs[3]["service"])
	}
	if logs[3]["request_id"] != "req-123" {
		t.Errorf("User logger should have request_id='req-123', got %v", logs[3]["request_id"])
	}
	if logs[3]["user_id"] != "user-456" {
		t.Errorf("User logger should have user_id='user-456', got %v", logs[3]["user_id"])
	}
}

// Test that demonstrates the bug before the fix
func TestContextIsolationBug_Demo(t *testing.T) {
	// This test demonstrates what would happen with the buggy implementation
	// Skip if the fix is already applied
	t.Skip("This test demonstrates the bug - skip if fix is applied")

	var buf bytes.Buffer
	logger := New().Output(&buf)

	// This would cause the bug in the original implementation
	ctxLogger := logger.With().Str("leaked", "value").Logger()

	// With the bug, this would unexpectedly have the "leaked" field
	logger.Info().Msg("Should not have leaked field")

	// This is expected behavior - should have the field
	ctxLogger.Info().Msg("Should have leaked field")

	// Parse and check (this test would fail with the original buggy code)
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	var firstLog map[string]interface{}
	json.Unmarshal([]byte(lines[0]), &firstLog)

	// This assertion would fail with the original bug
	if _, hasLeaked := firstLog["leaked"]; hasLeaked {
		t.Error("BUG DETECTED: Original logger has leaked context field!")
	}
}

// Performance test to ensure context isolation doesn't hurt performance
func BenchmarkContextCreation(b *testing.B) {
	logger := New().Output(&SafeBuffer{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctxLogger := logger.With().Str("key", "value").Int("iteration", i).Logger()
		ctxLogger.Info().Msg("benchmark")
	}
}

// Test concurrent context usage (should be safe after fix)
func TestConcurrentContextUsage(t *testing.T) {
	safeBuf := &SafeBuffer{}
	logger := New().Output(safeBuf)

	const numGoroutines = 50
	const messagesPerGoroutine = 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine creates its own context logger
			ctxLogger := logger.With().
				Int("goroutine_id", id). // Fix: Use Int instead of string conversion
				Str("component", "concurrent-test").
				Logger()

			for j := 0; j < messagesPerGoroutine; j++ {
				ctxLogger.Info().Int("message", j).Msg("concurrent context test")
			}
		}(i)
	}

	wg.Wait()

	output := safeBuf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	expectedLines := numGoroutines * messagesPerGoroutine
	actualLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			actualLines++
		}
	}

	if actualLines != expectedLines {
		t.Errorf("Expected %d log lines, got %d", expectedLines, actualLines)
		t.Logf("Output preview (first 500 chars): %s", output[:min(500, len(output))])
	}

	// Verify all lines have proper context isolation
	validContextLines := 0
	goroutineIDs := make(map[float64]int) // Track which goroutines we see

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			t.Errorf("Failed to parse log line %d: %v\nLine content: %s", i, err, line)
			continue
		}

		// All lines should have the goroutine_id context
		goroutineIDRaw, hasGoroutineID := logEntry["goroutine_id"]
		if !hasGoroutineID {
			t.Errorf("Context logger line %d missing goroutine_id field: %v", i, logEntry)
			continue
		}

		// Verify goroutine_id is a valid number
		goroutineID, ok := goroutineIDRaw.(float64) // JSON numbers are float64
		if !ok {
			t.Errorf("goroutine_id should be a number, got %T: %v", goroutineIDRaw, goroutineIDRaw)
			continue
		}

		if goroutineID < 0 || goroutineID >= float64(numGoroutines) {
			t.Errorf("Invalid goroutine_id: %v (should be 0-%d)", goroutineID, numGoroutines-1)
			continue
		}

		// Count occurrences per goroutine
		goroutineIDs[goroutineID]++

		// All lines should have component field
		if component, hasComponent := logEntry["component"]; !hasComponent {
			t.Errorf("Context logger line %d missing component field: %v", i, logEntry)
			continue
		} else if component != "concurrent-test" {
			t.Errorf("Expected component='concurrent-test', got %v", component)
			continue
		}

		// All lines should have message field
		if _, hasMessage := logEntry["message"]; !hasMessage {
			t.Errorf("Context logger line %d missing message field: %v", i, logEntry)
			continue
		}

		validContextLines++
	}

	if validContextLines != actualLines {
		t.Errorf("Expected %d lines with proper context, got %d", actualLines, validContextLines)
	}

	// Verify we got messages from all goroutines
	if len(goroutineIDs) != numGoroutines {
		t.Errorf("Expected messages from %d goroutines, got %d", numGoroutines, len(goroutineIDs))
	}

	// Verify each goroutine produced the expected number of messages
	for goroutineID, count := range goroutineIDs {
		if count != messagesPerGoroutine {
			t.Errorf("Goroutine %v produced %d messages, expected %d", goroutineID, count, messagesPerGoroutine)
		}
	}

	t.Logf("Successfully verified context isolation across %d concurrent loggers with %d total messages",
		len(goroutineIDs), validContextLines)
}
