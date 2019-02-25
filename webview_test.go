package webview_test

import (
	"github.com/KnownUnown/webview"
	"net/url"
)

func ExampleWebView() {
	v := webview.New()

	u, _ := url.Parse("https://google.com")
	v.Navigate(u)

	v.Loop()
}
