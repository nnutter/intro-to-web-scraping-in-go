package http_test

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/intro-to-web-scraping-in-go/internal/demo"
)

func Test(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:8080",
	}
	demo.StartServer(t, &u)

	boxScoreURL := u
	boxScoreURL.Path = "/box_score.html"

	t.Run("WebScrapingStack", func(t *testing.T) {
		resp, err := http.Get(boxScoreURL.String())
		require.NoError(t, err)
		t.Cleanup(func() {
			_ = resp.Body.Close()
		})

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		t.Log(string(body))
	})
}
