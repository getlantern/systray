// +build windows

package systray

import (
	"io/ioutil"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"path/filepath"
)

// Helpful sources: https://github.com/golang/exp/blob/master/shiny/driver/internal/win32

// https://msdn.microsoft.com/en-us/library/windows/desktop/bb762159
const (
	NIM_ADD    = 0x00000000
	NIM_MODIFY = 0x00000001
	NIM_DELETE = 0x00000002
)

const (
	NIF_MESSAGE = 0x00000001
	NIF_ICON    = 0x00000002
	NIF_TIP     = 0x00000004
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644931(v=vs.85).aspx
const (
	WM_USER       = 0x0400
	WM_COMMAND    = 0x0111
	WM_DESTROY    = 0x0002
	WM_ENDSESSION = 0x16
	WM_RBUTTONUP  = 0x0205
	WM_LBUTTONUP  = 0x0202

	// custom messages
	WM_SYSTRAY_MESSAGE = WM_USER + 1
)

var (
	k32                  = windows.NewLazySystemDLL("Kernel32.dll")
	s32                  = windows.NewLazySystemDLL("Shell32.dll")
	u32                  = windows.NewLazySystemDLL("User32.dll")
	pGetModuleHandle     = k32.NewProc("GetModuleHandleW")
	pShellNotifyIcon     = s32.NewProc("Shell_NotifyIconW")
	pCreatePopupMenu     = u32.NewProc("CreatePopupMenu")
	pCreateWindowEx      = u32.NewProc("CreateWindowExW")
	pDefWindowProc       = u32.NewProc("DefWindowProcW")
	pDeleteMenu          = u32.NewProc("DeleteMenu")
	pDestroyWindow       = u32.NewProc("DestroyWindow")
	pDispatchMessage     = u32.NewProc("DispatchMessageW")
	pGetCursorPos        = u32.NewProc("GetCursorPos")
	pGetMenuItemID       = u32.NewProc("GetMenuItemID")
	pGetMessage          = u32.NewProc("GetMessageW")
	pInsertMenuItem      = u32.NewProc("InsertMenuItemW")
	pLoadIcon            = u32.NewProc("LoadIconW")
	pLoadImage           = u32.NewProc("LoadImageW")
	pLoadCursor          = u32.NewProc("LoadCursorW")
	pPostMessage         = u32.NewProc("PostMessageW")
	pPostQuitMessage     = u32.NewProc("PostQuitMessage")
	pRegisterClass       = u32.NewProc("RegisterClassExW")
	pSetForegroundWindow = u32.NewProc("SetForegroundWindow")
	pSetMenuInfo         = u32.NewProc("SetMenuInfo")
	pSetMenuItemInfo     = u32.NewProc("SetMenuItemInfoW")
	pShowWindow          = u32.NewProc("ShowWindow")
	pTrackPopupMenu      = u32.NewProc("TrackPopupMenu")
	pTranslateMessage    = u32.NewProc("TranslateMessage")
	pUnregisterClass     = u32.NewProc("UnregisterClassW")
	pUpdateWindow        = u32.NewProc("UpdateWindow")
)

// Contains window class information.
// It is used with the RegisterClassEx and GetClassInfoEx functions.
// https://msdn.microsoft.com/en-us/library/ms633577.aspx
type wndClassEx struct {
	Size, Style                        uint32
	WndProc                            uintptr
	ClsExtra, WndExtra                 int32
	Instance, Icon, Cursor, Background windows.Handle
	MenuName, ClassName                *uint16
	IconSm                             windows.Handle
}

// Registers a window class for subsequent use in calls to the CreateWindow or CreateWindowEx function.
// https://msdn.microsoft.com/en-us/library/ms633587.aspx
func (w *wndClassEx) Register() error {
	w.Size = uint32(unsafe.Sizeof(*w))
	res, _, err := pRegisterClass.Call(uintptr(unsafe.Pointer(w)))
	if res == 0 {
		return err
	}
	return nil
}

// Unregisters a window class, freeing the memory required for the class.
// https://msdn.microsoft.com/en-us/library/ms644899.aspx
func (w *wndClassEx) Unregister() error {
	res, _, err := pUnregisterClass.Call(
		uintptr(unsafe.Pointer(w.ClassName)),
		uintptr(w.Instance),
	)
	if res == 0 {
		return err
	}
	return nil
}

// Contains information that the system needs to display notifications in the notification area.
// Used by Shell_NotifyIcon.
// https://msdn.microsoft.com/en-us/library/windows/desktop/bb773352(v=vs.85).aspx
type notifyIconData struct {
	Size                       uint32
	Wnd                        windows.Handle
	ID, Flags, CallbackMessage uint32
	Icon                       windows.Handle
	Tip                        [128]uint16
	State, StateMask           uint32
	Info                       [256]uint16
	Timeout, Version           uint32
	InfoTitle                  [64]uint16
	InfoFlags                  uint32
	GuidItem                   windows.GUID
	BalloonIcon                windows.Handle
}

// Contains information about a menu item
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647578(v=vs.85).aspx
type menuItemInfo struct {
	Size, Mask, Type, State     uint32
	ID                          uint32
	SubMenu, Checked, Unchecked windows.Handle
	ItemData                    uintptr
	TypeData                    *uint16
	Cch                         uint32
	Item                        windows.Handle
}

// The POINT structure defines the x- and y- coordinates of a point.
// https://msdn.microsoft.com/en-us/library/windows/desktop/dd162805(v=vs.85).aspx
type point struct {
	X, Y int32
}

// Contains information about loaded resources
type winTray struct {
	instance,
	icon,
	cursor,
	window,
	menu windows.Handle

	loadedImages map[string]windows.Handle
	nid          *notifyIconData
	wcex         *wndClassEx
}

// Loads an image from file and shows it in tray.
// LoadImage: https://msdn.microsoft.com/en-us/library/windows/desktop/ms648045(v=vs.85).aspx
// Shell_NotifyIcon: https://msdn.microsoft.com/en-us/library/windows/desktop/bb762159(v=vs.85).aspx
func (t *winTray) setIcon(src string) error {
	const (
		IMAGE_ICON = 1 // Loads an icon
	)
	const (
		LR_LOADFROMFILE = 0x00000010 // Loads the stand-alone image from the file
	)

	// Save and reuse handles of loaded images
	h, ok := t.loadedImages[src]
	if !ok {
		srcPtr, err := windows.UTF16PtrFromString(src)
		if err != nil {
			return err
		}
		res, _, err := pLoadImage.Call(
			0,
			uintptr(unsafe.Pointer(srcPtr)),
			IMAGE_ICON,
			64,
			64,
			LR_LOADFROMFILE,
		)
		if res == 0 {
			return err
		}
		h = windows.Handle(res)
		t.loadedImages[src] = h
	}

	t.nid.Icon = h
	t.nid.Flags = NIF_ICON
	t.nid.Size = uint32(unsafe.Sizeof(*t.nid))

	showIconRes, _, err := pShellNotifyIcon.Call(
		uintptr(NIM_MODIFY),
		uintptr(unsafe.Pointer(t.nid)),
	)
	if showIconRes == 0 {
		return err
	}

	return nil
}

// Sets tooltip on icon.
// Shell_NotifyIcon: https://msdn.microsoft.com/en-us/library/windows/desktop/bb762159(v=vs.85).aspx
func (t *winTray) setTooltip(src string) error {
	b, err := windows.UTF16FromString(src)
	if err != nil {
		return err
	}
	copy(t.nid.Tip[:], b[:])
	t.nid.Flags = NIF_TIP
	t.nid.Size = uint32(unsafe.Sizeof(*t.nid))
	showIconRes, _, err := pShellNotifyIcon.Call(
		uintptr(NIM_MODIFY),
		uintptr(unsafe.Pointer(t.nid)),
	)
	if showIconRes == 0 {
		return err
	}

	return nil
}

var wt winTray

// WindowProc callback function that processes messages sent to a window.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms633573(v=vs.85).aspx
func (t *winTray) wndProc(hWnd windows.Handle, message uint32, wParam, lParam uintptr) (lResult uintptr) {
	switch message {
	case WM_COMMAND:
		menuId := int32(wParam)
		if menuId != -1 {
			systrayMenuItemSelected(menuId)
		}
	case WM_DESTROY:
		// same as WM_ENDSESSION, but throws 0 exit code after all
		defer pPostQuitMessage.Call(uintptr(int32(0)))
		fallthrough
	case WM_ENDSESSION:
		if t.nid != nil {
			pShellNotifyIcon.Call(
				NIM_DELETE,
				uintptr(unsafe.Pointer(t.nid)),
			)
		}
		if systrayExit != nil {
			systrayExit()
		}
	case WM_SYSTRAY_MESSAGE:
		switch lParam {
		case WM_RBUTTONUP, WM_LBUTTONUP:
			t.showMenu()
		}
	default:
		// Calls the default window procedure to provide default processing for any window messages that an application does not process.
		// https://msdn.microsoft.com/en-us/library/windows/desktop/ms633572(v=vs.85).aspx
		lResult, _, _ = pDefWindowProc.Call(
			uintptr(hWnd),
			uintptr(message),
			uintptr(wParam),
			uintptr(lParam),
		)
	}
	return
}

func (t *winTray) initInstance() error {
	const IDI_APPLICATION = 32512
	const IDC_ARROW = 32512 // Standard arrow
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms633548(v=vs.85).aspx
	const SW_HIDE = 0
	const CW_USEDEFAULT = 0x80000000
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms632600(v=vs.85).aspx
	const (
		WS_CAPTION     = 0x00C00000
		WS_MAXIMIZEBOX = 0x00010000
		WS_MINIMIZEBOX = 0x00020000
		WS_OVERLAPPED  = 0x00000000
		WS_SYSMENU     = 0x00080000
		WS_THICKFRAME  = 0x00040000

		WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	)
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ff729176
	const (
		CS_HREDRAW = 0x0002
		CS_VREDRAW = 0x0001
	)

	const (
		className  = "SystrayClass"
		windowName = ""
	)
	wt.loadedImages = make(map[string]windows.Handle)

	instanceHandle, _, err := pGetModuleHandle.Call(0)
	if instanceHandle == 0 {
		return err
	}
	t.instance = windows.Handle(instanceHandle)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms648072(v=vs.85).aspx
	iconHandle, _, err := pLoadIcon.Call(0, uintptr(IDI_APPLICATION))
	if iconHandle == 0 {
		return err
	}
	t.icon = windows.Handle(iconHandle)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms648391(v=vs.85).aspx
	cursorHandle, _, err := pLoadCursor.Call(0, uintptr(IDC_ARROW))
	if cursorHandle == 0 {
		return err
	}
	t.cursor = windows.Handle(cursorHandle)

	classNamePtr, err := windows.UTF16PtrFromString(className)
	if err != nil {
		return err
	}

	windowNamePtr, err := windows.UTF16PtrFromString(windowName)
	if err != nil {
		return err
	}

	t.wcex = &wndClassEx{
		Style:      CS_HREDRAW | CS_VREDRAW,
		WndProc:    windows.NewCallback(t.wndProc),
		Instance:   t.instance,
		Icon:       t.icon,
		Cursor:     t.cursor,
		Background: windows.Handle(6), //(COLOR_WINDOW + 1)
		ClassName:  classNamePtr,
		IconSm:     t.icon,
	}
	if err := t.wcex.Register(); err != nil {
		return err
	}

	windowHandle, _, err := pCreateWindowEx.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(classNamePtr)),
		uintptr(unsafe.Pointer(windowNamePtr)),
		uintptr(WS_OVERLAPPEDWINDOW),
		uintptr(CW_USEDEFAULT),
		uintptr(CW_USEDEFAULT),
		uintptr(CW_USEDEFAULT),
		uintptr(CW_USEDEFAULT),
		uintptr(0),
		uintptr(0),
		uintptr(t.instance),
		uintptr(0),
	)
	if windowHandle == 0 {
		return err
	}
	t.window = windows.Handle(windowHandle)

	pShowWindow.Call(
		uintptr(t.window),
		uintptr(SW_HIDE),
	)

	pUpdateWindow.Call(
		uintptr(t.window),
	)

	wt.nid = &notifyIconData{
		Wnd:             windows.Handle(wt.window),
		ID:              100,
		Flags:           NIF_MESSAGE,
		CallbackMessage: WM_SYSTRAY_MESSAGE,
	}
	t.nid.Size = uint32(unsafe.Sizeof(*t.nid))

	showIconRes, _, err := pShellNotifyIcon.Call(
		uintptr(NIM_ADD),
		uintptr(unsafe.Pointer(t.nid)),
	)
	if showIconRes == 0 {
		return err
	}

	return nil
}

func (t *winTray) createMenu() error {
	const MIM_APPLYTOSUBMENUS = 0x80000000 // Settings apply to the menu and all of its submenus

	menuHandle, _, err := pCreatePopupMenu.Call()
	if menuHandle == 0 {
		return err
	}
	t.menu = windows.Handle(menuHandle)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647575(v=vs.85).aspx
	mi := struct {
		Size, Mask, Style, Max uint32
		Background             windows.Handle
		ContextHelpID          uint32
		MenuData               uintptr
	}{
		Mask: MIM_APPLYTOSUBMENUS,
	}
	mi.Size = uint32(unsafe.Sizeof(mi))

	res, _, err := pSetMenuInfo.Call(
		uintptr(t.menu),
		uintptr(unsafe.Pointer(&mi)),
	)
	if res == 0 {
		return err
	}
	return nil
}

func (t *winTray) addOrUpdateMenuItem(menuId int32, title string, disabled, checked bool) error {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647578(v=vs.85).aspx
	const (
		MIIM_FTYPE  = 0x00000100
		MIIM_STRING = 0x00000040
		MIIM_ID     = 0x00000002
		MIIM_STATE  = 0x00000001
	)
	const MFT_STRING = 0x00000000
	const (
		MFS_CHECKED  = 0x00000008
		MFS_DISABLED = 0x00000003
	)
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return err
	}

	mi := menuItemInfo{
		Mask:     MIIM_FTYPE | MIIM_STRING | MIIM_ID | MIIM_STATE,
		Type:     MFT_STRING,
		ID:       uint32(menuId),
		TypeData: titlePtr,
		Cch:      uint32(len(title)),
	}
	if disabled {
		mi.State |= MFS_DISABLED
	}
	if checked {
		mi.State |= MFS_CHECKED
	}
	mi.Size = uint32(unsafe.Sizeof(mi))

	// The return value is the identifier of the specified menu item.
	// If the menu item identifier is NULL or if the specified item opens a submenu, the return value is -1.
	res, _, err := pGetMenuItemID.Call(uintptr(t.menu), uintptr(menuId))
	if int32(res) == -1 {
		res, _, err = pInsertMenuItem.Call(
			uintptr(t.menu),
			uintptr(menuId),
			1,
			uintptr(unsafe.Pointer(&mi)),
		)
		if res == 0 {
			return err
		}
	} else {
		res, _, err = pSetMenuItemInfo.Call(
			uintptr(t.menu),
			uintptr(menuId),
			0,
			uintptr(unsafe.Pointer(&mi)),
		)
		if res == 0 {
			return err
		}
	}

	return nil
}

func (t *winTray) addSeparatorMenuItem(menuId int32) error {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647578(v=vs.85).aspx
	const (
		MIIM_FTYPE = 0x00000100
		MIIM_ID    = 0x00000002
		MIIM_STATE = 0x00000001
	)
	const MFT_SEPARATOR = 0x00000800

	mi := menuItemInfo{
		Mask: MIIM_FTYPE | MIIM_ID | MIIM_STATE,
		Type: MFT_SEPARATOR,
		ID:   uint32(menuId),
	}

	mi.Size = uint32(unsafe.Sizeof(mi))

	res, _, err := pInsertMenuItem.Call(
		uintptr(t.menu),
		uintptr(menuId),
		1,
		uintptr(unsafe.Pointer(&mi)),
	)
	if res == 0 {
		return err
	}

	return nil
}

func (t *winTray) hideMenuItem(menuId int32) error {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms647629(v=vs.85).aspx
	const MF_BYCOMMAND = 0x00000000

	res, _, err := pDeleteMenu.Call(
		uintptr(t.menu),
		uintptr(uint32(menuId)),
		MF_BYCOMMAND,
	)
	if res == 0 {
		return err
	}

	return nil
}

func (t *winTray) showMenu() error {
	const (
		TPM_BOTTOMALIGN = 0x0020
		TPM_LEFTALIGN   = 0x0000
	)
	p := point{}
	res, _, err := pGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	if res == 0 {
		return err
	}
	pSetForegroundWindow.Call(uintptr(t.window))

	res, _, err = pTrackPopupMenu.Call(
		uintptr(t.menu),
		TPM_BOTTOMALIGN|TPM_LEFTALIGN,
		uintptr(p.X),
		uintptr(p.Y),
		0,
		uintptr(t.window),
		0,
	)
	if res == 0 {
		return err
	}

	return nil
}

func nativeLoop() {
	if err := wt.initInstance(); err != nil {
		fmt.Printf("initInstance failed: %s", err)
	}

	if err := wt.createMenu(); err != nil {
		fmt.Printf("createMenu failed: %s", err)
		return
	}

	defer func() {
		pDestroyWindow.Call(uintptr(wt.window))
		wt.wcex.Unregister()
	}()

	if systrayReady != nil {
		systrayReady()
	}

	// Main message pump.
	m := struct {
		WindowHandle windows.Handle
		Message      uint32
		Wparam       uintptr
		Lparam       uintptr
		Time         uint32
		Pt           point
	}{}
	for {
		ret, _, err := pGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		res := int32(ret)
		if res == -1 {
			log.Errorf("win32 GetMessage failed: %v", err)
			return
		} else if res == 0 { // WM_QUIT
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func quit() {
	const WM_CLOSE = 0x0010

	pPostMessage.Call(
		uintptr(wt.window),
		WM_CLOSE,
		0,
		0,
	)
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	bh := md5.Sum(iconBytes)
	dataHash := hex.EncodeToString(bh[:])
	iconFilePath := filepath.Join(os.TempDir(), "systray_temp_icon_"+dataHash)

	if _, err := os.Stat(iconFilePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(iconFilePath, iconBytes, 0644); err != nil {
			log.Errorf("Unable to write icon data to temp file: %v", err)
			return
		}
	}

	if err := wt.setIcon(iconFilePath); err != nil {
		log.Errorf("Unable to set icon: %v", err)
		return
	}
}

// SetTitle sets the systray title, only available on Mac.
func SetTitle(title string) {
	// do nothing
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
	if err := wt.setTooltip(tooltip); err != nil {
		log.Errorf("Unable to set tooltip: %v", err)
		return
	}
}

func addOrUpdateMenuItem(item *MenuItem) {
	err := wt.addOrUpdateMenuItem(item.id, item.title, item.disabled, item.checked)
	if err != nil {
		log.Errorf("Unable to addOrUpdateMenuItem: %v", err)
		return
	}
}

func addSeparator(id int32) {
	err := wt.addSeparatorMenuItem(id)
	if err != nil {
		log.Errorf("Unable to addSeparator: %v", err)
		return
	}
}

func hideMenuItem(item *MenuItem) {
	err := wt.hideMenuItem(item.id)
	if err != nil {
		log.Errorf("Unable to hideMenuItem: %v", err)
		return
	}
}

func showMenuItem(item *MenuItem) {
	addOrUpdateMenuItem(item)
}
