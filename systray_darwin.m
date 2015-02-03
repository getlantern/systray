#import <Cocoa/Cocoa.h>
#include "systray.h"

@interface MenuItem : NSObject
{
  @public
    NSString* menuId;
    NSString* title;
    NSString* tooltip;
}
-(id) initWithId: (const char*)theMenuId
       withTitle: (const char*)theTitle
     withTooltip: (const char*)theTooltip;
@end
@implementation MenuItem
-(id) initWithId: (const char*)theMenuId
       withTitle: (const char*)theTitle
     withTooltip: (const char*)theTooltip
{
  menuId = [[NSString alloc] initWithCString:theMenuId
                                     encoding:NSUTF8StringEncoding];
  title = [[NSString alloc] initWithCString:theTitle
                                    encoding:NSUTF8StringEncoding];
  tooltip = [[NSString alloc] initWithCString:theTooltip
                                      encoding:NSUTF8StringEncoding];
  return self;
}

-(void) dealloc {
  [menuId release];
  [title release];
  [tooltip release];
  [super dealloc];
}
@end

@interface AppDelegate: NSObject <NSApplicationDelegate>
   - (void) addMenuItem:(MenuItem*) item;
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

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification
{
  self->statusItem = [[[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength] retain];
  NSZone *menuZone = [NSMenu menuZone];
  self->menu = [[[NSMenu allocWithZone:menuZone] init] autorelease];
  [self->statusItem setMenu:self->menu];
  systray_ready();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification
{
  [[NSStatusBar systemStatusBar] removeStatusItem: statusItem];
}

- (void)setIcon:(NSImage *)image {
  [statusItem setImage:image];
  [image release];
}

- (void)setTitle:(NSString *)title {
  [statusItem setTitle:title];
  [title release];
}

- (void)setTooltip:(NSString *)tooltip {
  [statusItem setToolTip:tooltip];
  [tooltip release];
}

- (IBAction)menuHandler:(id)sender
{
  NSString* menuId = [sender representedObject];
  systray_menu_item_selected((char*)[menuId cStringUsingEncoding: NSUTF8StringEncoding]);
}

- (void) addMenuItem:(MenuItem*) item
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
  [item release];
}

- (void) quit
{
  [NSApp terminate:self];
}

@end

int nativeLoop(void) {
    AppDelegate *delegate = [[[AppDelegate alloc] init] autorelease];
    [[NSApplication sharedApplication] setDelegate:delegate];
    [NSApp run];
    return EXIT_SUCCESS;
}

void runInMainThread(SEL method, id object) {
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

void setTitle(char* ctitle) {
  NSString* title = [[NSString alloc] initWithCString:ctitle
                                              encoding:NSUTF8StringEncoding];
  free(ctitle);
  runInMainThread(@selector(setTitle:), (id)title);
}

void setTooltip(char* ctooltip) {
  NSString* tooltip = [[NSString alloc] initWithCString:ctooltip
                                                encoding:NSUTF8StringEncoding];
  free(ctooltip);
  runInMainThread(@selector(setTooltip:), (id)tooltip);
}

void addMenuItem(char* menuId, char* title, char* tooltip) {
  MenuItem* item = [[MenuItem alloc] initWithId: menuId withTitle: title withTooltip: tooltip] ;
  free(menuId);
  free(title);
  free(tooltip);
  runInMainThread(@selector(addMenuItem:), (id)item);
}

void quit() {
  runInMainThread(@selector(quit), nil);
}
