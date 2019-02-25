webview
=======

[![GoDoc](https://godoc.org/github.com/KnownUnown/webview?status.svg)](https://godoc.org/github.com/KnownUnown/webview)

This package provides minimal bindings for WKWebView on OSX. The API is loosely 
based off of [zserge/webview](https://github.com/zserge/webview); however, this 
package only supports one platform, OSX.

All of the functions in this package must be called from the main thread. To 
achieve that, you can use `Dispatch`, which schedules a function to be called back 
sometime in the future.

Remember to call `Loop` after you initialize the `WebView`.