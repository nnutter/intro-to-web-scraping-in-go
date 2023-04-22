package html

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

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
	resp, err := http.Get(boxScoreURL.String())
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	doc, err := html.Parse(resp.Body)
	require.NoError(t, err)

	tables := findNodes(doc, "table")
	for _, table := range tables {
		logNode(t, table)
	}
}

func findNodes(node *html.Node, d string) []*html.Node {
	if node == nil {
		return nil
	}
	children := findNodes(node.FirstChild, d)
	siblings := findNodes(node.NextSibling, d)
	var nodes []*html.Node
	if node.Data == d {
		nodes = append(nodes, node)
	}
	if len(children) > 0 {
		nodes = append(nodes, children...)
	}
	if len(siblings) > 0 {
		nodes = append(nodes, siblings...)
	}
	return nodes
}

func logNode(t *testing.T, n *html.Node) {
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
