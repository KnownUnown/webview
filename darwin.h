#import <WebKit/WebKit.h>
#import <Cocoa/Cocoa.h>

extern NSAutoreleasePool* appPool;

typedef struct webview {
    WKWebView* view;
} *webview;

typedef struct _webview_dispatch_args {
    int id;
} *_webview_dispatch_args;

typedef struct webview_config {
    const char* windowName;
    const char* userAgent;
} *webview_config;

@interface InteropScriptMessageHandler : NSObject <WKScriptMessageHandler>
@end

webview webview_create();

void webview_load(webview v, const char* url);
void webview_add_user_script(webview v, const char* source);
void webview_add_user_style(webview v, const char* source);
void webview_eval(webview v, const char* source);
void webview_close(webview v);

void _webview_dispatch(int id);
void webview_loop();
