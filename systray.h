extern void systray_ready();
extern void systray_on_exit();
extern void systray_menu_item_selected(int menu_id);
int nativeLoop(void);

void setIcon(const char* iconBytes, int length);
void setTitle(char* title);
void setTooltip(char* tooltip);
void add_or_update_menu_item(int menuId, char* title, char* tooltip, short disabled, short checked);
void add_separator(int menuId);
void hide_menu_item(int menuId);
void show_menu_item(int menuId);
void quit();
