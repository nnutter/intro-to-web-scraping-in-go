package http_test

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/nnutter/intro-to-web-scraping-in-go/internal/demo"
)

func Test(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:8080",
	}
	startDemoServer(t, &u)

	boxScoreURL := u
	boxScoreURL.Path = "/box_score.html"
	resp, err := http.Get(boxScoreURL.String())
	if err != nil {
		t.Fatalf("failed to get box score")
	}
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read box score")
	}

	t.Log(string(bs))
}

func startDemoServer(t *testing.T, u *url.URL) {
	server, err := demo.ListenAndServe(u.Host)
	if err != nil {
		t.Fatalf("failed to start demo server")
	}
	t.Cleanup(func() {
		_ = server.Close()
	})
}