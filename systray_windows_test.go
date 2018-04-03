// +build windows

package systray

import (
	"testing"
	"sync/atomic"
	"time"
	"golang.org/x/sys/windows"
	"unsafe"
	"runtime"
)

func TestBaseWindowsTray(t *testing.T) {
	runtime.LockOSThread()

	if err := wt.initInstance(); err != nil {
		t.Fatalf("initInstance failed: %s", err)
	}

	if err := wt.createMenu(); err != nil {
		t.Fatalf("createMenu failed: %s", err)
	}

	defer func() {
		pDestroyWindow.Call(uintptr(wt.window))
		wt.wcex.Unregister()
	}()

	if err := wt.setIcon("example/icon/iconwin.ico"); err != nil {
		t.Errorf("SetIcon failed: %s", err)
	}

	if err := wt.setTooltip("Cyrillic tooltip тест:)"); err != nil {
		t.Errorf("SetIcon failed: %s", err)
	}

	var id int32 = 0
	err := wt.addOrUpdateMenuItem(atomic.AddInt32(&id, 1), "Simple enabled", false, false)
	if err != nil {
		t.Errorf("mergeMenuItem failed: %s", err)
	}
	err = wt.addOrUpdateMenuItem(atomic.AddInt32(&id, 1), "Simple disabled", true, false)
	if err != nil {
		t.Errorf("mergeMenuItem failed: %s", err)
	}
	err = wt.addSeparatorMenuItem(atomic.AddInt32(&id, 1))
	if err != nil {
		t.Errorf("addSeparatorMenuItem failed: %s", err)
	}
	err = wt.addOrUpdateMenuItem(atomic.AddInt32(&id, 1), "Simple checked enabled", false, true)
	if err != nil {
		t.Errorf("mergeMenuItem failed: %s", err)
	}
	err = wt.addOrUpdateMenuItem(atomic.AddInt32(&id, 1), "Simple checked disabled", true, true)
	if err != nil {
		t.Errorf("mergeMenuItem failed: %s", err)
	}

	err = wt.hideMenuItem(1)
	if err != nil {
		t.Errorf("hideMenuItem failed: %s", err)
	}

	err = wt.hideMenuItem(100)
	if err == nil {
		t.Error("hideMenuItem failed: must return error on invalid item id")
	}

	err = wt.addOrUpdateMenuItem(2, "Simple disabled update", true, false)
	if err != nil {
		t.Errorf("mergeMenuItem failed: %s", err)
	}

	time.AfterFunc(3*time.Second, quit)

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
			t.Errorf("win32 GetMessage failed: %v", err)
			return
		} else if res == 0 {
			break
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}
