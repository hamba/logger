package logger_test

import (
	"bytes"
	"testing"

	"github.com/hamba/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncWriter(t *testing.T) {
	var buf bytes.Buffer

	w := logger.NewSyncWriter(&buf)

	n, err := w.Write([]byte("test"))

	require.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, "test", buf.String())
}
