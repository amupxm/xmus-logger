package xmuslogger

type HTTPRemoteWriter struct {
	endpoint string
	headers  map[string]string
	// Implementation details would go here
}

type HTTPOption func(*HTTPRemoteWriter)

func NewHTTPRemoteWriter(endpoint string, options ...HTTPOption) *HTTPRemoteWriter {
	w := &HTTPRemoteWriter{
		endpoint: endpoint,
		headers:  make(map[string]string),
	}
	for _, opt := range options {
		opt(w)
	}
	return w
}

func (h *HTTPRemoteWriter) Write(data []byte) error      { return nil } // TODO: Implement
func (h *HTTPRemoteWriter) WriteAsync(data []byte) error { return nil } // TODO: Implement
func (h *HTTPRemoteWriter) Flush() error                 { return nil }
func (h *HTTPRemoteWriter) Close() error                 { return nil }

func WithHTTPAuth(token string) HTTPOption {
	return func(w *HTTPRemoteWriter) {
		w.headers["Authorization"] = "Bearer " + token
	}
}

func WithHTTPHeaders(headers map[string]string) HTTPOption {
	return func(w *HTTPRemoteWriter) {
		for k, v := range headers {
			w.headers[k] = v
		}
	}
}
