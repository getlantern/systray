#import <Cocoa/Cocoa.h>
#include "systray.h"

#if __MAC_OS_X_VERSION_MIN_REQUIRED < __MAC_10_14

    #ifndef NSControlStateValueOff
      #define NSControlStateValueOff NSOffState
    #endif

    #ifndef NSControlStateValueOn
      #define NSControlStateValueOn NSOnState
    #endif

#endif

@interface MenuItem : NSObject
{
  @public
    NSNumber* menuId;
    NSString* title;
    NSString* tooltip;
    short disabled;
    short checked;
}
-(id) initWithId: (int)theMenuId
       withTitle: (const char*)theTitle
     withTooltip: (const char*)theTooltip
    withDisabled: (short)theDisabled
     withChecked: (short)theChecked;
     @end
     @implementation MenuItem
     -(id) initWithId: (int)theMenuId
            withTitle: (const char*)theTitle
          withTooltip: (const char*)theTooltip
         withDisabled: (short)theDisabled
          withChecked: (short)theChecked
{
  menuId = [NSNumber numberWithInt:theMenuId];
  title = [[NSString alloc] initWithCString:theTitle
                                   encoding:NSUTF8StringEncoding];
  tooltip = [[NSString alloc] initWithCString:theTooltip
                                     encoding:NSUTF8StringEncoding];
  disabled = theDisabled;
  checked = theChecked;
  return self;
}
@end

@interface AppDelegate: NSObject <NSApplicationDelegate>
  - (void) add_or_update_menu_item:(MenuItem*) item;
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
  self->statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
  self->menu = [[NSMenu alloc] init];
  [self->menu setAutoenablesItems: FALSE];
  [self->statusItem setMenu:self->menu];
  systray_ready();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification
{
  systray_on_exit();
}

- (void)setIcon:(NSImage *)image {
  statusItem.button.image = image;
  [self updateTitleButtonStyle];
}

- (void)setTitle:(NSString *)title {
  statusItem.button.title = title;
  [self updateTitleButtonStyle];
}

-(void)updateTitleButtonStyle {
  if (statusItem.button.image != nil) {
    if ([statusItem.button.title length] == 0) {
      statusItem.button.imagePosition = NSImageOnly;
    } else {
      statusItem.button.imagePosition = NSImageLeft;
    }
  } else {
    statusItem.button.imagePosition = NSNoImage;
  }
}


- (void)setTooltip:(NSString *)tooltip {
  statusItem.button.toolTip = tooltip;
}

- (IBAction)menuHandler:(id)sender
{
  NSNumber* menuId = [sender representedObject];
  systray_menu_item_selected(menuId.intValue);
}

- (void) add_or_update_menu_item:(MenuItem*) item
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
  if (item->disabled == 1) {
    menuItem.enabled = FALSE;
  } else {
    menuItem.enabled = TRUE;
  }
  if (item->checked == 1) {
    menuItem.state = NSControlStateValueOn;
  } else {
    menuItem.state = NSControlStateValueOff;
  }
}

- (void) add_separator:(NSNumber*) menuId
{
  [menu addItem: [NSMenuItem separatorItem]];
}

- (void) hide_menu_item:(NSNumber*) menuId
{
  NSMenuItem* menuItem;
  int existedMenuIndex = [menu indexOfItemWithRepresentedObject: menuId];
  if (existedMenuIndex == -1) {
    return;
  }
  menuItem = [menu itemAtIndex: existedMenuIndex];
  [menuItem setHidden:TRUE];
}

- (void)setMenuItemIcon:(NSArray*)imageAndMenuId {
  NSImage* image = [imageAndMenuId objectAtIndex:0];
  NSNumber* menuId = [imageAndMenuId objectAtIndex:1];

  NSMenuItem* menuItem;
  int existedMenuIndex = [menu indexOfItemWithRepresentedObject: menuId];
  if (existedMenuIndex == -1) {
    return;
  }
  menuItem = [menu itemAtIndex: existedMenuIndex];
  menuItem.image = image;
}

- (void) show_menu_item:(NSNumber*) menuId
{
  NSMenuItem* menuItem;
  int existedMenuIndex = [menu indexOfItemWithRepresentedObject: menuId];
  if (existedMenuIndex == -1) {
    return;
  }
  menuItem = [menu itemAtIndex: existedMenuIndex];
  [menuItem setHidden:FALSE];
}

- (void) quit
{
  [NSApp terminate:self];
}

@end

int nativeLoop(void) {
  AppDelegate *delegate = [[AppDelegate alloc] init];
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
  [image setSize:NSMakeSize(16, 16)];
  runInMainThread(@selector(setIcon:), (id)image);
}

void setMenuItemIcon(const char* iconBytes, int length, int menuId) {
  NSData* buffer = [NSData dataWithBytes: iconBytes length:length];
  NSImage *image = [[NSImage alloc] initWithData:buffer];
  [image setSize:NSMakeSize(16, 16)];

  NSNumber *mId = [NSNumber numberWithInt:menuId];
  runInMainThread(@selector(setMenuItemIcon:), @[image, (id)mId]);
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

void add_or_update_menu_item(int menuId, char* title, char* tooltip, short disabled, short checked) {
  MenuItem* item = [[MenuItem alloc] initWithId: menuId withTitle: title withTooltip: tooltip withDisabled: disabled withChecked: checked];
  free(title);
  free(tooltip);
  runInMainThread(@selector(add_or_update_menu_item:), (id)item);
}

void add_separator(int menuId) {
  NSNumber *mId = [NSNumber numberWithInt:menuId];
  runInMainThread(@selector(add_separator:), (id)mId);
}

void hide_menu_item(int menuId) {
  NSNumber *mId = [NSNumber numberWithInt:menuId];
  runInMainThread(@selector(hide_menu_item:), (id)mId);
}

void show_menu_item(int menuId) {
  NSNumber *mId = [NSNumber numberWithInt:menuId];
  runInMainThread(@selector(show_menu_item:), (id)mId);
}

void quit() {
  runInMainThread(@selector(quit), nil);
}
