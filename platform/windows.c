#include <stdio.h>
#include <windows.h>
#include <shellapi.h>

// Below are windows macros but not found in current header file,
// so we define it manually
#ifndef WM_MENUCOMMAND
#define WM_MENUCOMMAND 0x0126
#endif

#ifndef MIM_APPLYTOSUBMENUS
#define MIM_APPLYTOSUBMENUS 0x80000000
#endif

#ifndef MIM_STYLE
#define MIM_STYLE 0x00000010
#endif

#ifndef MNS_NOTIFYBYPOS
#define MNS_NOTIFYBYPOS 0x08000000
#endif


// Our own macros
#define WM_SYSTRAY_MESSAGE (WM_USER + 1)

#define MAX_LOADSTRING 100

HINSTANCE hInst;
TCHAR szTitle[MAX_LOADSTRING];
TCHAR szWindowClass[MAX_LOADSTRING];
NOTIFYICONDATA nid;
HWND hWnd;
HMENU hSubMenu;

HWND                InitInstance(HINSTANCE, int);
LRESULT CALLBACK    WndProc(HWND, UINT, WPARAM, LPARAM);

ATOM MyRegisterClass(HINSTANCE hInstance)
{
	WNDCLASSEX wcex;

	wcex.cbSize = sizeof(WNDCLASSEX);

	wcex.style          = CS_HREDRAW | CS_VREDRAW;
	wcex.lpfnWndProc    = WndProc;
	wcex.cbClsExtra     = 0;
	wcex.cbWndExtra     = 0;
	wcex.hInstance      = hInstance;
	wcex.hIcon          = LoadIcon(NULL, IDI_APPLICATION);
	wcex.hCursor        = LoadCursor(NULL, IDC_ARROW);
	wcex.hbrBackground  = (HBRUSH)(COLOR_WINDOW+1);
	wcex.lpszMenuName   = 0;
	wcex.lpszClassName  = szWindowClass;
	wcex.hIconSm        = LoadIcon(NULL, IDI_APPLICATION);

	return RegisterClassEx(&wcex);
}

HWND InitInstance(HINSTANCE hInstance, int nCmdShow)
{
	HWND hWnd;

	hInst = hInstance;

	hWnd = CreateWindow(szWindowClass, szTitle, WS_OVERLAPPEDWINDOW,
			CW_USEDEFAULT, 0, CW_USEDEFAULT, 0, NULL, NULL, hInstance, NULL);

	if (!hWnd)
	{
		return 0;
	}

	ShowWindow(hWnd, nCmdShow);
	UpdateWindow(hWnd);

	return hWnd;
}

void ShowMenu(HWND hWnd)
{
	POINT p;
	GetCursorPos(&p);
	SetForegroundWindow(hWnd); // Win32 bug work-around
	TrackPopupMenu(hSubMenu, TPM_BOTTOMALIGN | TPM_LEFTALIGN, p.x, p.y, 0, hWnd, NULL);

}

char* GetMenuItemId(int index) {
	MENUITEMINFO menuItemInfo;
	menuItemInfo.cbSize = sizeof(MENUITEMINFO);
	menuItemInfo.fMask = MIIM_DATA;
	GetMenuItemInfo(hSubMenu, index, TRUE, &menuItemInfo);
	return (char*)menuItemInfo.dwItemData;
}

LRESULT CALLBACK WndProc(HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam)
{
	switch (message)
	{
		case WM_MENUCOMMAND:
			systray_menu_item_selected(GetMenuItemId(wParam));
			break;
		case WM_DESTROY:
			Shell_NotifyIcon(NIM_DELETE, &nid);
			PostQuitMessage(0);
			break;
		case WM_SYSTRAY_MESSAGE:
			switch(lParam)
			{
				case WM_RBUTTONUP:
					ShowMenu(hWnd);
					break;
				case WM_LBUTTONUP:
					ShowMenu(hWnd);
					break;
				default:
					return DefWindowProc(hWnd, message, wParam, lParam);
			};
			break;
		default:
			return DefWindowProc(hWnd, message, wParam, lParam);
	}
	return 0;
}


int nativeLoop(void) {
	HINSTANCE hInstance = GetModuleHandle(NULL);

	wcscpy((wchar_t*)szWindowClass, (wchar_t*)TEXT("MyClass"));
	MyRegisterClass(hInstance);

	hWnd = InitInstance(hInstance, FALSE); // Don't show window
	if (!hWnd)
	{
		return;
	}

	hSubMenu = CreatePopupMenu();
	MENUINFO menuInfo;
	menuInfo.cbSize = sizeof(MENUINFO);
	menuInfo.fMask = MIM_APPLYTOSUBMENUS | MIM_STYLE;
	menuInfo.dwStyle = MNS_NOTIFYBYPOS;
	SetMenuInfo(hSubMenu, &menuInfo);

	nid.cbSize = sizeof(NOTIFYICONDATA);
	nid.hWnd = hWnd;
	nid.uID = 100;
	nid.uCallbackMessage = WM_SYSTRAY_MESSAGE;
	nid.uFlags = NIF_MESSAGE;
	Shell_NotifyIcon(NIM_ADD, &nid);

	systray_ready();

	// Main message loop:
	MSG msg;
	while (GetMessage(&msg, NULL, 0, 0))
	{
		TranslateMessage(&msg);
		DispatchMessage(&msg);
	}   
	return EXIT_SUCCESS;
}


void setIcon(const char* iconBytes, int length) {
	// Let's load up the tray icon
	HICON hIcon;
	{
		// This is really hacky, but LoadImage won't let me load an image from memory.
		// So we have to write out a temporary file, load it from there, then delete the file.

		// From http://msdn.microsoft.com/en-us/library/windows/desktop/aa363875.aspx
		TCHAR szTempFileName[MAX_PATH+1];
		TCHAR lpTempPathBuffer[MAX_PATH+1];
		int dwRetVal = GetTempPath(MAX_PATH+1,        // length of the buffer
				lpTempPathBuffer); // buffer for path
		if (dwRetVal > MAX_PATH+1 || (dwRetVal == 0))
		{
			return; // Failure
		}

		//  Generates a temporary file name.
		int uRetVal = GetTempFileName(lpTempPathBuffer, // directory for tmp files
				TEXT("_tmpicon"), // temp file name prefix
				0,                // create unique name
				szTempFileName);  // buffer for name
		if (uRetVal == 0)
		{
			return; // Failure
		}

		// Dump the icon to the temp file
		FILE* fIcon = fopen(szTempFileName, "wb");
		fwrite(iconBytes, 1, length, fIcon);
		fclose(fIcon);

		// Load the image from the file
		hIcon = LoadImage(NULL, szTempFileName, IMAGE_ICON, 64, 64, LR_LOADFROMFILE);

		// Delete the temp file
		remove(szTempFileName);
	}

	nid.hIcon = hIcon;
	nid.uFlags = NIF_ICON;
	Shell_NotifyIcon(NIM_MODIFY, &nid);
}

// Don't support for Windows
void setTitle(char* ctitle) {
	free(ctitle);
}

void setTooltip(char* ctooltip) {
	strcpy(nid.szTip, ctooltip); // MinGW seems to use ANSI
	nid.uFlags = NIF_TIP;
	Shell_NotifyIcon(NIM_MODIFY, &nid);
	free(ctooltip);
}

void addMenuItem(char* menuId, char* title, char* tooltip) {
	MENUITEMINFO menuItemInfo;
	menuItemInfo.cbSize = sizeof(MENUITEMINFO);
	menuItemInfo.fMask = MIIM_FTYPE | MIIM_STRING | MIIM_DATA;
	menuItemInfo.fType = MFT_STRING;
	menuItemInfo.dwTypeData = title;
	menuItemInfo.cch = strlen(title) + 1;
	menuItemInfo.dwItemData = (ULONG_PTR)menuId;

	int itemCount = GetMenuItemCount(hSubMenu);
	int i;
	for (i = 0; i < itemCount; i++) {
		char * idString = GetMenuItemId(i);
		if (strcmp(menuId, idString) == 0) {
			break;
		}
	}
	if (i < itemCount) {
		SetMenuItemInfo(hSubMenu, i, TRUE, &menuItemInfo);
	} else {
		InsertMenuItem(hSubMenu, 0, TRUE, &menuItemInfo);
	}

	// free(menuId);
	free(title);
	free(tooltip);
}

void quit() {
	PostMessage(hWnd, WM_DESTROY, 0, 0);
}
