#import <Cocoa/Cocoa.h>
#import "systray.h"

@interface MenuItem : NSObject
{
  @public
  NSString* name;
  NSString* title;
  NSString* tooltip;
}
@end
@implementation MenuItem
@end

@interface AppDelegate: NSObject <NSApplicationDelegate>
- (IBAction)clicked:(id)sender;
- (void) addMenu:(MenuItem*) item;
- (void) createMenu;
- (IBAction)menuHandler:(id)sender;
@property (assign) IBOutlet NSWindow *window;
@end

@implementation AppDelegate
{
    NSStatusItem *statusItem;
    NSMenu *menu;
}

@synthesize window = _window;

- (id)init
{
  [self createMenu];
  return self;
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
  statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
  [statusItem setAction:@selector(clicked:)];
  NSImage *image = [[NSImage alloc] initWithContentsOfFile:@"icon.png"];
  if (image == nil) {
    NSLog(@"load icon error");
  }
  [statusItem setImage:image];
  [statusItem setMenu:menu];
  [statusItem setToolTip:@"Test Tray"];
  [menu autorelease];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification
{
  [[NSStatusBar systemStatusBar] removeStatusItem: statusItem];
}

- (IBAction)clicked:(id)sender {
    NSMutableDictionary *cmd = [NSMutableDictionary dictionaryWithObjectsAndKeys:@"clicked", @"action", nil];
    NSLog(@"clicked");
}

- (IBAction)menuHandler:(id)sender
{
  MenuItem* item = [sender representedObject];
  callMe((char*)[item->name cStringUsingEncoding: NSUTF8StringEncoding]);
}

- (void)createMenu
{
  NSZone *menuZone = [NSMenu menuZone];
  self->menu = [[NSMenu allocWithZone:menuZone] init];
}

- (void) addMenu:(MenuItem*) item
{
  NSMenuItem* menuItem = [menu addItemWithTitle:item->title
                               action:@selector(menuHandler:)
                        keyEquivalent:@""];
  [menuItem setToolTip:item->tooltip];
  [menuItem setRepresentedObject: item];
  [menuItem setTarget:self];
}

- (void) quit
{
  [NSApp terminate:self];
}

@end

int nativeLoop(void) {
  [NSAutoreleasePool new];
  AppDelegate *delegate = [[[AppDelegate alloc] init] autorelease];
  [[NSApplication sharedApplication] setDelegate:delegate];
  [NSApp run];
  NSLog(@"Quiting...");
  return EXIT_SUCCESS;
}

void addMenu(char* name, char* title, char* tooltip) {
    MenuItem* item = [[MenuItem alloc] init];
    item->name = [[NSString alloc] initWithCString:name encoding:NSUTF8StringEncoding];
    item->title = [[NSString alloc] initWithCString:title encoding:NSUTF8StringEncoding];
    item->tooltip = [[NSString alloc] initWithCString:tooltip encoding:NSUTF8StringEncoding];
    [(AppDelegate*)[NSApp delegate] performSelectorOnMainThread:@selector(addMenu:) withObject:(id)item waitUntilDone: YES];
}

void quit() {
    [(AppDelegate*)[NSApp delegate] performSelectorOnMainThread:@selector(quit) withObject:nil waitUntilDone: YES];
}
