#import <Cocoa/Cocoa.h>

extern void OnAction(const char* item);

@interface TrayMenu : NSObject {
  @private
    NSStatusItem *_statusItem;
}
@end

@implementation TrayMenu

- (void) openWebsite:(id)sender {
  NSURL *url = [NSURL URLWithString:@"https://getlantern.org"];
  [[NSWorkspace sharedWorkspace] openURL:url];
  [url release];
}

- (IBAction)menuHandler:(id)sender;
{
    const char *c = [[sender representedObject] UTF8String];
    OnAction(c);
}

- (NSMenu *) createMenu {
  NSZone *menuZone = [NSMenu menuZone];
  NSMenu *menu = [[NSMenu allocWithZone:menuZone] init];
  NSMenuItem *menuItem;

  // Add Quit Action
  menuItem = [menu addItemWithTitle:@"Quit"
                     action:@selector(menuHandler:)
                      keyEquivalent:@""];
  [menuItem setToolTip:@"Click to Quit this App"];
  [menuItem setRepresentedObject:@"quit"];
  [menuItem setTarget:self];

  return menu;
}

- (void) applicationDidFinishLaunching:(NSNotification *)notification {
  NSMenu *menu = [self createMenu];

  _statusItem = [[[NSStatusBar systemStatusBar]
  statusItemWithLength:NSSquareStatusItemLength] retain];
  [_statusItem setImage:[[NSImage alloc] initWithContentsOfFile:@"icon.png"]];
  [_statusItem setMenu:menu];
  [_statusItem setHighlightMode:YES];
  [_statusItem setToolTip:@"Test Tray"];
  [menu release];
}

@end

int StartApp(void) {
  NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
  [NSApplication sharedApplication];

  //HandleItem("Hi There from C");

  TrayMenu *menu = [[TrayMenu alloc] init];
  [NSApp setDelegate:menu];
  [NSApp run];

  [pool release];
  return EXIT_SUCCESS;
}