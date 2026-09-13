package gemini

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/stretchr/testify/require"
)

func TestListModelsByKeyDecodesGzipResponse(t *testing.T) {
	const body = `{"models":[{"name":"models/gemini-2.5-pro","displayName":"Gemini 2.5 Pro","supportedGenerationMethods":["generateContent"]}]}`
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	_, err := writer.Write([]byte(body))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(compressed.Bytes())
	}))
	defer server.Close()

	provider := NewGeminiProvider(&schemas.ProviderConfig{
		NetworkConfig: schemas.NetworkConfig{BaseURL: server.URL},
	}, testNoopLogger{})

	response, bifrostErr := provider.listModelsByKey(
		schemas.NewBifrostContext(context.Background(), schemas.NoDeadline),
		schemas.Key{Value: *schemas.NewSecretVar("dummy-key")},
		&schemas.BifrostListModelsRequest{Provider: schemas.Gemini, Unfiltered: true},
	)

	require.Nil(t, bifrostErr)
	require.NotNil(t, response)
	require.Len(t, response.Data, 1)
	require.Equal(t, "gemini/gemini-2.5-pro", response.Data[0].ID)
}
