package html

import (
	"net/http"
	"net/url"
	"testing"

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
	if err != nil {
		t.Fatalf("failed to get box score page")
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})
	doc, err := html.Parse(resp.Body)
	if err != nil {
		t.Fatalf("failed to parse box score page")
	}
	t.Log(doc)
}
