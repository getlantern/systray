#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#include "systray.h"

NSWindowController *windowController = nil;
WKWebView *webView = nil;

void configureAppWindow(char* title, int width, int height)
{
  if (windowController != nil) {
    // already configured, ignore
    return;
  }

  NSRect frame = NSMakeRect(0, 0, width, height);
  int mask = NSWindowStyleMaskTitled | NSWindowStyleMaskResizable | NSWindowStyleMaskClosable;
  id window =
    [[NSWindow alloc] initWithContentRect:frame
                                 styleMask:mask
                                   backing:NSBackingStoreBuffered
                                     defer:NO];
  [window setTitle:[[NSString alloc] initWithUTF8String:title]];
  [window makeKeyAndOrderFront:nil];
  [window center];

  NSView *contentView = [window contentView];
  webView = [[WKWebView alloc] initWithFrame:[contentView bounds]];
  [webView setTranslatesAutoresizingMaskIntoConstraints:NO];
  [contentView addSubview:webView];
  [contentView addConstraint:
    [NSLayoutConstraint constraintWithItem:webView
        attribute:NSLayoutAttributeWidth
        relatedBy:NSLayoutRelationEqual
        toItem:contentView
        attribute:NSLayoutAttributeWidth
        multiplier:1
        constant:0]];
  [contentView addConstraint:
    [NSLayoutConstraint constraintWithItem:webView
        attribute:NSLayoutAttributeHeight
        relatedBy:NSLayoutRelationEqual
        toItem:contentView
        attribute:NSLayoutAttributeHeight
        multiplier:1
        constant:0]];

  // Window controller:
  windowController = [[NSWindowController alloc] initWithWindow:window];
  free(title);
}

void doShowAppWindow(char* url)
{
  if (windowController == nil) {
    // no app window to open
    return;
  }
  id nsURL = [NSURL URLWithString:[[NSString alloc] initWithUTF8String:url]];
  id req = [[NSURLRequest alloc] initWithURL: nsURL
                                 cachePolicy: NSURLRequestUseProtocolCachePolicy
                                 timeoutInterval: 5];
  [webView loadRequest:req];
  [windowController.window orderFrontRegardless];
  free(url);
}

void showAppWindow(char* url)
{
  dispatch_async(dispatch_get_main_queue(), ^{
    doShowAppWindow(url);
  });
}
