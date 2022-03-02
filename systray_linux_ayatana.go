// +build linux,!legacy_appindicator
//go:build linux && !legacy_appindicator

package systray

/*
#cgo linux pkg-config: ayatana-appindicator3-0.1

#include "systray.h"
*/
import "C"
