package restapitest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPing(t *testing.T) {
	t.Parallel()
	resp, err := gClient.GetPing(gCtx)
	defer closeBody(resp)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, strings.HasPrefix(resp.Header.Get("Content-Type"), "text/plain"))
	assert.Equal(t, "pong", readAll(t, resp.Body))
}
