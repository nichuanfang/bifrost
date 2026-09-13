package mistral

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
	const body = `{"object":"list","data":[{"id":"mistral-small-latest","owned_by":"mistralai"}]}`
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

	provider := NewMistralProvider(&schemas.ProviderConfig{
		NetworkConfig: schemas.NetworkConfig{BaseURL: server.URL},
	}, &testLogger{})

	response, bifrostErr := provider.listModelsByKey(
		schemas.NewBifrostContext(context.Background(), schemas.NoDeadline),
		schemas.Key{Models: schemas.WhiteList{"*"}},
		&schemas.BifrostListModelsRequest{Provider: schemas.Mistral, Unfiltered: true},
	)

	require.Nil(t, bifrostErr)
	require.NotNil(t, response)
	require.Len(t, response.Data, 1)
	require.Equal(t, "mistral/mistral-small-latest", response.Data[0].ID)
}
