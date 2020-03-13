#include "stdbool.h"

extern void systray_ready();
extern void systray_on_exit();
extern void systray_menu_item_selected(int menu_id);
int nativeLoop(char* title, int width, int height);

void setIcon(const char* iconBytes, int length, bool template);
void setMenuItemIcon(const char* iconBytes, int length, int menuId, bool template);
void setTitle(char* title);
void setTooltip(char* tooltip);
void configureAppWindow(char* title, int width, int height);
void showAppWindow(char* url);
void add_or_update_menu_item(int menuId, char* title, char* tooltip, short disabled, short checked);
void add_or_update_submenu_item(int parent,int menuId, char* title, char* tooltip, short disabled, short checked);
void add_separator(int menuId);
void hide_menu_item(int menuId);
void show_menu_item(int menuId);
void quit();