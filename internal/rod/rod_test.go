package rod

import (
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

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

	t.Run("WebScrapingStack", func(t *testing.T) {
		tables := page.MustElements("table")
		for _, tableElement := range tables {
			t.Log(tableElement.String())
		}

		sel := "#contentarea > table.mytable"
		myTable := page.MustElement(sel)
		t.Log(sel)
		t.Log(myTable.String())
	})

	t.Run("ExtractScoreboard", func(t *testing.T) {
		scoreboardTable := page.MustElement("#contentarea > table.mytable")
		scoreboard, err := extractScoreboard(scoreboardTable)
		require.NoError(t, err)
		scoreboard.SetOutputMirror(os.Stdout)
		scoreboard.Render()
	})

	t.Run("Helpful Practices for Rod", func(t *testing.T) {
		t.Run("Element vs Elements", func(t *testing.T) {
			t.Run("Elements Returns Instantly", func(t *testing.T) {
				elements := page.MustElements("#does-not-exist")
				assert.Len(t, elements, 0)
			})
			t.Run("Element Waits for First", func(t *testing.T) {
				_ = page.MustElement("#does-not-exist")
			})
		})

		// https://go-rod.github.io/#/context-and-timeout?id=cancellation
		t.Run("Add a Deadline", func(t *testing.T) {
			t.Run("Using Go Context", func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				_ = page.Context(ctx).MustElement("#does-not-exist")
			})
			t.Run("Using Helper Function", func(t *testing.T) {
				_ = page.Timeout(3 * time.Second).MustElement("#does-not-exist")
			})
			t.Run("Clear Deadline", func(t *testing.T) {
				scoreboardTable := page.Timeout(2 * time.Second).MustElement("#contentarea > table.mytable")
				time.Sleep(2 * time.Second)
				assert.NotEmpty(t, scoreboardTable.CancelTimeout().MustText())
			})
		})

		// https://go-rod.github.io/#/selectors/README?id=by-text-content
		t.Run("Selector + Regex", func(t *testing.T) {
			boxScoreAnchor := page.MustElementR("a", `Box Score`)
			assert.Equal(t, "/game/box_score/5414758", *boxScoreAnchor.MustAttribute("href"))
		})

		// https://go-rod.github.io/#/events/README?id=get-the-event-details
		t.Run("Check HTTP Status Code", func(t *testing.T) {
			page, err := browser.Page(proto.TargetCreateTarget{})
			require.NoError(t, err)

			notFoundURL := u
			notFoundURL.Path = "/not-found.html"

			var e proto.NetworkResponseReceived
			wait := page.WaitEvent(&e)
			err = page.Navigate(notFoundURL.String())
			require.NoError(t, err)
			wait()
			assert.Equal(t, http.StatusNotFound, e.Response.Status)
		})

		// https://go-rod.github.io/#/selectors/README?id=race-selectors
		t.Run("Race Selectors", func(t *testing.T) {
			page := browser.MustPage("")

			page.Race().
				Element("div#roster").
				MustHandle(func(e *rod.Element) {
					// Handle roster page.
				}).
				Element("div#not-found").
				MustHandle(func(e *rod.Element) {
					// Handle not found error page.
				}).
				MustDo()
		})

		// https://go-rod.github.io/#/network?id=hijack-requests
		t.Run("Block Extraneous Resources", func(t *testing.T) {
			page, err := browser.Page(proto.TargetCreateTarget{})
			require.NoError(t, err)

			router := page.HijackRequests()
			router.MustAdd("*.png", func(ctx *rod.Hijack) {
				if ctx.Request.Type() == proto.NetworkResourceTypeImage {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
					return
				}
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			})
			router.MustAdd("*.css", func(ctx *rod.Hijack) {
				if ctx.Request.Type() == proto.NetworkResourceTypeStylesheet {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
					return
				}
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			})
			router.MustAdd("*.otf", func(ctx *rod.Hijack) {
				if ctx.Request.Type() == proto.NetworkResourceTypeFont {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
					return
				}
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			})
			router.MustAdd("*.js", func(ctx *rod.Hijack) {
				if ctx.Request.URL().Host != "stats.ncaa.org" && ctx.Request.Type() == proto.NetworkResourceTypeScript {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
					return
				}
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			})
			go router.Run()

			_ = page.MustNavigate(boxScoreURL.String())
		})

		// https://go-rod.github.io/#/browsers-pages?id=page-pool
		t.Run("Limit Memory Using a Page Pool", func(t *testing.T) {
			browser := rod.New().MustConnect()
			defer browser.MustClose()

			pool := rod.NewPagePool(2)
			defer func() {
				pool.Cleanup(func(p *rod.Page) { p.MustClose() })
			}()

			create := func() *rod.Page {
				return browser.MustPage()
			}
			wg, ctx := errgroup.WithContext(context.Background())
			for range "...." {
				wg.Go(func() error {
					page := pool.Get(create)
					defer func() {
						t.Log("returning page")
						pool.Put(page)
					}()

					t.Log("loading page")
					page.Context(ctx).MustNavigate(boxScoreURL.String()).MustWaitLoad()
					time.Sleep(2 * time.Second)

					return nil
				})
			}
			require.NoError(t, wg.Wait())
		})
	})
}

func extractScoreboard(t *rod.Element) (table.Writer, error) {
	trs := t.MustElements("tr")
	headerRow := trs[0]
	visitorRow := trs[1]
	homeRow := trs[2]

	tw := table.NewWriter()

	tdToRow := func(tr *rod.Element) table.Row {
		var row []interface{}
		for _, td := range tr.MustElements("td") {
			row = append(row, td.MustText())
		}
		return row
	}

	tw.AppendHeader(tdToRow(headerRow))
	tw.AppendRow(tdToRow(visitorRow))
	tw.AppendRow(tdToRow(homeRow))

	// Structs, tags, reflection.

	return tw, nil
}
