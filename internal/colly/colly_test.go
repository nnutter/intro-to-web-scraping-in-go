package colly

import (
	"net/url"
	"strings"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

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

	c := colly.NewCollector()

	c.OnHTML("table", func(table *colly.HTMLElement) {
		logElement(t, table.DOM.Nodes[0])
	})

	sel := "#contentarea > table.mytable"
	c.OnHTML(sel, func(table *colly.HTMLElement) {
		t.Log(sel)
		logElement(t, table.DOM.Nodes[0])
	})

	err := c.Visit(boxScoreURL.String())
	require.NoError(t, err)
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
