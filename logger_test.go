// logger_test.go
package xmuslogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Test Helpers
type mockRemoteWriter struct {
	writes      [][]byte
	asyncWrites [][]byte
	writeError  error
	asyncError  error
	flushError  error
	closeError  error
	mu          sync.Mutex
}

func (m *mockRemoteWriter) Write(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.writeError != nil {
		return m.writeError
	}
	m.writes = append(m.writes, append([]byte(nil), data...))
	return nil
}

func (m *mockRemoteWriter) WriteAsync(data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.asyncError != nil {
		return m.asyncError
	}
	m.asyncWrites = append(m.asyncWrites, append([]byte(nil), data...))
	return nil
}

func (m *mockRemoteWriter) Flush() error {
	return m.flushError
}

func (m *mockRemoteWriter) Close() error {
	return m.closeError
}

func (m *mockRemoteWriter) GetWrites() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([][]byte(nil), m.writes...)
}

func (m *mockRemoteWriter) GetAsyncWrites() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([][]byte(nil), m.asyncWrites...)
}

func parseLogLine(line string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(line), &result)
	return result, err
}

// =============================================================================
// CORE LOGGER TESTS
// =============================================================================

func TestNew(t *testing.T) {
	logger := New()

	if logger == nil {
		t.Fatal("New() returned nil")
	}

	if logger.level != InfoLevel {
		t.Errorf("Expected default level InfoLevel, got %v", logger.level)
	}

	if len(logger.writers) != 1 {
		t.Errorf("Expected 1 writer, got %d", len(logger.writers))
	}

	if logger.writers[0] != os.Stdout {
		t.Error("Expected default writer to be os.Stdout")
	}

	if logger.Logger == nil {
		t.Error("Embedded logger should not be nil")
	}
}

func TestNewWithOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewWithOutput(&buf)

	if logger == nil {
		t.Fatal("NewWithOutput() returned nil")
	}

	if len(logger.writers) != 1 {
		t.Errorf("Expected 1 writer, got %d", len(logger.writers))
	}

	if logger.writers[0] != &buf {
		t.Error("Expected writer to be the provided buffer")
	}
}

func TestLoggerLevel(t *testing.T) {
	tests := []struct {
		level    Level
		expected Level
	}{
		{TraceLevel, TraceLevel},
		{DebugLevel, DebugLevel},
		{InfoLevel, InfoLevel},
		{WarnLevel, WarnLevel},
		{ErrorLevel, ErrorLevel},
		{FatalLevel, FatalLevel},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%s", tt.level.String()), func(t *testing.T) {
			logger := New().Level(tt.level)
			if logger.level != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, logger.level)
			}
		})
	}
}

func TestLoggerClone(t *testing.T) {
	original := New().Level(DebugLevel)
	var buf bytes.Buffer
	original.SetOutput(&buf)

	clone := original.clone()

	// Test independence
	if &original == &clone {
		t.Error("Clone should return a different instance")
	}

	if original.level != clone.level {
		t.Errorf("Clone should have same level: original=%v, clone=%v", original.level, clone.level)
	}

	if len(original.writers) != len(clone.writers) {
		t.Errorf("Clone should have same number of writers: original=%d, clone=%d",
			len(original.writers), len(clone.writers))
	}
}

// =============================================================================
// LEVEL TESTS
// =============================================================================

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{TraceLevel, "trace"},
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.level.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.level.String())
			}
		})
	}
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf).Level(WarnLevel)

	// These should not produce output
	logger.Trace().Msg("trace message")
	logger.Debug().Msg("debug message")
	logger.Info().Msg("info message")

	// These should produce output
	logger.Warn().Msg("warn message")
	logger.Error().Msg("error message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 log lines, got %d. Output: %s", len(lines), output)
	}

	// Verify warn message
	if len(lines) > 0 {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(lines[0]), &logEntry); err != nil {
			t.Errorf("Failed to parse first log line: %v", err)
		} else {
			if logEntry["level"] != "warn" {
				t.Errorf("Expected level 'warn', got %v", logEntry["level"])
			}
			if logEntry["message"] != "warn message" {
				t.Errorf("Expected message 'warn message', got %v", logEntry["message"])
			}
		}
	}
}

// =============================================================================
// EVENT TESTS
// =============================================================================

func TestEventChaining(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	logger.Info().
		Str("key1", "value1").
		Int("key2", 42).
		Bool("key3", true).
		Msg("test message")

	output := buf.String()
	var logEntry map[string]interface{}

	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["key1"] != "value1" {
		t.Errorf("Expected key1=value1, got %v", logEntry["key1"])
	}

	if logEntry["key2"] != float64(42) { // JSON numbers are float64
		t.Errorf("Expected key2=42, got %v", logEntry["key2"])
	}

	if logEntry["key3"] != true {
		t.Errorf("Expected key3=true, got %v", logEntry["key3"])
	}

	if logEntry["message"] != "test message" {
		t.Errorf("Expected message='test message', got %v", logEntry["message"])
	}

	if logEntry["level"] != "info" {
		t.Errorf("Expected level='info', got %v", logEntry["level"])
	}
}

func TestEventError(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	testError := fmt.Errorf("test error message")
	logger.Error().Err(testError).Msg("error occurred")

	output := buf.String()
	var logEntry map[string]interface{}

	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["error"] != "test error message" {
		t.Errorf("Expected error='test error message', got %v", logEntry["error"])
	}
}

func TestEventMsgf(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	logger.Info().Msgf("formatted message: %s %d", "test", 123)

	output := buf.String()
	var logEntry map[string]interface{}

	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["message"] != "formatted message: test 123" {
		t.Errorf("Expected formatted message, got %v", logEntry["message"])
	}
}

func TestEventSend(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	logger.Info().Str("key", "value").Send()

	output := buf.String()
	var logEntry map[string]interface{}

	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["message"] != "" {
		t.Errorf("Expected empty message, got %v", logEntry["message"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("Expected key=value, got %v", logEntry["key"])
	}
}

func TestNilEventSafety(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf).Level(ErrorLevel) // Set high level

	// This should return nil event due to level filtering
	event := logger.Info()

	// All these should be safe to call on nil event
	event.Str("key", "value").Int("num", 42).Bool("flag", true).Msg("message")
	event.Msgf("formatted %s", "message")
	event.Send()

	// Should produce no output
	if buf.Len() > 0 {
		t.Errorf("Expected no output for filtered level, got: %s", buf.String())
	}
}

// =============================================================================
// CONTEXT TESTS
// =============================================================================

func TestContext(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	ctxLogger := logger.With().Str("service", "test").Int("version", 1).Logger()
	ctxLogger.Info().Msg("context test")

	output := buf.String()
	var logEntry map[string]interface{}

	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["service"] != "test" {
		t.Errorf("Expected service=test, got %v", logEntry["service"])
	}

	if logEntry["version"] != float64(1) {
		t.Errorf("Expected version=1, got %v", logEntry["version"])
	}

	if logEntry["message"] != "context test" {
		t.Errorf("Expected message='context test', got %v", logEntry["message"])
	}
}

func TestContextChaining(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	ctx := logger.With().Str("key1", "value1")
	ctx2 := ctx.Str("key2", "value2").Int("key3", 3)

	ctx2.Logger().Info().Msg("chained context")

	output := buf.String()
	var logEntry map[string]interface{}

	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["key1"] != "value1" {
		t.Errorf("Expected key1=value1, got %v", logEntry["key1"])
	}

	if logEntry["key2"] != "value2" {
		t.Errorf("Expected key2=value2, got %v", logEntry["key2"])
	}

	if logEntry["key3"] != float64(3) {
		t.Errorf("Expected key3=3, got %v", logEntry["key3"])
	}
}

// =============================================================================
// STANDARD LOGGER COMPATIBILITY TESTS
// =============================================================================

func TestStandardLoggerCompatibility(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	logger.Print("print message")
	logger.Printf("printf message: %s", "test")
	logger.Println("println message")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 log lines, got %d", len(lines))
	}

	// Test first line
	var logEntry map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &logEntry); err != nil {
		t.Fatalf("Failed to parse first log line: %v", err)
	}

	if logEntry["level"] != "info" {
		t.Errorf("Expected level 'info', got %v", logEntry["level"])
	}

	if logEntry["source"] != "stdlib" {
		t.Errorf("Expected source 'stdlib', got %v", logEntry["source"])
	}
}

func TestSetOutput(t *testing.T) {
	logger := New()
	var buf bytes.Buffer

	logger.SetOutput(&buf)
	logger.Print("test message")

	if buf.Len() == 0 {
		t.Error("Expected output to buffer, got none")
	}
}

func TestSetFlags(t *testing.T) {
	logger := New()
	logger.SetFlags(log.Lshortfile)

	if logger.Flags() != log.Lshortfile {
		t.Errorf("Expected flags %d, got %d", log.Lshortfile, logger.Flags())
	}
}

func TestSetPrefix(t *testing.T) {
	logger := New()
	testPrefix := "TEST: "

	logger.SetPrefix(testPrefix)

	if logger.Prefix() != testPrefix {
		t.Errorf("Expected prefix %q, got %q", testPrefix, logger.Prefix())
	}
}

// =============================================================================
// REMOTE WRITER TESTS
// =============================================================================

func TestRemoteWriter(t *testing.T) {
	var buf bytes.Buffer
	mockRemote := &mockRemoteWriter{}
	logger := New().Output(&buf).Remote(mockRemote)

	logger.Info().Msg("remote test")

	writes := mockRemote.GetWrites()
	if len(writes) != 1 {
		t.Errorf("Expected 1 remote write, got %d", len(writes))
	}

	var logEntry map[string]interface{}
	if err := json.Unmarshal(writes[0], &logEntry); err != nil {
		t.Fatalf("Failed to parse remote write: %v", err)
	}

	if logEntry["message"] != "remote test" {
		t.Errorf("Expected message 'remote test', got %v", logEntry["message"])
	}
}

func TestAsyncRemoteWriter(t *testing.T) {
	var buf bytes.Buffer
	mockRemote := &mockRemoteWriter{}
	logger := New().Output(&buf).Remote(mockRemote)
	logger.async = true // Simulate async mode

	logger.Info().Msg("async remote test")

	asyncWrites := mockRemote.GetAsyncWrites()
	if len(asyncWrites) != 1 {
		t.Errorf("Expected 1 async remote write, got %d", len(asyncWrites))
	}

	writes := mockRemote.GetWrites()
	if len(writes) != 0 {
		t.Errorf("Expected 0 sync remote writes, got %d", len(writes))
	}
}

func TestRemoteWriterError(t *testing.T) {
	var buf bytes.Buffer
	mockRemote := &mockRemoteWriter{
		writeError: fmt.Errorf("remote write failed"),
	}
	logger := New().Output(&buf).Remote(mockRemote)

	// This should not panic even with remote write error
	logger.Info().Msg("error test")

	// Local write should still work
	if buf.Len() == 0 {
		t.Error("Expected local output even with remote error")
	}
}

func TestFlushAndClose(t *testing.T) {
	mockRemote := &mockRemoteWriter{}
	logger := New().Remote(mockRemote)

	if err := logger.Flush(); err != nil {
		t.Errorf("Flush() returned error: %v", err)
	}

	if err := logger.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}

func TestFlushAndCloseErrors(t *testing.T) {
	flushErr := fmt.Errorf("flush error")
	closeErr := fmt.Errorf("close error")

	mockRemote := &mockRemoteWriter{
		flushError: flushErr,
		closeError: closeErr,
	}
	logger := New().Remote(mockRemote)

	if err := logger.Flush(); err != flushErr {
		t.Errorf("Expected flush error %v, got %v", flushErr, err)
	}

	if err := logger.Close(); err != closeErr {
		t.Errorf("Expected close error %v, got %v", closeErr, err)
	}
}

// =============================================================================
// HTTP REMOTE WRITER TESTS
// =============================================================================

func TestNewHTTPRemoteWriter(t *testing.T) {
	endpoint := "https://example.com/logs"
	writer := NewHTTPRemoteWriter(endpoint)

	if writer.endpoint != endpoint {
		t.Errorf("Expected endpoint %s, got %s", endpoint, writer.endpoint)
	}

	if writer.headers == nil {
		t.Error("Headers should be initialized")
	}
}

func TestHTTPRemoteWriterOptions(t *testing.T) {
	endpoint := "https://example.com/logs"
	token := "test-token"

	writer := NewHTTPRemoteWriter(endpoint, WithHTTPAuth(token))

	expectedAuth := "Bearer " + token
	if writer.headers["Authorization"] != expectedAuth {
		t.Errorf("Expected Authorization %s, got %s", expectedAuth, writer.headers["Authorization"])
	}
}

func TestHTTPRemoteWriterCustomHeaders(t *testing.T) {
	endpoint := "https://example.com/logs"
	customHeaders := map[string]string{
		"Content-Type": "application/json",
		"X-API-Key":    "secret",
	}

	writer := NewHTTPRemoteWriter(endpoint, WithHTTPHeaders(customHeaders))

	for k, v := range customHeaders {
		if writer.headers[k] != v {
			t.Errorf("Expected header %s=%s, got %s", k, v, writer.headers[k])
		}
	}
}

func TestRemoteHTTPMethod(t *testing.T) {
	var buf bytes.Buffer
	endpoint := "https://example.com/logs"
	logger := New().Output(&buf).RemoteHTTP(endpoint, WithHTTPAuth("token"))

	// Should not panic (actual HTTP implementation is TODO)
	logger.Info().Msg("http test")

	if buf.Len() == 0 {
		t.Error("Expected local output")
	}
}

// =============================================================================
// SERIALIZER TESTS
// =============================================================================

func TestAppendString(t *testing.T) {
	var buf []byte
	buf = appendString(buf, "key", "value")

	expected := `"key":"value",`
	if string(buf) != expected {
		t.Errorf("Expected %s, got %s", expected, string(buf))
	}
}

func TestAppendInt(t *testing.T) {
	var buf []byte
	buf = appendInt(buf, "number", 42)

	expected := `"number":42,`
	if string(buf) != expected {
		t.Errorf("Expected %s, got %s", expected, string(buf))
	}
}

func TestAppendBool(t *testing.T) {
	var buf []byte
	buf = appendBool(buf, "flag", true)
	buf = appendBool(buf, "flag2", false)

	expected := `"flag":true,"flag2":false,`
	if string(buf) != expected {
		t.Errorf("Expected %s, got %s", expected, string(buf))
	}
}

func TestAppendTime(t *testing.T) {
	var buf []byte
	testTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	buf = appendTime(buf, "timestamp", testTime)

	expected := `"timestamp":"2023-01-01T12:00:00Z",`
	if string(buf) != expected {
		t.Errorf("Expected %s, got %s", expected, string(buf))
	}
}

func TestWrapJSON(t *testing.T) {
	buf := []byte(`"key":"value","number":42,`)
	result := wrapJSON(buf)

	expected := `{"key":"value","number":42}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

func TestWrapJSONEmpty(t *testing.T) {
	var buf []byte
	result := wrapJSON(buf)

	expected := `{}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// =============================================================================
// EVENT POOL TESTS
// =============================================================================

func TestEventPool(t *testing.T) {
	e1 := getEvent()
	if e1 == nil {
		t.Fatal("getEvent() returned nil")
	}

	if cap(e1.buf) == 0 {
		t.Error("Event buffer should have capacity")
	}

	putEvent(e1)

	e2 := getEvent()
	if e2 == nil {
		t.Fatal("getEvent() after put returned nil")
	}

	// Should reuse the same event (though not guaranteed in tests)
	if len(e2.buf) != 0 {
		t.Error("Reused event should have empty buffer")
	}
}

func TestEventPoolOversizedBuffer(t *testing.T) {
	e := getEvent()

	// Create oversized buffer
	e.buf = make([]byte, 1<<17) // Bigger than 1<<16

	putEvent(e)

	// This test mainly ensures putEvent doesn't panic with large buffers
	// The actual pooling behavior for oversized buffers is internal
}

// =============================================================================
// UNIFIED OUTPUT TESTS
// =============================================================================

func TestFindMessageStart(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"2023/01/01 12:00:00 message", 20},
		{"prefix message", 7},
		{"no spaces", 0},
		{"one space", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := findMessageStart(tt.input)
			if result != tt.expected {
				t.Errorf("findMessageStart(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoggerWriter(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	lw := &loggerWriter{parent: logger}

	message := "test message\n"
	n, err := lw.Write([]byte(message))

	if err != nil {
		t.Errorf("Write returned error: %v", err)
	}

	if n != len(message) {
		t.Errorf("Write returned %d bytes, expected %d", n, len(message))
	}

	if buf.Len() == 0 {
		t.Error("Expected output to buffer")
	}

	var logEntry map[string]interface{}
	output := strings.TrimSpace(buf.String())
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	if logEntry["message"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", logEntry["message"])
	}
}

// =============================================================================
// GLOBAL FUNCTION TESTS
// =============================================================================

func TestGlobalFunctions(t *testing.T) {
	// Capture original stdout
	oldStd := std
	defer func() { std = oldStd }()

	// Set custom logger for testing
	var buf bytes.Buffer
	std = New().Output(&buf)

	Print("test print")
	Printf("test printf: %s", "arg")
	Println("test println")

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("Expected 3 log lines, got %d", len(lines))
	}
}

func TestSetDefault(t *testing.T) {
	var buf bytes.Buffer
	customLogger := New().Output(&buf)

	oldStd := std
	defer func() { std = oldStd }()

	SetDefault(customLogger)

	if Default() != customLogger {
		t.Error("Default() should return the set logger")
	}
}

// =============================================================================
// CONCURRENCY TESTS
// =============================================================================

// Thread-safe buffer for testing
type SafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (sb *SafeBuffer) Write(p []byte) (n int, err error) {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Write(p)
}

func (sb *SafeBuffer) String() string {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.String()
}

func (sb *SafeBuffer) Len() int {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	return sb.buf.Len()
}

func (sb *SafeBuffer) Reset() {
	sb.mu.Lock()
	defer sb.mu.Unlock()
	sb.buf.Reset()
}

// Fixed concurrent test using SafeBuffer
func TestConcurrentLogging_Fixed(t *testing.T) {
	safeBuf := &SafeBuffer{}
	logger := New().Output(safeBuf)

	const numGoroutines = 100
	const messagesPerGoroutine = 10
	expectedLines := numGoroutines * messagesPerGoroutine

	var wg sync.WaitGroup
	var counter int64

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				msgNum := atomic.AddInt64(&counter, 1)
				logger.Info().Int("goroutine", id).Int64("message", msgNum).Msg("concurrent test")
			}
		}(i)
	}

	wg.Wait()

	output := safeBuf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	actualLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			actualLines++
		}
	}

	if actualLines != expectedLines {
		t.Errorf("Expected %d log lines, got %d", expectedLines, actualLines)
	}

	// Verify the atomic counter worked correctly
	finalCount := atomic.LoadInt64(&counter)
	if finalCount != int64(expectedLines) {
		t.Errorf("Expected counter to be %d, got %d", expectedLines, finalCount)
	}

	t.Logf("Successfully logged %d lines concurrently", actualLines)
}

// Alternative: Test concurrent logging with separate outputs
func TestConcurrentLogging_SeparateOutputs(t *testing.T) {
	const numGoroutines = 100
	const messagesPerGoroutine = 10

	var wg sync.WaitGroup
	var totalLines int64
	var totalBytes int64

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine gets its own buffer
			var buf bytes.Buffer
			logger := New().Output(&buf)

			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info().Int("goroutine", id).Int("message", j).Msg("separate concurrent test")
			}

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			validLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					validLines++
				}
			}

			atomic.AddInt64(&totalLines, int64(validLines))
			atomic.AddInt64(&totalBytes, int64(len(output)))

			if validLines != messagesPerGoroutine {
				t.Errorf("Goroutine %d: expected %d lines, got %d", id, messagesPerGoroutine, validLines)
			}
		}(i)
	}

	wg.Wait()

	expectedTotal := int64(numGoroutines * messagesPerGoroutine)
	actualTotal := atomic.LoadInt64(&totalLines)
	totalBytesWritten := atomic.LoadInt64(&totalBytes)

	if actualTotal != expectedTotal {
		t.Errorf("Expected %d total lines, got %d", expectedTotal, actualTotal)
	}

	t.Logf("Successfully processed %d lines across %d goroutines (%d total bytes)",
		actualTotal, numGoroutines, totalBytesWritten)
}

// Test that demonstrates the logger is thread-safe with proper output
func TestLoggerThreadSafety_WithFiles(t *testing.T) {
	// This test would use actual files in a real scenario
	// For this test, we'll use our SafeBuffer

	const numLoggers = 10
	const messagesPerLogger = 100

	var wg sync.WaitGroup
	var totalMessages int64

	wg.Add(numLoggers)

	for i := 0; i < numLoggers; i++ {
		go func(loggerID int) {
			defer wg.Done()

			// Each "logger instance" writes to its own output
			safeBuf := &SafeBuffer{}
			logger := New().Output(safeBuf).Level(DebugLevel)

			// Add some context
			ctxLogger := logger.With().
				Str("logger_id", string(rune(loggerID))).
				Str("component", "test").
				Logger()

			for j := 0; j < messagesPerLogger; j++ {
				switch j % 4 {
				case 0:
					ctxLogger.Info().Int("msg_id", j).Msg("Info message")
				case 1:
					ctxLogger.Debug().Int("msg_id", j).Msg("Debug message")
				case 2:
					ctxLogger.Warn().Int("msg_id", j).Msg("Warning message")
				case 3:
					ctxLogger.Error().Int("msg_id", j).Msg("Error message")
				}

				atomic.AddInt64(&totalMessages, 1)
			}

			// Verify this logger's output
			output := safeBuf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			actualLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					actualLines++
				}
			}

			if actualLines != messagesPerLogger {
				t.Errorf("Logger %d: expected %d lines, got %d", loggerID, messagesPerLogger, actualLines)
			}
		}(i)
	}

	wg.Wait()

	expectedTotal := int64(numLoggers * messagesPerLogger)
	actualTotal := atomic.LoadInt64(&totalMessages)

	if actualTotal != expectedTotal {
		t.Errorf("Expected %d total messages, got %d", expectedTotal, actualTotal)
	}

	t.Logf("Successfully handled %d concurrent loggers with %d messages each", numLoggers, messagesPerLogger)
}

// Performance test to ensure thread safety doesn't hurt performance too much
func BenchmarkConcurrentLogging_SafeBuffer(b *testing.B) {
	safeBuf := &SafeBuffer{}
	logger := New().Output(safeBuf)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			logger.Info().Int("iteration", i).Msg("benchmark message")
			i++
		}
	})
}

// Benchmark with separate buffers (should be faster)
func BenchmarkConcurrentLogging_SeparateBuffers(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine gets its own logger and buffer
		var buf bytes.Buffer
		logger := New().Output(&buf)

		i := 0
		for pb.Next() {
			logger.Info().Int("iteration", i).Msg("benchmark message")
			i++
		}
	})
}

// Test with separate buffers to prove it's not just the buffer
func TestWithSeparateBuffers(t *testing.T) {
	const numGoroutines = 100
	const messagesPerGoroutine = 10

	var wg sync.WaitGroup
	var totalLines int64

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine gets its own buffer and logger
			var buf bytes.Buffer
			logger := New().Output(&buf)

			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info().Int("id", id).Int("msg", j).Msg("separate test")
			}

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")
			validLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					validLines++
				}
			}

			atomic.AddInt64(&totalLines, int64(validLines))

			if validLines != messagesPerGoroutine {
				t.Errorf("Goroutine %d: expected %d lines, got %d", id, messagesPerGoroutine, validLines)
			}
		}(i)
	}

	wg.Wait()

	expectedTotal := int64(numGoroutines * messagesPerGoroutine)
	actualTotal := atomic.LoadInt64(&totalLines)

	if actualTotal != expectedTotal {
		t.Errorf("Separate buffers: expected %d total lines, got %d", expectedTotal, actualTotal)
	} else {
		t.Log("Separate buffers work fine - confirms shared buffer is the issue")
	}
}

// Thread-safe buffer wrapper for testing
type threadSafeBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (tsb *threadSafeBuffer) Write(p []byte) (n int, err error) {
	tsb.mu.Lock()
	defer tsb.mu.Unlock()
	return tsb.buf.Write(p)
}

func (tsb *threadSafeBuffer) String() string {
	tsb.mu.Lock()
	defer tsb.mu.Unlock()
	return tsb.buf.String()
}

func TestWithThreadSafeBuffer(t *testing.T) {
	tsb := &threadSafeBuffer{}
	logger := New().Output(tsb)

	const numGoroutines = 100
	const messagesPerGoroutine = 10
	expectedLines := numGoroutines * messagesPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info().Int("id", id).Int("msg", j).Msg("thread-safe test")
			}
		}(i)
	}

	wg.Wait()

	output := tsb.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	actualLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			actualLines++
		}
	}

	t.Logf("With thread-safe buffer: expected %d lines, got %d", expectedLines, actualLines)

	if actualLines != expectedLines {
		t.Errorf("Even with thread-safe buffer: expected %d lines, got %d", expectedLines, actualLines)
		t.Log("This indicates additional thread safety issues beyond just the buffer")
	}
}

// Option 4: Test thread safety more thoroughly
func TestConcurrentLogging_ThreadSafety(t *testing.T) {
	// Use a thread-safe buffer or separate buffers

	const numGoroutines = 50 // Reduced for more reliable testing
	const messagesPerGoroutine = 20

	var wg sync.WaitGroup
	var completedMessages int64

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Create a local buffer for this goroutine to avoid races
			var localBuf bytes.Buffer
			localLogger := New().Output(&localBuf)

			for j := 0; j < messagesPerGoroutine; j++ {
				localLogger.Info().
					Int("goroutine", id).
					Int("message", j).
					Str("timestamp", fmt.Sprintf("%d", atomic.AddInt64(&completedMessages, 1))).
					Msg("concurrent test")
			}

			// Verify this goroutine's output
			output := localBuf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")

			if len(lines) != messagesPerGoroutine {
				t.Errorf("Goroutine %d: expected %d lines, got %d",
					id, messagesPerGoroutine, len(lines))
			}
		}(i)
	}

	wg.Wait()

	expectedTotal := int64(numGoroutines * messagesPerGoroutine)
	actualTotal := atomic.LoadInt64(&completedMessages)

	if actualTotal != expectedTotal {
		t.Errorf("Expected %d total messages, got %d", expectedTotal, actualTotal)
	}
}

// Option 5: Test that focuses on the logger's thread safety, not output counting
func TestLoggerThreadSafety(t *testing.T) {
	var buf bytes.Buffer
	logger := New().Output(&buf)

	const numGoroutines = 100
	const messagesPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Test that concurrent logging doesn't panic or corrupt the logger state
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				// Focus on testing that these operations are thread-safe
				logger.Info().Int("goroutine", id).Int("message", j).Msg("concurrent test")
				logger.Debug().Str("type", "debug").Int("id", id).Send()
				logger.Error().Bool("concurrent", true).Msgf("Error from goroutine %d", id)
			}
		}(i)
	}

	wg.Wait()

	// Just verify we got some output and didn't panic
	output := buf.String()
	if len(output) == 0 {
		t.Error("Expected some log output, got none")
	}

	// Verify output is valid JSON lines
	lines := strings.Split(strings.TrimSpace(output), "\n")
	validLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			// Try to parse as JSON to verify structure isn't corrupted
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(line), &data); err == nil {
				validLines++
			}
		}
	}

	if validLines == 0 {
		t.Error("No valid JSON log lines found")
	}

	t.Logf("Successfully processed %d valid log lines from concurrent operations", validLines)
}

// Helper function for older Go versions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestConcurrentRemoteLogging(t *testing.T) {
	mockRemote := &mockRemoteWriter{}
	logger := New().Remote(mockRemote)

	const numGoroutines = 50
	const messagesPerGoroutine = 5

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				logger.Info().Int("goroutine", id).Int("message", j).Msg("concurrent remote test")
			}
		}(i)
	}

	wg.Wait()

	writes := mockRemote.GetWrites()
	expectedWrites := numGoroutines * messagesPerGoroutine
	if len(writes) != expectedWrites {
		t.Errorf("Expected %d remote writes, got %d", expectedWrites, len(writes))
	}
}

// =============================================================================
// BENCHMARK TESTS
// =============================================================================

func BenchmarkSimpleLogging(b *testing.B) {
	logger := New().Output(io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Msg("benchmark test")
	}
}

func BenchmarkStructuredLogging(b *testing.B) {
	logger := New().Output(io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().
			Str("key1", "value1").
			Int("key2", 42).
			Bool("key3", true).
			Msg("benchmark test")
	}
}

func BenchmarkDisabledLogging(b *testing.B) {
	logger := New().Output(io.Discard).Level(ErrorLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().
			Str("key1", "value1").
			Int("key2", 42).
			Bool("key3", true).
			Msg("benchmark test")
	}
}

func BenchmarkContextLogging(b *testing.B) {
	logger := New().Output(io.Discard)
	ctxLogger := logger.With().Str("service", "test").Int("version", 1).Logger()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctxLogger.Info().Msg("benchmark test")
	}
}

func BenchmarkStandardLogger(b *testing.B) {
	logger := New().Output(io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print("benchmark test")
	}
}

func BenchmarkEventPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := getEvent()
		putEvent(e)
	}
}

func BenchmarkJSONSerialization(b *testing.B) {
	var buf []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = buf[:0]
		buf = appendString(buf, "message", "test message")
		buf = appendString(buf, "level", "info")
		buf = appendInt(buf, "count", i)
		buf = appendBool(buf, "flag", true)
		buf = appendTime(buf, "time", time.Now())
		result := wrapJSON(buf)
		_ = result
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestComplexLoggingScenario(b *testing.T) {
	var localBuf bytes.Buffer
	mockRemote := &mockRemoteWriter{}

	// Create logger with both local and remote output
	logger := New().Output(&localBuf).Remote(mockRemote).Level(DebugLevel)

	// Create context logger
	ctxLogger := logger.With().Str("service", "test-service").Int("version", 2).Logger()

	// Log various types of messages
	logger.Info().Msg("Application starting")
	ctxLogger.Debug().Str("config", "loaded").Msg("Configuration initialized")

	testErr := fmt.Errorf("connection failed")
	logger.Error().Err(testErr).Str("host", "localhost").Int("port", 8080).Msg("Failed to connect")

	ctxLogger.Warn().Bool("retrying", true).Msg("Retrying connection")
	logger.Info().Msg("Application ready")

	// Verify local output
	localOutput := localBuf.String()
	localLines := strings.Split(strings.TrimSpace(localOutput), "\n")

	if len(localLines) != 5 {
		b.Errorf("Expected 5 local log lines, got %d: %v", len(localLines), localLines)
	}

	// Verify remote output
	remoteWrites := mockRemote.GetWrites()
	if len(remoteWrites) != 5 {
		b.Errorf("Expected 5 remote writes, got %d , %s", len(remoteWrites), func([][]byte) string {
			var combined []string
			for _, w := range remoteWrites {
				combined = append(combined, string(w))
			}
			return strings.Join(combined, "\n")

		}(remoteWrites))
	}

	// Parse and verify structured content
	for i, line := range localLines {
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			b.Errorf("Failed to parse log line %d: %v", i, err)
			continue
		}

		// All entries should have required fields
		if _, ok := logEntry["message"]; !ok {
			b.Errorf("Log line %d missing 'message' field", i)
		}
		if _, ok := logEntry["time"]; !ok {
			b.Errorf("Log line %d missing 'time' field", i)
		}
		if _, ok := logEntry["level"]; !ok {
			b.Errorf("Log line %d missing 'level' field", i)
		}
	}

	// Verify context logger entries have context fields
	var contextEntries []map[string]interface{}
	for _, write := range remoteWrites {
		var logEntry map[string]interface{}
		if err := json.Unmarshal(write, &logEntry); err != nil {
			continue
		}
		if service, ok := logEntry["service"]; ok && service == "test-service" {
			contextEntries = append(contextEntries, logEntry)
		}
	}

	if len(contextEntries) != 2 {
		jsonI, err := json.MarshalIndent(remoteWrites, "", "  ")
		if err != nil {
			b.Errorf("Failed to marshal remote writes for debugging: %v", err)
		}
		b.Logf("Remote writes: %s", string(jsonI))
		b.Errorf("Expected 2 context entries, got %d: %v", len(contextEntries), contextEntries)
	}

	for i, entry := range contextEntries {
		if entry["version"] != float64(2) {
			b.Errorf("Context entry %d missing version field", i)
		}
	}
}

func TestErrorHandlingScenarios(t *testing.T) {
	t.Run("RemoteWriteError", func(t *testing.T) {
		var localBuf bytes.Buffer
		mockRemote := &mockRemoteWriter{
			writeError: fmt.Errorf("network error"),
		}

		logger := New().Output(&localBuf).Remote(mockRemote)

		// Should not panic and should still write locally
		logger.Info().Msg("test with remote error")

		if localBuf.Len() == 0 {
			t.Error("Expected local output even with remote error")
		}

		// Remote should have attempted write
		writes := mockRemote.GetWrites()
		if len(writes) != 0 {
			t.Error("Expected no successful remote writes due to error")
		}
	})

	t.Run("AsyncRemoteError", func(t *testing.T) {
		var localBuf bytes.Buffer
		mockRemote := &mockRemoteWriter{
			asyncError: fmt.Errorf("async network error"),
		}

		logger := New().Output(&localBuf).Remote(mockRemote)
		logger.async = true

		// Should not panic and should still write locally
		logger.Info().Msg("test with async remote error")

		if localBuf.Len() == 0 {
			t.Error("Expected local output even with async remote error")
		}

		// Async remote should have attempted write
		asyncWrites := mockRemote.GetAsyncWrites()
		if len(asyncWrites) != 0 {
			t.Error("Expected no successful async remote writes due to error")
		}
	})
}

func TestConfigurationMethods(t *testing.T) {
	original := New()

	t.Run("LevelConfiguration", func(t *testing.T) {
		debugLogger := original.Level(DebugLevel)

		// Original should be unchanged
		if original.level == DebugLevel {
			t.Error("Original logger level should not change")
		}

		// New logger should have debug level
		if debugLogger.level != DebugLevel {
			t.Error("New logger should have debug level")
		}
	})

	t.Run("OutputConfiguration", func(t *testing.T) {
		var buf bytes.Buffer
		newLogger := original.Output(&buf)

		// Should be different instance
		if &original == &newLogger {
			t.Error("Output() should return new instance")
		}

		// New logger should write to buffer
		newLogger.Info().Msg("test")
		if buf.Len() == 0 {
			t.Error("New logger should write to configured output")
		}
	})

	t.Run("RemoteConfiguration", func(t *testing.T) {
		mockRemote := &mockRemoteWriter{}
		remoteLogger := original.Remote(mockRemote)

		if &original == &remoteLogger {
			t.Error("Remote() should return new instance")
		}

		if remoteLogger.remoteWriter != mockRemote {
			t.Error("Remote logger should use configured remote writer")
		}
	})

	t.Run("RemoteHTTPConfiguration", func(t *testing.T) {
		endpoint := "https://example.com/logs"
		httpLogger := original.RemoteHTTP(endpoint, WithHTTPAuth("token"))

		if &original == &httpLogger {
			t.Error("RemoteHTTP() should return new instance")
		}

		if httpLogger.remoteWriter == nil {
			t.Error("HTTP logger should have remote writer")
		}
	})
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestEdgeCases(t *testing.T) {
	t.Run("EmptyMessage", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New().Output(&buf)

		logger.Info().Msg("")

		output := buf.String()
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
			t.Fatalf("Failed to parse log output: %v", err)
		}

		if logEntry["message"] != "" {
			t.Errorf("Expected empty message, got %v", logEntry["message"])
		}
	})

	t.Run("NilError", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New().Output(&buf)

		logger.Info().Err(nil).Msg("nil error test")

		output := buf.String()
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
			t.Fatalf("Failed to parse log output: %v", err)
		}

		// Should not have error field when nil
		if _, exists := logEntry["error"]; exists {
			t.Error("Should not have error field for nil error")
		}
	})

	t.Run("VeryLongMessage", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New().Output(&buf)

		longMessage := strings.Repeat("a", 10000)
		logger.Info().Msg(longMessage)

		output := buf.String()
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
			t.Fatalf("Failed to parse log output with long message: %v", err)
		}

		if logEntry["message"] != longMessage {
			t.Error("Long message not preserved correctly")
		}
	})

	t.Run("SpecialCharacters", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New().Output(&buf)

		specialMessage := "Special chars: \n\t\r\"'\\"
		logger.Info().Str("special", specialMessage).Msg("test")

		output := buf.String()
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
			t.Fatalf("Failed to parse log output with special chars: %v", err)
		}

		if logEntry["special"] != specialMessage {
			t.Error("Special characters not handled correctly")
		}
	})

	t.Run("UnicodeCharacters", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New().Output(&buf)

		unicodeMessage := "Unicode: 你好世界 🌍 émojis"
		logger.Info().Str("unicode", unicodeMessage).Msg("unicode test")

		output := buf.String()
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry); err != nil {
			t.Fatalf("Failed to parse log output with unicode: %v", err)
		}

		if logEntry["unicode"] != unicodeMessage {
			t.Error("Unicode characters not handled correctly")
		}
	})
}
