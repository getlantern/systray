extern void systray_clicked_call_back();
extern void systray_menu_item_call_back(char* menu_id);
int nativeLoop(void);

void setIcon(const char* iconBytes, int length);
void setTitle(const char* title);
void setTooltip(const char* tooltip);
void addMenu(char* menuId, char* title, char* tooltip);
void quit();
