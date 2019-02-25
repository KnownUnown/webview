#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#import "_cgo_export.h"
#import "darwin.h"

NSAutoreleasePool* appPool = nil;

@implementation InteropScriptMessageHandler
- (void)userContentController:(WKUserContentController *)userContentController
      didReceiveScriptMessage:(WKScriptMessage *)message {
    scriptMessageCallback((char *) [message.body description].UTF8String);
}
@end

webview webview_create() {
    NSWindow* window;
    NSRect frame;

    appPool = [[NSAutoreleasePool alloc] init];
    frame = NSMakeRect(0, 0, 480, 640);
    window = [[NSWindow alloc] initWithContentRect:frame
                                      styleMask:NSWindowStyleMaskTitled
                                               |NSWindowStyleMaskClosable
                                               |NSWindowStyleMaskMiniaturizable
                                      backing:NSBackingStoreBuffered
                                      defer:NO];
    [window setTitle:@"WebView"];
    [window center];

    WKWebViewConfiguration* config;
    id<WKScriptMessageHandler> messageHandler;
    WKWebView* view;

    config = [[[WKWebViewConfiguration alloc] init] autorelease];
    messageHandler = [[[InteropScriptMessageHandler alloc] init] autorelease];
    [config.userContentController addScriptMessageHandler:messageHandler name:@"webview"];
    [config.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
    view = [[[WKWebView alloc] initWithFrame:frame configuration:config] autorelease];

    [window setContentView:view];
    [window makeKeyAndOrderFront:nil];

    webview v = malloc(sizeof(webview));
    v->view = view;

    return v;
}

void webview_load(webview v, const char* url) {
    NSString* str = [NSString stringWithUTF8String:url];
    NSURL* nsurl = [NSURL URLWithString:str];
    NSURLRequest* req = [NSURLRequest requestWithURL:nsurl];

    [v->view loadRequest:req];
}

void webview_add_user_script(webview v, const char* source) {
    NSString* str = [NSString stringWithUTF8String:source];
    WKUserScript* script = [[WKUserScript alloc] initWithSource:str
        injectionTime:WKUserScriptInjectionTimeAtDocumentStart forMainFrameOnly:YES];

    [v->view.configuration.userContentController addUserScript:script];
}

void webview_add_user_style(webview v, const char* source) {
}

void webview_eval(webview v, const char* source) {
    NSString *str = [NSString stringWithUTF8String:source];

    [v->view evaluateJavaScript:str completionHandler:nil];
}

void webview_close(webview v) {
    [v->view.window close];
}

void _webview_dispatch_cb(_webview_dispatch_args args) {
    dispatchCallback(args->id);
    free(args);
}

void _webview_dispatch(int id) {
    _webview_dispatch_args args = malloc(sizeof(_webview_dispatch_args));
    args->id = id;

    dispatch_async_f(dispatch_get_main_queue(), (void *) args, (dispatch_function_t) _webview_dispatch_cb);
}


void webview_loop() {
    [NSApp run];

    [NSApp release];
    [appPool release];
}
