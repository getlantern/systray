#import <Cocoa/Cocoa.h>
#import "systray.h"

@interface MenuItem : NSObject
{
  @public
  NSString* name;
  NSString* title;
  NSString* tooltip;
  void(*callback)();
}
@end
@implementation MenuItem
@end

@interface AppDelegate: NSObject <NSApplicationDelegate>
- (IBAction)clicked:(id)sender;
- (void) updateTitle:(NSString*)title;
- (void) addMenu:(MenuItem*) item;
- (void) createMenu;
- (IBAction)menuHandler:(id)sender;
@property (assign) IBOutlet NSWindow *window;
@end

@implementation AppDelegate
{
    NSStatusItem *statusItem;
    NSMenuItem *doStuff;
    NSMenu *menu;
}

@synthesize window = _window;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
  [self createMenu];
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

- (IBAction)clicked:(id)sender {
    NSMutableDictionary *cmd = [NSMutableDictionary dictionaryWithObjectsAndKeys:@"clicked", @"action", nil];
    NSLog(@"clicked");
}

- (IBAction)menuHandler:(id)sender
{
  MenuItem* item = [sender representedObject];
  item->callback();
}

- (void)createMenu
{
  NSZone *menuZone = [NSMenu menuZone];
  self->menu = [[NSMenu allocWithZone:menuZone] init];
}

- (void) addMenu:(MenuItem*) item

{
  // Add DoStuff Action
  self->doStuff = [menu addItemWithTitle:item->title
                               action:@selector(menuHandler:)
                        keyEquivalent:@""];
  [self->doStuff setToolTip:item->tooltip];
  [self->doStuff setRepresentedObject: item];
  [self->doStuff setTarget:self];
}

- (void) updateTitle:(NSString*)title
{
    self->doStuff.title = title;
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

void addMenu(char* name, char* title, char* tooltip, void* callback) {
    MenuItem* item = [[MenuItem alloc] init];
    item->name = [[NSString alloc] initWithCString:name encoding:NSUTF8StringEncoding];
    item->title = [[NSString alloc] initWithCString:title encoding:NSUTF8StringEncoding];
    item->tooltip = [[NSString alloc] initWithCString:tooltip encoding:NSUTF8StringEncoding];
    item->callback = callMe;
    [(AppDelegate*)[NSApp delegate] performSelectorOnMainThread:@selector(addMenu:) withObject:(id)item waitUntilDone: YES];
}

void updateTitle(char* title) {
    NSString *titleString = [[NSString alloc] initWithCString:title encoding:NSUTF8StringEncoding];
    [(AppDelegate*)[NSApp delegate] performSelectorOnMainThread:@selector(updateTitle:) withObject:(id)titleString waitUntilDone: YES];
}
