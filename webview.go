package webview

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "darwin.h"
 */
import "C"
import (
	"net/url"
	"runtime"
	"sync"
	"unsafe"
)

type DispatchCallback func()
type ScriptCallback func(string)
type WebView struct {
	ref C.webview
}

type callbackTracker struct {
	fnIdx int
	fnMu sync.Mutex
	fns map[int]DispatchCallback
}

var callbacks = callbackTracker{fns: make(map[int]DispatchCallback)}
var scriptCallback ScriptCallback


// New creates a new WebView.
func New() *WebView {
	v := C.webview_create()
	return &WebView{ref: v}
}

// Navigate loads the provided URL in the WebView.
func (v *WebView) Navigate(u *url.URL) {
	cs := C.CString(u.String())
	defer C.free(unsafe.Pointer(cs))

	C.webview_load(v.ref, cs)
}

// AddUserScript registers the given string as a user
// script, that is, a script that is persistent
// across many Navigate calls or page loads.
func (v *WebView) AddUserScript(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.webview_add_user_script(v.ref, cs)
}

// Eval evaluates a JavaScript string.
func (v *WebView) Eval(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))

	C.webview_eval(v.ref, cs)
}

// Close closes the WebView.
func (v *WebView) Close() {
	C.webview_close(v.ref)
}

// Dispatch runs the given function on the main thread
// at some point in the future.
func (v *WebView) Dispatch(cb DispatchCallback) {
	callbacks.fnMu.Lock()
	defer callbacks.fnMu.Unlock()

	callbacks.fns[callbacks.fnIdx] = cb

	C._webview_dispatch(C.int(callbacks.fnIdx))
	callbacks.fnIdx++
}

// SetScriptCallback registers a function to be called
// whenever JavaScript calls `window.webkit.messageHandlers.webview.postMessage`.
func (v *WebView) SetScriptCallback(cb ScriptCallback) {
	scriptCallback = cb
}

// Loop pumps the NSApp run loop, namely, it calls
// `[NSApp run]`. This function must be called from
// the main thread to run the WebView after you create it.
// It blocks until the app exits.
func (*WebView) Loop() {
	C.webview_loop()
}

//export dispatchCallback
func dispatchCallback(i C.int) {
	id := int(i)

	callbacks.fnMu.Lock()
	defer callbacks.fnMu.Unlock()

	callbacks.fns[id]()
	callbacks.fns[id] = nil
}

//export scriptMessageCallback
func scriptMessageCallback(data *C.char) {
	d := C.GoString(data)
	if scriptCallback != nil {
		scriptCallback(d)
	}
}

func init() {
	runtime.LockOSThread()
}
