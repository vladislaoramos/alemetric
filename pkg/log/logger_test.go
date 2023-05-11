package logger

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("debug", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New("debug", &buf)

		logger.Debug("debug message")
		require.True(t, strings.HasSuffix(buf.String(), "debug message\"\n"))

		logger.Info("info message")
		require.True(t, strings.HasSuffix(buf.String(), "info message\"\n"))

		logger.Warn("warning message")
		require.True(t, strings.HasSuffix(buf.String(), "warning message\"\n"))

		logger.Error("error message")
		require.True(t, strings.HasSuffix(buf.String(), "error message\"\n"))
	})

	t.Run("info", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New("info", &buf)

		logger.Info("info message")
		require.True(t, strings.HasSuffix(buf.String(), "info message\"\n"))

		logger.Warn("warning message")
		require.True(t, strings.HasSuffix(buf.String(), "warning message\"\n"))

		logger.Error("error message")
		require.True(t, strings.HasSuffix(buf.String(), "error message\"\n"))
	})

	t.Run("warn", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New("warn", &buf)

		logger.Warn("warning message")
		require.True(t, strings.HasSuffix(buf.String(), "warning message\"\n"))

		logger.Error("error message")
		require.True(t, strings.HasSuffix(buf.String(), "error message\"\n"))
	})

	t.Run("error", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New("error", &buf)

		logger.Error("error message")
		require.True(t, strings.HasSuffix(buf.String(), "error message\"\n"))
	})

	t.Run("default", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New("default", &buf)

		logger.Info("info message")
		require.True(t, strings.HasSuffix(buf.String(), "info message\"\n"))
	})
}
