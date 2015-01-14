extern void callMe(char* name);
extern void (*theCallMe)(const char* name);

int nativeLoop(void);

void addMenu(char* name, char* title, char* tooltip);

void quit();
