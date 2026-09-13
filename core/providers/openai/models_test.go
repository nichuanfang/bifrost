package openai

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestListModelsByKeyResponseShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		body        string
		wantID      string
		wantOwnedBy string
		wantContext int
	}{
		{
			name:        "OpenAI envelope",
			body:        `{"object":"list","data":[{"id":"gpt-5","owned_by":"openai","context_window":128000}]}`,
			wantID:      "test/gpt-5",
			wantOwnedBy: "openai",
			wantContext: 128000,
		},
		{
			name:        "Together array",
			body:        `[{"id":"zai-org/GLM-5.2","organization":"Z.ai","context_length":131072}]`,
			wantID:      "test/zai-org/GLM-5.2",
			wantOwnedBy: "Z.ai",
			wantContext: 131072,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()

			response, bifrostErr := ListModelsByKey(
				schemas.NewBifrostContext(context.Background(), schemas.NoDeadline),
				&fasthttp.Client{},
				server.URL,
				schemas.Key{Models: schemas.WhiteList{"*"}},
				false,
				nil,
				schemas.ModelProvider("test"),
				false,
				false,
			)

			require.Nil(t, bifrostErr)
			require.Len(t, response.Data, 1)
			require.Equal(t, test.wantID, response.Data[0].ID)
			require.Equal(t, schemas.Ptr(test.wantOwnedBy), response.Data[0].OwnedBy)
			require.Equal(t, schemas.Ptr(test.wantContext), response.Data[0].ContextLength)
		})
	}
}

func TestListModelsByKeyDecodesGzipResponse(t *testing.T) {
	t.Parallel()

	const body = `{"object":"list","data":[{"id":"gpt-5","owned_by":"openai"}]}`
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	_, err := writer.Write([]byte(body))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(compressed.Bytes())
	}))
	defer server.Close()

	response, bifrostErr := ListModelsByKey(
		schemas.NewBifrostContext(context.Background(), schemas.NoDeadline),
		&fasthttp.Client{},
		server.URL,
		schemas.Key{Models: schemas.WhiteList{"*"}},
		false,
		nil,
		schemas.ModelProvider("test"),
		false,
		false,
	)

	require.Nil(t, bifrostErr)
	require.NotNil(t, response)
	require.Len(t, response.Data, 1)
	require.Equal(t, "test/gpt-5", response.Data[0].ID)
}

func TestListModelsByKeyReportsResponseDecodeErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		encoding string
		body     string
		wantErr  string
	}{
		{name: "damaged gzip", encoding: "gzip", body: `{"object":"list"}`, wantErr: schemas.ErrProviderResponseDecode},
		{name: "unknown encoding", encoding: "compress", body: `{"object":"list"}`, wantErr: schemas.ErrProviderResponseDecode},
		{name: "ordinary JSON", body: `{"object":"list","data":[{"id":"gpt-5"}]}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if test.encoding != "" {
					w.Header().Set("Content-Encoding", test.encoding)
				}
				_, _ = w.Write([]byte(test.body))
			}))
			defer server.Close()

			response, bifrostErr := ListModelsByKey(
				schemas.NewBifrostContext(context.Background(), schemas.NoDeadline),
				&fasthttp.Client{},
				server.URL,
				schemas.Key{Models: schemas.WhiteList{"*"}},
				false,
				nil,
				schemas.ModelProvider("test"),
				false,
				false,
			)

			if test.wantErr == "" {
				require.Nil(t, bifrostErr)
				require.NotNil(t, response)
				require.Len(t, response.Data, 1)
				return
			}

			require.Nil(t, response)
			require.NotNil(t, bifrostErr)
			require.NotNil(t, bifrostErr.Error)
			require.Equal(t, test.wantErr, bifrostErr.Error.Message)
		})
	}
}
