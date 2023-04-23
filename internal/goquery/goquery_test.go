package goquery

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
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

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		require.NoError(t, err)

		tables := doc.Find("table")
		tables.Each(func(i int, selection *goquery.Selection) {
			logElement(t, selection.Nodes[0])
		})

		sel := "#contentarea > table.mytable"
		table := doc.Find(sel).First()
		t.Log(sel)
		logElement(t, table.Nodes[0])
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
