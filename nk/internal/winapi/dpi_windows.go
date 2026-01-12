//go:build windows

package winapi

import (
	"syscall"
	"unsafe"
)

var (
	procGetDpiForWindow               = moduser32.NewProc("GetDpiForWindow")
	procSetProcessDpiAwarenessContext = moduser32.NewProc("SetProcessDpiAwarenessContext")
	procSetWindowPos                  = moduser32.NewProc("SetWindowPos")
	procGetDeviceCaps                 = modgdi32.NewProc("GetDeviceCaps")
)

// SetPerMonitorDPIAwareV2 enables per-monitor DPI awareness v2 if the API exists.
func SetPerMonitorDPIAwareV2() bool {
	if err := procSetProcessDpiAwarenessContext.Find(); err != nil {
		return false
	}
	r, _, _ := procSetProcessDpiAwarenessContext.Call(DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2)
	return r != 0
}

// DpiForWindow returns the current DPI for hwnd, using GetDpiForWindow when available.
// This is intended to be called rarely (e.g. once at window creation), not per-frame.
func DpiForWindow(hwnd syscall.Handle) uint32 {
	if err := procGetDpiForWindow.Find(); err == nil {
		r, _, _ := procGetDpiForWindow.Call(uintptr(hwnd))
		if r != 0 {
			return uint32(r)
		}
	}

	dc, _ := GetDC(hwnd)
	if dc != 0 {
		defer ReleaseDC(hwnd, dc) // best-effort
		r, _, _ := procGetDeviceCaps.Call(uintptr(dc), uintptr(LOGPIXELSX))
		if r != 0 {
			return uint32(r)
		}
	}
	return 96
}

// SetWindowPosFromRect applies the suggested window rect (e.g. from WM_DPICHANGED lParam).
func SetWindowPosFromRect(hwnd syscall.Handle, r *RECT, flags uint32) bool {
	if r == nil {
		return false
	}
	if err := procSetWindowPos.Find(); err != nil {
		return false
	}
	w := r.Right - r.Left
	h := r.Bottom - r.Top
	ret, _, _ := procSetWindowPos.Call(
		uintptr(hwnd),
		0,
		uintptr(r.Left),
		uintptr(r.Top),
		uintptr(w),
		uintptr(h),
		uintptr(flags),
	)
	return ret != 0
}

func RectFromLPARAM(lParam uintptr) *RECT {
	return (*RECT)(unsafe.Pointer(lParam))
}
