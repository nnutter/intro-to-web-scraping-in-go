package cascadia

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/andybalholm/cascadia"
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

		dom, err := html.Parse(resp.Body)
		require.NoError(t, err)

		sel, err := cascadia.Parse("table")
		require.NoError(t, err)

		tables := cascadia.QueryAll(dom, sel)
		for _, table := range tables {
			logElement(t, table)
		}

		sel, err = cascadia.Parse("#contentarea > table.mytable")
		require.NoError(t, err)
		table := cascadia.Query(dom, sel)
		t.Log(sel.String())
		logElement(t, table)
	})
}

func logElement(t *testing.T, n *html.Node) {
	var s strings.Builder
	s.WriteString(n.Data + "[")
	for i, a := range n.Attr {
		if i > 0 {
			s.WriteString(" ")
		}
		s.WriteString(a.Key + `="` + a.Val + `"`)
	}
	s.WriteString("]")
	t.Log(s.String())
}
