package api

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWriter_Hijack_Supported(t *testing.T) {
	// httptest.ResponseRecorder does not implement Hijacker by default,
	// so we wrap it with a mock that does.
	rw := &responseWriter{
		ResponseWriter: &hijackableRecorder{ResponseRecorder: httptest.NewRecorder()},
		statusCode:     http.StatusOK,
	}

	// Verify that responseWriter satisfies http.Hijacker
	var _ http.Hijacker = rw

	conn, brw, err := rw.Hijack()
	require.NoError(t, err)
	assert.Nil(t, conn) // our mock returns nil
	assert.Nil(t, brw)
}

func TestResponseWriter_Hijack_NotSupported(t *testing.T) {
	rw := &responseWriter{
		ResponseWriter: httptest.NewRecorder(), // does NOT implement Hijacker
		statusCode:     http.StatusOK,
	}

	_, _, err := rw.Hijack()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not implement http.Hijacker")
}

// hijackableRecorder wraps httptest.ResponseRecorder and implements http.Hijacker.
type hijackableRecorder struct {
	*httptest.ResponseRecorder
}

func (h *hijackableRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, nil
}
