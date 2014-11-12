#import <Cocoa/Cocoa.h>

extern void OnAction(const char* item);

@interface TrayMenu : NSObject {
  @private
    NSStatusItem *_statusItem;
    NSMenuItem *doStuff;
}

- (void) updateTitle:(NSString*)title;

@end

@implementation TrayMenu

- (IBAction)menuHandler:(id)sender;
{
    const char *c = [[sender representedObject] UTF8String];
    OnAction(c);
}

- (NSMenu *) createMenu
{
  NSZone *menuZone = [NSMenu menuZone];
  NSMenu *menu = [[NSMenu allocWithZone:menuZone] init];

  // Add DoStuff Action
  self->doStuff = [menu addItemWithTitle:@"Change Me"
                               action:@selector(menuHandler:)
                        keyEquivalent:@""];
  [self->doStuff setToolTip:@"Click to change the title"];
  [self->doStuff setRepresentedObject:@"dostuff"];
  [self->doStuff setTarget:self];

  // Add Quit Action
  NSMenuItem *quit = [menu addItemWithTitle:@"Quit"
                     action:@selector(menuHandler:)
                      keyEquivalent:@""];
  [quit setToolTip:@"Click to Quit this App"];
  [quit setRepresentedObject:@"quit"];
  [quit setTarget:self];

  return menu;
}

- (void) updateTitle:(NSString*)title
{
    self->doStuff.title = title;
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

TrayMenu *menu;

int StartApp(void) {
  NSAutoreleasePool *pool = [[NSAutoreleasePool alloc] init];
  [NSApplication sharedApplication];

  //HandleItem("Hi There from C");

  menu = [[TrayMenu alloc] init];
  [NSApp setDelegate:menu];
  [NSApp run];

  [pool release];
  return EXIT_SUCCESS;
}

void updateTitle(char* title) {
    NSString *titleString = [[NSString alloc] initWithCString:title encoding:NSUTF8StringEncoding];
    [menu performSelectorOnMainThread:@selector(updateTitle:) withObject:(id)titleString waitUntilDone: YES];
}