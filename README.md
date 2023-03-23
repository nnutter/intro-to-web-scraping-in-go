# Intro to Web Scraping in Go

Working our way up the tech stack:

1. [http](https://github.com/nnutter/intro-to-web-scraping-in-go/blob/master/internal/http): We have to be able to retrieve the content from the web.
2. [html](https://github.com/nnutter/intro-to-web-scraping-in-go/blob/master/internal/html): In the context of web scraping that content is usually HTML so it helps to parse it.
3. [cascadia](https://github.com/nnutter/intro-to-web-scraping-in-go/blob/master/internal/cascadia): Allows us to navigate the HTML using the common language of CSS selectors.
4. [goquery](https://github.com/nnutter/intro-to-web-scraping-in-go/blob/master/internal/goquery): Takes that down a somewhat familiar path for many by mimicking jQuery.
5. [colly](https://github.com/nnutter/intro-to-web-scraping-in-go/blob/master/internal/colly): Builds on `goquery` to provide a web scraping API.
6. [rod](https://github.com/nnutter/intro-to-web-scraping-in-go/blob/master/internal/rod): Provides an API to interact with the DevTools Protocol allowing you to drive a browser such as Google Chrome and all the features that come with that such as JavaScript interop.
