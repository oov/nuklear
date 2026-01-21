package w32

/*
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

void poll_events() {
	MSG m;
	while (PeekMessageW(&m, 0, 0, 0, PM_REMOVE)) {
		TranslateMessage(&m);
		DispatchMessageW(&m);
	}
}
*/
import "C"

import (
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/golang-ui/nuklear/nk/internal/winapi"
)

const (
	MouseButtonLeft           = 1
	MouseButtonLeftModified   = 2
	MouseButtonLeftDouble     = 4
	MouseButtonRight          = 8
	MouseButtonRightModified  = 16
	MouseButtonRightDouble    = 32
	MouseButtonMiddle         = 64
	MouseButtonMiddleModified = 128
	MouseButtonMiddleDouble   = 256
)

// Cursor types for SetCursor
const (
	CursorArrow   = iota // Default arrow cursor
	CursorHand           // Hand/pointer cursor (for clickable items)
	CursorSizeWE         // Horizontal resize cursor
	CursorSizeNS         // Vertical resize cursor
	CursorSizeAll        // Move/resize all directions cursor
	CursorWait           // Wait/hourglass cursor
)

// System cursor IDs (Windows IDC_* constants)
const (
	idcArrow   = 32512
	idcHand    = 32649
	idcSizeWE  = 32644
	idcSizeNS  = 32645
	idcSizeAll = 32646
	idcWait    = 32514
)

var (
	// Cached system cursors
	systemCursors [6]syscall.Handle
)

func init() {
	// Load system cursors once
	ids := []uintptr{idcArrow, idcHand, idcSizeWE, idcSizeNS, idcSizeAll, idcWait}
	for i, id := range ids {
		systemCursors[i], _ = winapi.LoadCursor(0, id)
	}
}

func Init() error {
	// Best-effort: enable per-monitor DPI awareness so WM_DPICHANGED is delivered.
	winapi.SetPerMonitorDPIAwareV2()
	return nil
}

func Terminate() {

}

type Window struct {
	Handle syscall.Handle
	Keys   map[int]struct{}
	Mouse  struct {
		X      int
		Y      int
		Button int
	}
	Chars []rune

	dpi    uint32 // updated on WM_DPICHANGED
	cursor int    // current cursor type (CursorArrow, CursorHand, etc.)

	shouldClose             bool
	dropHandler             DropCallback
	sizeHandler             SizeCallback
	paintHandler            PaintCallback
	dpiHandler              DPIChangedCallback
	keyHandler              KeyCallback
	mouseButtonHandler      MouseButtonCallback
	mouseMoveHandler        MouseMoveCallback
	mouseWheelYHandler      MouseWheelCallback
	mouseDoubleClickHandler MouseDoubleClickCallback
}

// SetCursor sets the mouse cursor type for this window.
// Use CursorArrow, CursorHand, CursorSizeWE, etc.
func (w *Window) SetCursor(cursor int) {
	if cursor >= 0 && cursor < len(systemCursors) {
		w.cursor = cursor
	}
}

func (w *Window) wndProc(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) {
	switch uMsg {
	case winapi.WM_CLOSE:
		w.SetShouldClose(true)
		return 0
	case winapi.WM_SETCURSOR:
		// Handle cursor change only for client area (HTCLIENT)
		if winapi.LOWORD(lParam) == winapi.HTCLIENT {
			cursor := w.cursor
			if cursor >= 0 && cursor < len(systemCursors) && systemCursors[cursor] != 0 {
				winapi.SetCursor(systemCursors[cursor])
				return 1 // Indicate we handled the message
			}
		}
		// Let DefWindowProc handle non-client area cursors
	case winapi.WM_SIZE:
		// Fallback path: if WM_DPICHANGED wasn't delivered (e.g. non-top-level/parented window),
		// re-check DPI on resize and fire the handler on change.
		{
			dpi := winapi.DpiForWindow(hwnd)
			prev := atomic.LoadUint32(&w.dpi)
			if dpi != 0 && dpi != prev {
				atomic.StoreUint32(&w.dpi, dpi)
				if w.dpiHandler != nil {
					w.dpiHandler(w, dpi)
				}
			}
		}
		if w.sizeHandler != nil {
			w.sizeHandler(w, int(winapi.LOWORD(lParam)), int(winapi.HIWORD(lParam)))
		}
	case winapi.WM_DPICHANGED:
		dpi := uint32(winapi.LOWORD(wParam))
		if dpi != 0 {
			atomic.StoreUint32(&w.dpi, dpi)
			if w.dpiHandler != nil {
				w.dpiHandler(w, dpi)
			}
		}
		// Apply suggested window rect (per Microsoft guidance).
		r := winapi.RectFromLPARAM(lParam)
		winapi.SetWindowPosFromRect(hwnd, r, winapi.SWP_NOZORDER|winapi.SWP_NOACTIVATE)
		return 0
	case winapi.WM_PAINT:
		if w.paintHandler != nil {
			var ps winapi.PAINTSTRUCT
			hdc := winapi.BeginPaint(hwnd, &ps)
			w.paintHandler(w, hdc)
			winapi.EndPaint(hwnd, &ps)
			return 1
		}
	case winapi.WM_ERASEBKGND:
		return 1
	case winapi.WM_CHAR:
		w.Chars = append(w.Chars, rune(wParam))
		return 0
	case winapi.WM_KILLFOCUS:
		w.Keys = map[int]struct{}{}
		w.Mouse.Button = MouseButtonLeftModified | MouseButtonRightModified | MouseButtonMiddleModified
		return 0
	case winapi.WM_KEYDOWN, winapi.WM_KEYUP:
		switch key := int(wParam); key {
		case winapi.VK_SHIFT,
			winapi.VK_CONTROL,
			winapi.VK_DELETE,
			winapi.VK_RETURN,
			winapi.VK_TAB,
			winapi.VK_LEFT,
			winapi.VK_RIGHT,
			winapi.VK_BACK,
			winapi.VK_HOME,
			winapi.VK_END,
			winapi.VK_NEXT,
			winapi.VK_PRIOR,
			winapi.VK_C,
			winapi.VK_V,
			winapi.VK_X,
			winapi.VK_Z,
			winapi.VK_R,
			winapi.VK_B,
			winapi.VK_E:
			down := lParam&0x80000000 == 0
			if down {
				w.Keys[key] = struct{}{}
			} else {
				delete(w.Keys, key)
			}
			if w.keyHandler != nil {
				w.keyHandler(w, key, down)
			}
			return 1
		}
	case winapi.WM_LBUTTONDBLCLK:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button |= MouseButtonLeftDouble
		if w.mouseDoubleClickHandler != nil {
			w.mouseDoubleClickHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_LBUTTON)
		}
		return 1
	case winapi.WM_LBUTTONDOWN:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button |= MouseButtonLeft + MouseButtonLeftModified
		if w.mouseButtonHandler != nil {
			w.mouseButtonHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_LBUTTON, true)
		}
		winapi.SetCapture(hwnd)
		return 1
	case winapi.WM_LBUTTONUP:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button = (w.Mouse.Button | MouseButtonLeftModified) &^ MouseButtonLeft
		if w.mouseButtonHandler != nil {
			w.mouseButtonHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_LBUTTON, false)
		}
		winapi.ReleaseCapture()
		return 1
	case winapi.WM_RBUTTONDBLCLK:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button |= MouseButtonRightDouble
		if w.mouseDoubleClickHandler != nil {
			w.mouseDoubleClickHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_RBUTTON)
		}
		return 1
	case winapi.WM_RBUTTONDOWN:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button |= MouseButtonRight + MouseButtonRightModified
		if w.mouseButtonHandler != nil {
			w.mouseButtonHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_RBUTTON, true)
		}
		winapi.SetCapture(hwnd)
		return 1
	case winapi.WM_RBUTTONUP:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button = (w.Mouse.Button | MouseButtonRightModified) &^ MouseButtonRight
		if w.mouseButtonHandler != nil {
			w.mouseButtonHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_RBUTTON, false)
		}
		winapi.ReleaseCapture()
		return 1
	case winapi.WM_MBUTTONDBLCLK:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button |= MouseButtonMiddleDouble
		if w.mouseDoubleClickHandler != nil {
			w.mouseDoubleClickHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_MBUTTON)
		}
		return 1
	case winapi.WM_MBUTTONDOWN:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button |= MouseButtonMiddle + MouseButtonMiddleModified
		if w.mouseButtonHandler != nil {
			w.mouseButtonHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_MBUTTON, true)
		}
		winapi.SetCapture(hwnd)
		return 1
	case winapi.WM_MBUTTONUP:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		w.Mouse.Button = (w.Mouse.Button | MouseButtonMiddleModified) &^ MouseButtonMiddle
		if w.mouseButtonHandler != nil {
			w.mouseButtonHandler(w, w.Mouse.X, w.Mouse.Y, winapi.VK_MBUTTON, false)
		}
		winapi.ReleaseCapture()
		return 1
	case winapi.WM_MOUSEMOVE:
		w.Mouse.X = int(int16(winapi.LOWORD(lParam)))
		w.Mouse.Y = int(int16(winapi.HIWORD(lParam)))
		if w.mouseMoveHandler != nil {
			w.mouseMoveHandler(w, w.Mouse.X, w.Mouse.Y)
		}
		return 1
	case winapi.WM_MOUSEWHEEL:
		if w.mouseWheelYHandler != nil {
			w.mouseWheelYHandler(w, float32(int16(winapi.HIWORD(wParam)))/winapi.WHEEL_DELTA)
		}
		return 1
	case winapi.WM_DROPFILES:
		drop := syscall.Handle(wParam)
		n := winapi.DragQueryFile(drop, 0xFFFFFFFF, nil, 0)
		files := make([]string, n)
		for i := 0; i < n; i++ {
			l := winapi.DragQueryFile(drop, i, nil, 0) + 1
			s := make([]uint16, l)
			winapi.DragQueryFile(drop, i, &s[0], l)
			files[i] = syscall.UTF16ToString(s)
		}
		winapi.DragFinish(drop)
		if w.dropHandler != nil {
			w.dropHandler(w, files)
		}
		return 0
	}
	return winapi.DefWindowProc(hwnd, uMsg, wParam, lParam)
}

func CreateWindow(width, height int, title string, dummy *int, dummy2 *int) (*Window, error) {
	winClsName, err := syscall.UTF16PtrFromString("NuklearWindowClass")
	if err != nil {
		return nil, err
	}
	hInstance, err := winapi.GetModuleHandle(nil)
	if err != nil {
		return nil, err
	}
	hIcon, err := winapi.LoadIcon(0, winapi.IDI_APPLICATION)
	if err != nil {
		return nil, err
	}
	hCursor, err := winapi.LoadCursor(0, winapi.IDC_ARROW)
	if err != nil {
		return nil, err
	}
	hBrush, err := winapi.GetStockObject(winapi.NULL_BRUSH)
	if err != nil {
		return nil, err
	}
	w := &Window{
		Keys: map[int]struct{}{},
		dpi:  96,
	}
	atm, err := winapi.RegisterClassEx(&winapi.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(winapi.WNDCLASSEX{})),
		Style:         0, //winapi.CS_DBLCLKS,
		LpfnWndProc:   syscall.NewCallback(w.wndProc),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     hInstance,
		HIcon:         hIcon,
		HCursor:       hCursor,
		HbrBackground: hBrush,
		LpszMenuName:  nil,
		LpszClassName: winClsName,
		HIconSm:       hIcon,
	})
	if err != nil {
		return nil, err
	}

	titleStr, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return nil, err
	}
	style := uint32(winapi.WS_OVERLAPPEDWINDOW)
	exStyle := uint32(winapi.WS_EX_APPWINDOW)
	rect := winapi.RECT{0, 0, int32(width), int32(height)}
	if err = winapi.AdjustWindowRectEx(&rect, style, 0, exStyle); err != nil {
		return nil, err
	}
	h, err := winapi.CreateWindowEx(exStyle,
		(*uint16)(unsafe.Pointer(uintptr(winapi.MAKELONG(atm, 0)))),
		titleStr,
		style,
		winapi.CW_USEDEFAULT, winapi.CW_USEDEFAULT,
		rect.Right-rect.Left, rect.Bottom-rect.Top,
		0, 0, hInstance, 0)
	if err != nil {
		return nil, err
	}
	w.Handle = h
	atomic.StoreUint32(&w.dpi, winapi.DpiForWindow(w.Handle))
	return w, nil
}

type DropCallback func(w *Window, names []string)

func (w *Window) SetDropCallback(f DropCallback) (prev DropCallback) {
	prev = w.dropHandler
	w.dropHandler = f
	winapi.DragAcceptFiles(w.Handle, f != nil)
	return prev
}

type SizeCallback func(w *Window, width, height int)

func (w *Window) SetSizeCallback(f SizeCallback) (prev SizeCallback) {
	prev = w.sizeHandler
	w.sizeHandler = f
	return prev
}

type PaintCallback func(w *Window, hdc syscall.Handle)

func (w *Window) SetPaintCallback(f PaintCallback) (prev PaintCallback) {
	prev = w.paintHandler
	w.paintHandler = f
	return prev
}

type DPIChangedCallback func(w *Window, dpi uint32)

func (w *Window) SetDPIChangedCallback(f DPIChangedCallback) (prev DPIChangedCallback) {
	prev = w.dpiHandler
	w.dpiHandler = f
	return prev
}

func (w *Window) DPI() uint32 {
	dpi := atomic.LoadUint32(&w.dpi)
	if dpi == 0 {
		return 96
	}
	return dpi
}

func (w *Window) Scale() float32 {
	return float32(w.DPI()) / 96.0
}

type KeyCallback func(w *Window, key int, down bool)

func (w *Window) SetKeyCallback(f KeyCallback) (prev KeyCallback) {
	prev = w.keyHandler
	w.keyHandler = f
	return prev
}

type MouseButtonCallback func(w *Window, x, y int, button int, down bool)

func (w *Window) SetMouseButtonCallback(f MouseButtonCallback) (prev MouseButtonCallback) {
	prev = w.mouseButtonHandler
	w.mouseButtonHandler = f
	return prev
}

type MouseMoveCallback func(w *Window, x, y int)

func (w *Window) SetMouseMoveCallback(f MouseMoveCallback) (prev MouseMoveCallback) {
	prev = w.mouseMoveHandler
	w.mouseMoveHandler = f
	return prev
}

type MouseWheelCallback func(w *Window, v float32)

func (w *Window) SetMouseWheelYCallback(f MouseWheelCallback) (prev MouseWheelCallback) {
	prev = w.mouseWheelYHandler
	w.mouseWheelYHandler = f
	return prev
}

type MouseDoubleClickCallback func(w *Window, x, y int, button int)

func (w *Window) SetMouseDoubleClickCallback(f MouseDoubleClickCallback) (prev MouseDoubleClickCallback) {
	prev = w.mouseDoubleClickHandler
	w.mouseDoubleClickHandler = f
	return prev
}

func (w *Window) GetSize() (width, height int) {
	var r winapi.RECT
	if err := winapi.GetClientRect(w.Handle, &r); err != nil {
		panic(err.Error())
	}
	return int(r.Right - r.Left), int(r.Bottom - r.Top)
}

func (w *Window) Show() error {
	winapi.ShowWindow(w.Handle, winapi.SW_SHOW)
	return nil
}

func (w *Window) Hide() error {
	winapi.ShowWindow(w.Handle, winapi.SW_HIDE)
	return nil
}

func (w *Window) Visible() bool {
	return winapi.IsWindowVisible(w.Handle)
}

func (w *Window) ShouldClose() bool {
	return w.shouldClose
}

func (w *Window) SetShouldClose(v bool) {
	w.shouldClose = v
}

func PollEvents() {
	// If implements message loop at Go, its may cause a crash for some reasons...
	C.poll_events()
	// var m winapi.MSG
	// for winapi.PeekMessage(&m, 0, 0, 0, winapi.PM_REMOVE) {
	// 	winapi.TranslateMessage(&m)
	// 	winapi.DispatchMessage(&m)
	// }
}
