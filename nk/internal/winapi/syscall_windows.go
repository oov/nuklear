package winapi

import (
	"errors"
	"syscall"
)

const (
	IDI_APPLICATION = 32512

	IDC_ARROW = 32512

	BLACK_BRUSH = 4
	NULL_BRUSH  = 5

	CS_DBLCLKS = 8

	WS_OVERLAPPEDWINDOW = 0xcf0000

	WS_EX_APPWINDOW = 0x40000

	CW_USEDEFAULT = 0x80000000 - 0x100000000

	SW_HIDE = 0
	SW_SHOW = 5

	FW_DONTCARE   = 0
	FW_THIN       = 100
	FW_EXTRALIGHT = 200
	FW_ULTRALIGHT = 200
	FW_LIGHT      = 300
	FW_NORMAL     = 400
	FW_REGULAR    = 400
	FW_MEDIUM     = 500
	FW_SEMIBOLD   = 600
	FW_DEMIBOLD   = 600
	FW_BOLD       = 700
	FW_EXTRABOLD  = 800
	FW_ULTRABOLD  = 800
	FW_HEAVY      = 900
	FW_BLACK      = 900

	DEFAULT_CHARSET     = 1
	OUT_DEFAULT_PRECIS  = 0
	CLIP_DEFAULT_PRECIS = 0
	DEFAULT_QUALITY     = 0
	DEFAULT_PITCH       = 0
	FF_DONTCARE         = 0

	WM_SIZE          = 5
	WM_KILLFOCUS     = 8
	WM_PAINT         = 15
	WM_CLOSE         = 16
	WM_ERASEBKGND    = 20
	WM_KEYDOWN       = 256
	WM_KEYUP         = 257
	WM_CHAR          = 258
	WM_SYSKEYDOWN    = 260
	WM_SYSKEYUP      = 261
	WM_LBUTTONDBLCLK = 515
	WM_LBUTTONDOWN   = 513
	WM_LBUTTONUP     = 514
	WM_RBUTTONDBLCLK = 518
	WM_RBUTTONDOWN   = 516
	WM_RBUTTONUP     = 517
	WM_MBUTTONDBLCLK = 521
	WM_MBUTTONDOWN   = 519
	WM_MBUTTONUP     = 520
	WM_MOUSEMOVE     = 512
	WM_MOUSEWHEEL    = 522
	WM_DROPFILES     = 563

	WHEEL_DELTA = 120

	PM_REMOVE = 1

	VK_LBUTTON    = 1
	VK_RBUTTON    = 2
	VK_CANCEL     = 3
	VK_MBUTTON    = 4
	VK_XBUTTON1   = 5
	VK_XBUTTON2   = 6
	VK_BACK       = 8
	VK_TAB        = 9
	VK_CLEAR      = 12
	VK_RETURN     = 13
	VK_SHIFT      = 16
	VK_CONTROL    = 17
	VK_MENU       = 18
	VK_PAUSE      = 19
	VK_CAPITAL    = 20
	VK_KANA       = 21
	VK_HANGEUL    = 21
	VK_HANGUL     = 21
	VK_JUNJA      = 23
	VK_FINAL      = 24
	VK_HANJA      = 25
	VK_KANJI      = 25
	VK_ESCAPE     = 27
	VK_CONVERT    = 28
	VK_NONCONVERT = 29
	VK_ACCEPT     = 30
	VK_MODECHANGE = 31
	VK_SPACE      = 32
	VK_PRIOR      = 33
	VK_NEXT       = 34
	VK_END        = 35
	VK_HOME       = 36
	VK_LEFT       = 37
	VK_UP         = 38
	VK_RIGHT      = 39
	VK_DOWN       = 40
	VK_SELECT     = 41
	VK_PRINT      = 42
	VK_EXECUTE    = 43
	VK_SNAPSHOT   = 44
	VK_INSERT     = 45
	VK_DELETE     = 46
	VK_HELP       = 47
	VK_0          = 48
	VK_1          = 49
	VK_2          = 50
	VK_3          = 51
	VK_4          = 52
	VK_5          = 53
	VK_6          = 54
	VK_7          = 55
	VK_8          = 56
	VK_9          = 57
	VK_A          = 65
	VK_B          = 66
	VK_C          = 67
	VK_D          = 68
	VK_E          = 69
	VK_F          = 70
	VK_G          = 71
	VK_H          = 72
	VK_I          = 73
	VK_J          = 74
	VK_K          = 75
	VK_L          = 76
	VK_M          = 77
	VK_N          = 78
	VK_O          = 79
	VK_P          = 80
	VK_Q          = 81
	VK_R          = 82
	VK_S          = 83
	VK_T          = 84
	VK_U          = 85
	VK_V          = 86
	VK_W          = 87
	VK_X          = 88
	VK_Y          = 89
	VK_Z          = 90
	VK_LWIN       = 91
	VK_RWIN       = 92
	VK_APPS       = 93
	VK_SLEEP      = 95
	VK_NUMPAD0    = 96
	VK_NUMPAD1    = 97
	VK_NUMPAD2    = 98
	VK_NUMPAD3    = 99
	VK_NUMPAD4    = 100
	VK_NUMPAD5    = 101
	VK_NUMPAD6    = 102
	VK_NUMPAD7    = 103
	VK_NUMPAD8    = 104
	VK_NUMPAD9    = 105
	VK_MULTIPLY   = 106
	VK_ADD        = 107
	VK_SEPARATOR  = 108
	VK_SUBTRACT   = 109
	VK_DECIMAL    = 110
	VK_DIVIDE     = 111
	VK_F1         = 112
	VK_F2         = 113
	VK_F3         = 114
	VK_F4         = 115
	VK_F5         = 116
	VK_F6         = 117
	VK_F7         = 118
	VK_F8         = 119
	VK_F9         = 120
	VK_F10        = 121
	VK_F11        = 122
	VK_F12        = 123
	VK_F13        = 124
	VK_F14        = 125
	VK_F15        = 126
	VK_F16        = 127
	VK_F17        = 128
	VK_F18        = 129
	VK_F19        = 130
	VK_F20        = 131
	VK_F21        = 132
	VK_F22        = 133
	VK_F23        = 134
	VK_F24        = 135

	MM_MAX_NUMAXES = 16

	ETO_CLIPPED = 4

	SRCCOPY = 0x00CC0020

	TRANSPARENT = 1

	CLR_INVALID = 0xffffffff

	DIB_RGB_COLORS = 0

	BI_RGB = 0
)

type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     syscall.Handle
	HIcon         syscall.Handle
	HCursor       syscall.Handle
	HbrBackground syscall.Handle
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       syscall.Handle
}

func MAKELONG(a, b uint16) uint32 {
	return uint32(a) | uint32(b)<<16
}

func LOWORD(a uintptr) uint16 { return uint16(a & 0xffff) }
func HIWORD(a uintptr) uint16 { return uint16((a >> 16) & 0xffff) }

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type POINT struct {
	X int32
	Y int32
}

type SIZE struct {
	CX int32
	CY int32
}

type MSG struct {
	HWND    syscall.Handle
	Message uint32
	Wparam  uintptr
	Lparam  uintptr
	Time    uint32
	Pt      POINT
}

type PAINTSTRUCT struct {
	HDC       syscall.Handle
	Erase     int32
	Paint     RECT
	Restore   int32
	IncUpdate int32
	Reserved  [32]byte
}

type DESIGNVECTOR struct {
	Reserved uint32
	NumAxes  uint32
	Values   [MM_MAX_NUMAXES]int32
}

type TEXTMETRIC struct {
	Height           int32
	Ascent           int32
	Descent          int32
	InternalLeading  int32
	ExternalLeading  int32
	AveCharWidth     int32
	MaxCharWidth     int32
	Weight           int32
	Overhang         int32
	DigitizedAspectX int32
	DigitizedAspectY int32
	FirstChar        int16
	LastChar         int16
	DefaultChar      int16
	BreakChar        int16
	Italic           byte
	Underlined       byte
	StruckOut        byte
	PitchAndFamily   byte
	CharSet          byte
}

type COLORREF uint32

type BITMAPINFOHEADER struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint16
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type RGBQUAD struct {
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

type BITMAPINFO struct {
	Header BITMAPINFOHEADER
	Colors *RGBQUAD
}

//sys	GetDC(hwnd syscall.Handle) (dc syscall.Handle, err error) = user32.GetDC
//sys	ReleaseDC(hwnd syscall.Handle, dc syscall.Handle) (err error) = user32.ReleaseDC
//sys	SendMessage(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) = user32.SendMessageW
//sys	PostMessage(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (err error) = user32.PostMessageW
//sys	RegisterClassEx(wc *WNDCLASSEX) (atom uint16, err error) = user32.RegisterClassExW
//sys	LoadCursor(hInstance syscall.Handle, cursorName uintptr) (cursor syscall.Handle, err error) = user32.LoadCursorW
//sys	LoadIcon(hInstance syscall.Handle, iconName uintptr) (icon syscall.Handle, err error) = user32.LoadIconW
//sys	CreateWindowEx(exstyle uint32, className *uint16, windowText *uint16, style uint32, x int32, y int32, width int32, height int32, parent syscall.Handle, menu syscall.Handle, hInstance syscall.Handle, lpParam uintptr) (hwnd syscall.Handle, err error) = user32.CreateWindowExW
//sys	AdjustWindowRectEx(rect *RECT, style uint32, menu int32, exStyle uint32) (err error) = user32.AdjustWindowRectEx
//sys	DefWindowProc(hwnd syscall.Handle, uMsg uint32, wParam uintptr, lParam uintptr) (lResult uintptr) = user32.DefWindowProcW
//sys	DestroyWindow(hwnd syscall.Handle) (err error) = user32.DestroyWindow
//sys	SetCapture(hwnd syscall.Handle) (prev syscall.Handle) = user32.SetCapture
//sys	ReleaseCapture() (err error) = user32.ReleaseCapture
//sys	GetMessage(msg *MSG, hwnd syscall.Handle, msgfiltermin uint32, msgfiltermax uint32) (ret int32, err error) [failretval==-1] = user32.GetMessageW
//sys	PeekMessage(msg *MSG, hwnd syscall.Handle, msgfiltermin uint32, msgfiltermax uint32, removeflag uint32) (ok bool) = user32.PeekMessageW
//sys	TranslateMessage(msg *MSG) (ret int32) = user32.TranslateMessage
//sys	DispatchMessage(msg *MSG) (ret int32) = user32.DispatchMessageW
//sys	GetClientRect(hwnd syscall.Handle, rect *RECT) (err error) = user32.GetClientRect
//sys	GetDesktopWindow() (hwnd syscall.Handle) = user32.GetDesktopWindow
//sys	ShowWindow(hwnd syscall.Handle, cmdshow int32) (wasVisible bool) = user32.ShowWindow
//sys	IsWindowVisible(hwnd syscall.Handle) (visible bool) = user32.IsWindowVisible
//sys	FillRect(hdc syscall.Handle, rect *RECT, br syscall.Handle) (err error) = user32.FillRect
//sys	BeginPaint(hwnd syscall.Handle, paint *PAINTSTRUCT) (hdc syscall.Handle) = user32.BeginPaint
//sys	EndPaint(hwnd syscall.Handle, paint *PAINTSTRUCT) = user32.EndPaint

//sys	GetModuleHandle(moduleName *uint16) (module syscall.Handle, err error) = kernel32.GetModuleHandleW

//sys	CreateCompatibleDC(hdc syscall.Handle) (dc syscall.Handle, err error) = gdi32.CreateCompatibleDC
//sys	DeleteDC(hdc syscall.Handle) (err error) = gdi32.DeleteDC
//sys	CreateDIBSection(hdc syscall.Handle, bmi *BITMAPINFO, usage uint, p **byte, section syscall.Handle, offset uint32) (bitmap syscall.Handle, err error) = gdi32.CreateDIBSection
//sys	GetStockObject(fnObject int32) (gdiobj syscall.Handle, err error) = gdi32.GetStockObject
//sys	SelectObject(hdc syscall.Handle, gdiobj syscall.Handle) (prev syscall.Handle) = gdi32.SelectObject
//sys	DeleteObject(gdiobj syscall.Handle) (err error) = gdi32.DeleteObject
//sys	GetTextMetrics(hdc syscall.Handle, tm *TEXTMETRIC) (err error) = gdi32.GetTextMetricsW
//sys	GetTextExtentPoint32(hdc syscall.Handle, str *uint16, strlen int, size *SIZE) (err error) = gdi32.GetTextExtentPoint32W
//sys	AddFontMemResourceEx(font *byte, len uint32, pdv *DESIGNVECTOR, fonts *uint32) (handle syscall.Handle) = gdi32.AddFontMemResourceEx
//sys	RemoveFontMemResourceEx(handle syscall.Handle) (ok bool) = gdi32.RemoveFontMemResourceEx
//sys	CreateFont(height int32, width int32, escapement int32, orientation int32, weight int32, italic uint32, underline uint32, strikeOut uint32, charSet uint32, outputPrecision uint32, clipPrecision uint32, quality uint32, pitchAndFamily uint32, face *uint16) (font syscall.Handle, err error) = gdi32.CreateFontW
//sys	ExtTextOut(hdc syscall.Handle, x int, y int, options uint32, rect *RECT, str *uint16, strlen int, dx *int) (err error) = gdi32.ExtTextOutW
//sys	BitBlt(hdc syscall.Handle, dx int, dy int, w int, h int, src syscall.Handle, sx int, sy int, rop uint32) (err error) = gdi32.BitBlt
//sys	SetBkMode(hdc syscall.Handle, mode int) (prevMode int, err error) = gdi32.SetBkMode
//sys	SetTextColor(hdc syscall.Handle, col COLORREF) (prevCol COLORREF, err error) [failretval==CLR_INVALID] = gdi32.SetTextColor
//sys	CreateRectRgn(left int, top int, right int, bottom int) (hrgn syscall.Handle, err error) = gdi32.CreateRectRgn
//sys	SelectClipRgn(hdc syscall.Handle, hrgn syscall.Handle) (c int, err error) = gdi32.SelectClipRgn

//sys	DragAcceptFiles(hwnd syscall.Handle, accept bool) = shell32.DragAcceptFiles
//sys	DragQueryFile(drop syscall.Handle, file int, str *uint16, strlen int) (n int) = shell32.DragQueryFileW
//sys	DragFinish(drop syscall.Handle) = shell32.DragFinish

//sys	GdiplusStartup(token *uintptr, input *GdiplusStartupInput, output *GdiplusStartupOutput) (status GpStatus) = gdiplus.GdiplusStartup
//sys	GdiplusShutdown(token uintptr) = gdiplus.GdiplusShutdown

//sys	GdipCreateFromHDC(hdc syscall.Handle, g *GpGraphics) (status GpStatus) = gdiplus.GdipCreateFromHDC
//sys	GdipCreateFromHWND(hwnd syscall.Handle, g *GpGraphics) (status GpStatus) = gdiplus.GdipCreateFromHWND
//sys	GdipDeleteGraphics(g GpGraphics) (status GpStatus) = gdiplus.GdipDeleteGraphics
//sys	GdipGetDC(g GpGraphics, hdc *syscall.Handle) (status GpStatus) = gdiplus.GdipGetDC
//sys	GdipReleaseDC(g GpGraphics, hdc syscall.Handle) (status GpStatus) = gdiplus.GdipReleaseDC
//sys	GdipSetSmoothingMode(g GpGraphics, mode GpSmoothingMode) (status GpStatus) = gdiplus.GdipSetSmoothingMode
//sys	GdipGraphicsClear(g GpGraphics, col GpARGB) (status GpStatus) = gdiplus.GdipGraphicsClear
//sys	GdipSetClipRectI(g GpGraphics, x int, y int, w int, h int, cm GpCombineMode) (status GpStatus) = gdiplus.GdipSetClipRectI

//sys	GdipCreateBitmapFromGraphics(width int, height int, g GpGraphics, img *GpImage) (status GpStatus) = gdiplus.GdipCreateBitmapFromGraphics
//sys	GdipCreateBitmapFromScan0(width int, height int, stride int, format GpPixelFormat, scan0 uintptr, img *GpImage) (status GpStatus) = gdiplus.GdipCreateBitmapFromScan0
//sys	GdipDisposeImage(img GpImage) (status GpStatus) = gdiplus.GdipDisposeImage
//sys	GdipGetImageGraphicsContext(img GpImage, g *GpGraphics) (status GpStatus) = gdiplus.GdipGetImageGraphicsContext
//sys	GdipDrawImageI(g GpGraphics, img GpImage, x int, y int) (status GpStatus) = gdiplus.GdipDrawImageI
//sys	GdipDrawImageRectI(g GpGraphics, img GpImage, x int, y int, w int, h int) (status GpStatus) = gdiplus.GdipDrawImageRectI
//sys	GdipDrawImageRectRectI(g GpGraphics, img GpImage, dx int, dy int, dw int, dh int, sx int, sy int, sw int, sh int, unit GpUnit, ia GpImageAttributes, callback GpDrawImageAbort, cbdata unsafe.Pointer) (status GpStatus) = gdiplus.GdipDrawImageRectRectI

//sys	GdipCreateSolidFill(col GpARGB, br *GpBrush) (status GpStatus) = gdiplus.GdipCreateSolidFill
//sys	GdipSetSolidFillColor(br GpBrush, col GpARGB) (status GpStatus) = gdiplus.GdipSetSolidFillColor
//sys	GdipDeleteBrush(br GpBrush) (status GpStatus) = gdiplus.GdipDeleteBrush
//sys	GdipFillRectangleI(g GpGraphics, br GpBrush, x int, y int, width int, height int) (status GpStatus) = gdiplus.GdipFillRectangleI
//sys	GdipFillPieI(g GpGraphics, br GpBrush, x int, y int, width int, height int, startAngle float32, sweepAngle float32) (status GpStatus) = gdiplus.GdipFillPieI
//sys	GdipFillEllipseI(g GpGraphics, br GpBrush, x int, y int, width int, height int) (status GpStatus) = gdiplus.GdipFillEllipseI
//sys	GdipFillPolygonI(g GpGraphics, br GpBrush, points *GpPoint, plen int, m GpFillMode) (status GpStatus) = gdiplus.GdipFillPolygonI

//sys	GdipCreatePen1(argb GpARGB, width float32, unit GpUnit, pen *GpPen) (status GpStatus) = gdiplus.GdipCreatePen1
//sys	GdipDeletePen(pen GpPen) (status GpStatus) = gdiplus.GdipDeletePen
//sys	GdipSetPenWidth(pen GpPen, width float32) (status GpStatus) = gdiplus.GdipSetPenWidth
//sys	GdipSetPenColor(pen GpPen, argb GpARGB) (status GpStatus) = gdiplus.GdipSetPenColor
//sys	GdipDrawLineI(g GpGraphics, pen GpPen, x1 int, y1 int, x2 int, y2 int) (status GpStatus) = gdiplus.GdipDrawLineI
//sys	GdipDrawRectangleI(g GpGraphics, pen GpPen, x int, y int, width int, height int) (status GpStatus) = gdiplus.GdipDrawRectangleI
//sys	GdipDrawArcI(g GpGraphics, pen GpPen, x int, y int, width int, height int, startAngle float32, sweepAngle float32) (status GpStatus) = gdiplus.GdipDrawArcI
//sys	GdipDrawEllipseI(g GpGraphics, pen GpPen, x int, y int, width int, height int) (status GpStatus) = gdiplus.GdipDrawEllipseI
//sys	GdipDrawPolygonI(g GpGraphics, pen GpPen, points *GpPoint, plen int) (status GpStatus) = gdiplus.GdipDrawPolygonI
//sys	GdipDrawBezierI(g GpGraphics, pen GpPen, x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, x4 int, y4 int) (status GpStatus) = gdiplus.GdipDrawBezierI

type GdiplusStartupInput struct {
	GdiplusVersion           uint32
	DebugEventCallback       uintptr
	SuppressBackgroundThread int32
	SuppressExternalCodecs   int32
}

type GdiplusStartupOutput struct {
	NotificationHook   uintptr
	NotificationUnhook uintptr
}

type GpStatus uintptr

const (
	GpStatusOk                        = GpStatus(0)
	GpStatusGenericError              = GpStatus(1)
	GpStatusInvalidParameter          = GpStatus(2)
	GpStatusOutOfMemory               = GpStatus(3)
	GpStatusObjectBusy                = GpStatus(4)
	GpStatusInsufficientBuffer        = GpStatus(5)
	GpStatusNotImplemented            = GpStatus(6)
	GpStatusWin32Error                = GpStatus(7)
	GpStatusWrongState                = GpStatus(8)
	GpStatusAborted                   = GpStatus(9)
	GpStatusFileNotFound              = GpStatus(10)
	GpStatusValueOverflow             = GpStatus(11)
	GpStatusAccessDenied              = GpStatus(12)
	GpStatusUnknownImageFormat        = GpStatus(13)
	GpStatusFontFamilyNotFound        = GpStatus(14)
	GpStatusFontStyleNotFound         = GpStatus(15)
	GpStatusNotTrueTypeFont           = GpStatus(16)
	GpStatusUnsupportedGdiplusVersion = GpStatus(17)
	GpStatusGdiplusNotInitialized     = GpStatus(18)
	GpStatusPropertyNotFound          = GpStatus(19)
	GpStatusPropertyNotSupported      = GpStatus(20)
	GpStatusProfileNotFound           = GpStatus(21)
)

type GpUnit int

const (
	GpUnitWorld      = 0
	GpUnitDisplay    = 1
	GpUnitPixel      = 2
	GpUnitPoint      = 3
	GpUnitInch       = 4
	GpUnitDocument   = 5
	GpUnitMillimeter = 6
)

type GpSmoothingMode uintptr

const (
	GpQualityModeInvalid = -1
	GpQualityModeDefault = 0
	GpQualityModeLow     = 1
	GpQualityModeHigh    = 2

	GpSmoothingModeInvalid      = GpQualityModeInvalid
	GpSmoothingModeDefault      = GpQualityModeDefault
	GpSmoothingModeHighSpeed    = GpQualityModeLow
	GpSmoothingModeHighQuality  = GpQualityModeHigh
	GpSmoothingModeNone         = 3
	GpSmoothingModeAntiAlias    = 4
	GpSmoothingModeAntiAlias8x4 = GpSmoothingModeAntiAlias
	GpSmoothingModeAntiAlias8x8 = 5
)

type GpFillMode uintptr

const (
	GpFillModeAlternate = 0
	GpFillModeWinding   = 1
)

type GpCombineMode uintptr

const (
	GpCombineModeReplace    = GpCombineMode(0)
	GpCombineModeIntersect  = GpCombineMode(1)
	GpCombineModeUnion      = GpCombineMode(2)
	GpCombineModeXor        = GpCombineMode(3)
	GpCombineModeExclude    = GpCombineMode(4)
	GpCombineModeComplement = GpCombineMode(5)
)

type GpImageAttributes uintptr
type GpDrawImageAbort uintptr

type GpARGB uint32

type GpPoint struct {
	X int32
	Y int32
}

type GpRectF struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type GpBrush uintptr

func (b GpBrush) Close() error {
	return GdipDeleteBrush(b).Err()
}

type GpPen uintptr

func (p GpPen) Close() error {
	return GdipDeletePen(p).Err()
}

type GpGraphics uintptr

func (g GpGraphics) Close() error {
	return GdipDeleteGraphics(g).Err()
}

type GpPixelFormat int32

const (
	GpPixelFormatGDI        = GpPixelFormat(0x00020000) // Is a GDI-supported format
	GpPixelFormatAlpha      = GpPixelFormat(0x00040000) // Has an alpha component
	GpPixelFormatPAlpha     = GpPixelFormat(0x00080000) // Pre-multiplied alpha
	GpPixelFormatCanonical  = GpPixelFormat(0x00200000)
	GpPixelFormat32bppARGB  = GpPixelFormat(10 | (32 << 8) | GpPixelFormatAlpha | GpPixelFormatGDI | GpPixelFormatCanonical)
	GpPixelFormat32bppPARGB = GpPixelFormat(11 | (32 << 8) | GpPixelFormatAlpha | GpPixelFormatPAlpha | GpPixelFormatGDI)
)

type GpImage uintptr

func (i GpImage) Close() error {
	return GdipDisposeImage(i).Err()
}

func (s GpStatus) Err() error {
	switch s {
	case GpStatusOk:
		return nil
	case GpStatusGenericError:
		return errors.New("winapi(gdiplus): GenericError")
	case GpStatusInvalidParameter:
		return errors.New("winapi(gdiplus): InvalidParameter")
	case GpStatusOutOfMemory:
		return errors.New("winapi(gdiplus): OutOfMemory")
	case GpStatusObjectBusy:
		return errors.New("winapi(gdiplus): ObjectBusy")
	case GpStatusInsufficientBuffer:
		return errors.New("winapi(gdiplus): InsufficientBuffer")
	case GpStatusNotImplemented:
		return errors.New("winapi(gdiplus): NotImplemented")
	case GpStatusWin32Error:
		return errors.New("winapi(gdiplus): Win32Error")
	case GpStatusWrongState:
		return errors.New("winapi(gdiplus): WrongState")
	case GpStatusAborted:
		return errors.New("winapi(gdiplus): Aborted")
	case GpStatusFileNotFound:
		return errors.New("winapi(gdiplus): FileNotFound")
	case GpStatusValueOverflow:
		return errors.New("winapi(gdiplus): ValueOverflow")
	case GpStatusAccessDenied:
		return errors.New("winapi(gdiplus): AccessDenied")
	case GpStatusUnknownImageFormat:
		return errors.New("winapi(gdiplus): UnknownImageFormat")
	case GpStatusFontFamilyNotFound:
		return errors.New("winapi(gdiplus): FontFamilyNotFound")
	case GpStatusFontStyleNotFound:
		return errors.New("winapi(gdiplus): FontStyleNotFound")
	case GpStatusNotTrueTypeFont:
		return errors.New("winapi(gdiplus): NotTrueTypeFont")
	case GpStatusUnsupportedGdiplusVersion:
		return errors.New("winapi(gdiplus): UnsupportedGdiplusVersion")
	case GpStatusGdiplusNotInitialized:
		return errors.New("winapi(gdiplus): GdiplusNotInitialized")
	case GpStatusPropertyNotFound:
		return errors.New("winapi(gdiplus): PropertyNotFound")
	case GpStatusPropertyNotSupported:
		return errors.New("winapi(gdiplus): PropertyNotSupported")
	case GpStatusProfileNotFound:
		return errors.New("winapi(gdiplus): ProfileNotFound")
	}
	return errors.New("winapi(gdiplus): Unexpected error")
}
