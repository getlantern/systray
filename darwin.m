#import <Cocoa/Cocoa.h>
#import "systray.h"

@interface MenuItem : NSObject
{
  @public
    NSString* menuId;
    NSString* title;
    NSString* tooltip;
}
@end
@implementation MenuItem
@end

@interface AppDelegate: NSObject <NSApplicationDelegate>
   - (void) addMenu:(MenuItem*) item;
   - (IBAction)menuHandler:(id)sender;
   @property (assign) IBOutlet NSWindow *window;
@end

@implementation AppDelegate
{
  NSStatusItem *statusItem;
  NSMenu *menu;
  NSCondition* cond;
}

@synthesize window = _window;

- (id) init {
  return self;
}

- (void) waitForFinishLaunching {
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
  self->statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
  NSZone *menuZone = [NSMenu menuZone];
  self->menu = [[NSMenu allocWithZone:menuZone] init];
  [self->statusItem setMenu:self->menu];
  [menu autorelease];
  systray_ready();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification
{
  [[NSStatusBar systemStatusBar] removeStatusItem: statusItem];
}

- (void)setIcon:(NSImage *)image {
  [statusItem setImage:image];
}

- (void)setTitle:(NSString *)title {
  [statusItem setTitle:title];
}

- (void)setTooltip:(NSString *)tooltip {
  [statusItem setToolTip:tooltip];
}

- (IBAction)menuHandler:(id)sender
{
  NSString* menuId = [sender representedObject];
  systray_menu_item_selected((char*)[menuId cStringUsingEncoding: NSUTF8StringEncoding]);
}

- (void) addMenu:(MenuItem*) item
{
  NSMenuItem* menuItem;
  int existedMenuIndex = [menu indexOfItemWithRepresentedObject: item->menuId];
  if (existedMenuIndex == -1) {
    menuItem = [menu addItemWithTitle:item->title action:@selector(menuHandler:) keyEquivalent:@""];
    [menuItem setTarget:self];
    [menuItem setRepresentedObject: item->menuId];
  }
  else {
    menuItem = [menu itemAtIndex: existedMenuIndex];
    [menuItem setTitle:item->title];
  }
  [menuItem setToolTip:item->tooltip];
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
  return EXIT_SUCCESS;
}

void runInMainThread(SEL method, id object) {
  [(AppDelegate*)[NSApp delegate] waitForFinishLaunching];
  [(AppDelegate*)[NSApp delegate]
    performSelectorOnMainThread:method
                     withObject:object
                  waitUntilDone: YES];
}

void setIcon(const char* iconBytes, int length) {
  NSData* buffer = [NSData dataWithBytes: iconBytes length:length];
  NSImage *image = [[NSImage alloc] initWithData:buffer];
  runInMainThread(@selector(setIcon:), (id)image);
}

void setTitle(const char* ctitle) {
  NSString* title = [[NSString alloc] initWithCString:ctitle encoding:NSUTF8StringEncoding];
  runInMainThread(@selector(setTitle:), (id)title);
}

void setTooltip(const char* ctooltip) {
  NSString* tooltip = [[NSString alloc] initWithCString:ctooltip encoding:NSUTF8StringEncoding];
  runInMainThread(@selector(setTooltip:), (id)tooltip);
}

void addMenu(char* menuId, char* title, char* tooltip) {
  MenuItem* item = [[MenuItem alloc] init];
  item->menuId = [[NSString alloc] initWithCString:menuId encoding:NSUTF8StringEncoding];
  item->title = [[NSString alloc] initWithCString:title encoding:NSUTF8StringEncoding];
  item->tooltip = [[NSString alloc] initWithCString:tooltip encoding:NSUTF8StringEncoding];
  runInMainThread(@selector(addMenu:), (id)item);
}

void quit() {
  runInMainThread(@selector(quit), nil);
}
