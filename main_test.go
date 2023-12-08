package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"gotest.tools/v3/assert"
)

func TestServersStarted(t *testing.T) {
	serversCount := 10
	for i := 0; i < serversCount; i++ {
		startHttpServer(i, fmt.Sprintf(":807%d", i), fmt.Sprintf("hello from %d", i))
	}
	for i := 0; i < serversCount; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:807%d", i))
		assert.NilError(t, err)
		defer resp.Body.Close()
		bytes, err := io.ReadAll(resp.Body)
		assert.NilError(t, err)
		assert.Equal(t, string(bytes), fmt.Sprintf("hello from %d", i))
	}
}
