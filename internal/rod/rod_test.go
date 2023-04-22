package rod

import (
	"net/url"
	"testing"

	"github.com/go-rod/rod"

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

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(boxScoreURL.String())
	defer page.MustClose()

	tables := page.MustElements("table")
	for _, table := range tables {
		t.Log(table.String())
	}

	sel := "#contentarea > table:nth-child(5)"
	table := page.MustElement(sel)
	t.Log(sel)
	t.Log(table.String())

	sel = "#contentarea > table.mytable"
	table = page.MustElement(sel)
	t.Log(sel)
	t.Log(table.String())
}
