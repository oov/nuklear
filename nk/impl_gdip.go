//go:build gdip
// +build gdip

package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT -Wno-implicit-function-declaration
#cgo windows LDFLAGS: -Wl,--allow-multiple-definition -lgdi32
#include <string.h>

#include "nuklear.h"

#define NK_IMPLEMENTATION
#define NK_GDI_IMPLEMENTATION
#include "nuklear.h"

// this code is part of stb_truetype.h - v1.17 - public domain
// authored from 2009-2016 by Sean Barrett / RAD Game Tools
NK_INTERN const char *nk_tt_GetFontNameString(const struct nk_tt_fontinfo *font, int *length, int platformID, int encodingID, int languageID, int nameID)
{
   nk_int i,count,stringOffset;
   const nk_byte *fc = font->data;
   nk_uint offset = font->fontstart;
   nk_uint nm = nk_tt__find_table(fc, offset, "name");
   if (!nm) return NULL;

   count = nk_ttUSHORT(fc+nm+2);
   stringOffset = nm + nk_ttUSHORT(fc+nm+4);
   for (i=0; i < count; ++i) {
      nk_uint loc = nm + 6 + 12 * i;
      if (platformID == nk_ttUSHORT(fc+loc+0) && encodingID == nk_ttUSHORT(fc+loc+2)
          && languageID == nk_ttUSHORT(fc+loc+4) && nameID == nk_ttUSHORT(fc+loc+6)) {
         *length = nk_ttUSHORT(fc+loc+8);
         return (const char *) (fc+stringOffset+nk_ttUSHORT(fc+loc+10));
      }
   }
   return NULL;
}

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

struct nk_gdip_font {
	HANDLE hdc;
	HANDLE hfont;
};

static float gdip_text_width(nk_handle h, float height, const char* str, int len) {
	struct nk_gdip_font *f = h.ptr;
	int wsz = MultiByteToWideChar(CP_UTF8, 0, str, len, NULL, 0);
    WCHAR* wstr = (WCHAR*)_alloca(wsz * sizeof(WCHAR));
	SIZE sz;
	MultiByteToWideChar(CP_UTF8, 0, str, len, wstr, wsz);
	if (GetTextExtentPoint32W(f->hdc, wstr, wsz, &sz)) {
		return (float)sz.cx;
	}
	return 0;
}

struct nk_user_font *create_gdip_font(HANDLE hdc, HANDLE hfont, float height) {
	struct nk_gdip_font *f = malloc(sizeof(struct nk_gdip_font));
	f->hdc = hdc;
	f->hfont = hfont;
	struct nk_user_font *uf = malloc(sizeof(struct nk_user_font));
	uf->userdata.ptr = f;
	uf->height = height;
	uf->width = &gdip_text_width;
	return uf;
}

void destroy_gdip_font(struct nk_user_font *uf) {
	struct nk_gdip_font *f = (struct nk_gdip_font*)uf->userdata.ptr;
	free(f);
	free(uf);
}

// Font chain for fallback support
#define NK_GDIP_FONT_CHAIN_MAGIC 0x4643484E  // "FCHN" - Font CHaiN

struct nk_gdip_font_chain {
	unsigned int magic;      // Magic number to identify font chain
	int count;               // Number of fonts in chain
	int capacity;            // Capacity of fonts array
	struct nk_gdip_font **fonts;  // Array of font pointers
};

// Check if a font has a glyph for a specific UTF-16 character
static int gdip_font_has_glyph(struct nk_gdip_font *f, WCHAR wch) {
	WORD gi;
	DWORD ret = GetGlyphIndicesW((HDC)f->hdc, &wch, 1, &gi, GGI_MARK_NONEXISTING_GLYPHS);
	if (ret == GDI_ERROR) return 1;  // Assume has glyph on error
	return gi != 0xFFFF;
}

// Find the appropriate font for a character in the chain
static int gdip_font_chain_find_font(struct nk_gdip_font_chain *fc, WCHAR wch) {
	for (int i = 0; i < fc->count; i++) {
		if (gdip_font_has_glyph(fc->fonts[i], wch)) {
			return i;
		}
	}
	return 0;  // Default to first font if none has the glyph
}

// Text width calculation for font chain
static float gdip_font_chain_text_width(nk_handle h, float height, const char* str, int len) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)h.ptr;
	if (fc->count == 0) return 0;

	int wsz = MultiByteToWideChar(CP_UTF8, 0, str, len, NULL, 0);
	if (wsz == 0) return 0;
	WCHAR* wstr = (WCHAR*)_alloca(wsz * sizeof(WCHAR));
	MultiByteToWideChar(CP_UTF8, 0, str, len, wstr, wsz);

	float total_width = 0;
	int seg_start = 0;
	int current_font = gdip_font_chain_find_font(fc, wstr[0]);

	for (int i = 1; i <= wsz; i++) {
		int font_idx = (i < wsz) ? gdip_font_chain_find_font(fc, wstr[i]) : -1;

		if (font_idx != current_font || i == wsz) {
			// End of segment, measure it
			SIZE sz;
			struct nk_gdip_font *f = fc->fonts[current_font];
			HGDIOBJ old = SelectObject((HDC)f->hdc, (HGDIOBJ)f->hfont);
			if (GetTextExtentPoint32W((HDC)f->hdc, wstr + seg_start, i - seg_start, &sz)) {
				total_width += (float)sz.cx;
			}
			SelectObject((HDC)f->hdc, old);

			seg_start = i;
			current_font = font_idx;
		}
	}

	return total_width;
}

struct nk_user_font *create_gdip_font_chain(float height) {
	struct nk_gdip_font_chain *fc = malloc(sizeof(struct nk_gdip_font_chain));
	fc->magic = NK_GDIP_FONT_CHAIN_MAGIC;
	fc->count = 0;
	fc->capacity = 4;
	fc->fonts = malloc(sizeof(struct nk_gdip_font*) * fc->capacity);

	struct nk_user_font *uf = malloc(sizeof(struct nk_user_font));
	uf->userdata.ptr = fc;
	uf->height = height;
	uf->width = &gdip_font_chain_text_width;
	return uf;
}

void gdip_font_chain_add(struct nk_user_font *uf, HANDLE hdc, HANDLE hfont) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)uf->userdata.ptr;

	if (fc->count >= fc->capacity) {
		fc->capacity *= 2;
		fc->fonts = realloc(fc->fonts, sizeof(struct nk_gdip_font*) * fc->capacity);
	}

	struct nk_gdip_font *f = malloc(sizeof(struct nk_gdip_font));
	f->hdc = hdc;
	f->hfont = hfont;
	fc->fonts[fc->count++] = f;
}

void destroy_gdip_font_chain(struct nk_user_font *uf) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)uf->userdata.ptr;
	for (int i = 0; i < fc->count; i++) {
		free(fc->fonts[i]);
	}
	free(fc->fonts);
	free(fc);
	free(uf);
}

int is_gdip_font_chain(struct nk_user_font *uf) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)uf->userdata.ptr;
	return fc->magic == NK_GDIP_FONT_CHAIN_MAGIC;
}

int gdip_font_chain_count(struct nk_user_font *uf) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)uf->userdata.ptr;
	return fc->count;
}

struct nk_gdip_font *gdip_font_chain_get(struct nk_user_font *uf, int index) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)uf->userdata.ptr;
	if (index < 0 || index >= fc->count) return NULL;
	return fc->fonts[index];
}

// Segment information for font chain text rendering
struct nk_gdip_text_segment {
	int start;      // Start index in UTF-16 string
	int length;     // Length in UTF-16 characters
	int font_index; // Index of font to use
	int x_offset;   // X offset from start of text
};

// Get segments for rendering text with font chain
// Returns number of segments, fills segments array (caller must provide enough space)
// max_segments should be at least equal to the UTF-16 string length
int gdip_font_chain_get_segments(struct nk_user_font *uf, const WCHAR *wstr, int wsz,
                                  struct nk_gdip_text_segment *segments, int max_segments) {
	struct nk_gdip_font_chain *fc = (struct nk_gdip_font_chain*)uf->userdata.ptr;
	if (fc->count == 0 || wsz == 0) return 0;

	int seg_count = 0;
	int seg_start = 0;
	int current_font = gdip_font_chain_find_font(fc, wstr[0]);
	int x_offset = 0;

	for (int i = 1; i <= wsz; i++) {
		int font_idx = (i < wsz) ? gdip_font_chain_find_font(fc, wstr[i]) : -1;

		if (font_idx != current_font || i == wsz) {
			// End of segment
			if (seg_count < max_segments) {
				segments[seg_count].start = seg_start;
				segments[seg_count].length = i - seg_start;
				segments[seg_count].font_index = current_font;
				segments[seg_count].x_offset = x_offset;

				// Calculate width for x_offset of next segment
				SIZE sz;
				struct nk_gdip_font *f = fc->fonts[current_font];
				HGDIOBJ old = SelectObject((HDC)f->hdc, (HGDIOBJ)f->hfont);
				if (GetTextExtentPoint32W((HDC)f->hdc, wstr + seg_start, i - seg_start, &sz)) {
					x_offset += sz.cx;
				}
				SelectObject((HDC)f->hdc, old);

				seg_count++;
			}

			seg_start = i;
			current_font = font_idx;
		}
	}

	return seg_count;
}

*/
import "C"
import (
	"errors"
	"image"
	"syscall"
	"unsafe"

	"github.com/golang-ui/nuklear/nk/internal"
	"github.com/golang-ui/nuklear/nk/internal/winapi"
	"github.com/golang-ui/nuklear/nk/w32"
)

type PlatformInitOption int

const (
	PlatformDefault PlatformInitOption = iota
	PlatformInstallCallbacks
)

var gpToken uintptr

func gpMust(s winapi.GpStatus) {
	if err := s.Err(); err != nil {
		panic(err)
	}
}

func NkPlatformInit(win *w32.Window, opt PlatformInitOption) *Context {
	gpMust(winapi.GdiplusStartup(&gpToken, &winapi.GdiplusStartupInput{
		GdiplusVersion:         1,
		SuppressExternalCodecs: 1,
	}, nil))

	state.win = win
	if opt == PlatformInstallCallbacks {
		win.SetSizeCallback(func(_ *w32.Window, w, h int) {
			if w == 0 || h == 0 {
				if state.bmp != nil {
					state.bmp.Close()
					state.bmp = nil
				}
			} else if state.bmp == nil || state.bmp.Width != w || state.bmp.Height != h {
				bmp, err := internal.NewBitmap(w, h)
				if err != nil {
					panic(err)
				}
				if state.bmp != nil {
					state.bmp.Close()
				}
				state.bmp = bmp
			}
		})
		win.SetPaintCallback(func(_ *w32.Window, hdc syscall.Handle) {
			if state.bmp == nil {
				return
			}
			if err := state.bmp.UseHDC(func(bmphdc syscall.Handle) error {
				return winapi.BitBlt(hdc, 0, 0, state.bmp.Width, state.bmp.Height, bmphdc, 0, 0, winapi.SRCCOPY)
			}); err != nil {
				panic(err)
			}
		})
		win.SetMouseWheelYCallback(func(_ *w32.Window, v float32) {
			state.wheelY += v
		})
	}

	hdc, err := winapi.GetDC(win.Handle)
	if err != nil {
		panic(err)
	}

	w, h := win.GetSize()
	bmp, err := internal.NewBitmap(w, h)
	if err != nil {
		panic(err)
	}

	var pen winapi.GpPen
	gpMust(winapi.GdipCreatePen1(0, 1, winapi.GpUnitPixel, &pen))
	var brush winapi.GpBrush
	gpMust(winapi.GdipCreateSolidFill(0, &brush))

	state.hdc = hdc
	state.pen = pen
	state.brush = brush
	state.bmp = bmp
	state.ctx = NewContext()
	NkInitDefault(state.ctx, nil)
	return state.ctx
}

func NkPlatformShutdown() {
	if state.hrgn != 0 {
		winapi.DeleteObject(state.hrgn)
	}
	state.pen.Close()
	state.brush.Close()
	state.bmp.Close()
	winapi.ReleaseDC(state.win.Handle, state.hdc)

	NkFree(state.ctx)
	state = nil

	winapi.GdiplusShutdown(gpToken)
}

type GdipFont struct {
	pf *internal.PrivateFont
	f  *internal.Font
	uf *C.struct_nk_user_font
}

func (f *GdipFont) Close() error {
	if f.f != nil {
		f.f.Close()
	}
	if f.pf != nil {
		f.pf.Close()
	}
	if f.uf != nil {
		C.destroy_gdip_font(f.uf)
	}
	return nil
}

func (f *GdipFont) Handle() *UserFont {
	return NewUserFontRef(unsafe.Pointer(f.uf))
}

func getFontNameString(f *C.struct_nk_tt_fontinfo, platformID, encodingID, languageID, nameID C.int) (string, bool) {
	var l C.int
	p := C.nk_tt_GetFontNameString(f, &l, platformID, encodingID, languageID, nameID)
	if p == nil {
		return "", false
	}
	sp := (*[^uint32(0) >> 1]uint8)(unsafe.Pointer(p))
	s := make([]uint16, l/2)
	for i := range s {
		s[i] = uint16(sp[i*2])<<8 | uint16(sp[i*2+1])
	}
	return syscall.UTF16ToString(s), true
}

func NkCreateFontFromBytes(data []byte, height int) (*GdipFont, error) {
	font := (*C.struct_nk_tt_fontinfo)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_nk_tt_fontinfo{}))))
	defer C.free(unsafe.Pointer(font))
	C.nk_tt_InitFont(font, (*C.uchar)(unsafe.Pointer(&data[0])), 0)
	const (
		// https://www.microsoft.com/typography/otspec/name.htm#nameIDs
		NID_FULL_FONT_NAME = 4
	)
	s, ok := getFontNameString(font, C.NK_TT_PLATFORM_ID_MICROSOFT, C.NK_TT_MS_EID_UNICODE_BMP, C.NK_TT_MS_LANG_ENGLISH, NID_FULL_FONT_NAME)
	if !ok {
		return nil, errors.New("nk: cannot find valid english font face name")
	}

	pf, err := internal.NewPrivateFont(s, data)
	if err != nil {
		return nil, err
	}
	f, err := pf.NewFont(height)
	if err != nil {
		pf.Close()
		return nil, err
	}
	gf := GdipFont{
		pf: pf,
		f:  f,
		uf: C.create_gdip_font(C.HANDLE(unsafe.Pointer(f.DC)), C.HANDLE(unsafe.Pointer(f.Handle)), C.float(float32(f.Height))),
	}
	return &gf, nil
}

// fontChainEntry represents an entry in the font chain.
// It can be either a GdipFont (from TTF bytes) or a system font.
type fontChainEntry struct {
	gdipFont   *GdipFont      // non-nil if created from bytes
	systemFont *internal.Font // non-nil if created from system font name
}

func (e *fontChainEntry) Close() error {
	if e.gdipFont != nil {
		return e.gdipFont.Close()
	}
	if e.systemFont != nil {
		return e.systemFont.Close()
	}
	return nil
}

// FontChain provides font fallback support for multi-language text rendering.
// When rendering text, if the primary font doesn't contain a glyph for a character,
// the chain will try subsequent fonts in order until one with the glyph is found.
type FontChain struct {
	entries []*fontChainEntry
	height  int
	uf      *C.struct_nk_user_font
}

// NewFontChain creates a new font chain with the specified height.
// The primary font should be added first using AddFont or AddSystemFont.
func NewFontChain(height int) *FontChain {
	return &FontChain{
		entries: make([]*fontChainEntry, 0, 4),
		height:  height,
		uf:      C.create_gdip_font_chain(C.float(float32(height))),
	}
}

// AddFont adds a font to the chain. Fonts are tried in the order they are added.
// The first font added is the primary font.
// The FontChain takes ownership of the font and will close it when the chain is closed.
func (fc *FontChain) AddFont(font *GdipFont) {
	fc.entries = append(fc.entries, &fontChainEntry{gdipFont: font})
	C.gdip_font_chain_add(fc.uf, C.HANDLE(unsafe.Pointer(font.f.DC)), C.HANDLE(unsafe.Pointer(font.f.Handle)))
}

// AddFontFromBytes creates a font from TTF data and adds it to the chain.
func (fc *FontChain) AddFontFromBytes(data []byte) error {
	font, err := NkCreateFontFromBytes(data, fc.height)
	if err != nil {
		return err
	}
	fc.AddFont(font)
	return nil
}

// AddSystemFont adds a system font by name to the chain.
// Example names: "MS UI Gothic", "Arial", "Segoe UI", "Yu Gothic UI"
func (fc *FontChain) AddSystemFont(name string) error {
	f, err := internal.NewFont(name, fc.height)
	if err != nil {
		return err
	}
	fc.entries = append(fc.entries, &fontChainEntry{systemFont: f})
	C.gdip_font_chain_add(fc.uf, C.HANDLE(unsafe.Pointer(f.DC)), C.HANDLE(unsafe.Pointer(f.Handle)))
	return nil
}

// Handle returns the nuklear user font handle for this font chain.
func (fc *FontChain) Handle() *UserFont {
	return NewUserFontRef(unsafe.Pointer(fc.uf))
}

// Close releases all resources associated with the font chain.
// This also closes all fonts that were added to the chain.
func (fc *FontChain) Close() error {
	for _, e := range fc.entries {
		e.Close()
	}
	fc.entries = nil
	if fc.uf != nil {
		C.destroy_gdip_font_chain(fc.uf)
		fc.uf = nil
	}
	return nil
}

// FontCount returns the number of fonts in the chain.
func (fc *FontChain) FontCount() int {
	return len(fc.entries)
}

type GdipImage struct {
	Width  int
	Height int
	img    winapi.GpImage
	p      unsafe.Pointer
	mem    unsafe.Pointer
}

func (img *GdipImage) Handle() Handle {
	return NkHandlePtr(img.p)
}

func (img *GdipImage) Close() error {
	if img.img != 0 {
		img.img.Close()
	}
	if img.p != nil {
		C.free(img.p)
	}
	if img.mem != nil {
		C.free(img.mem)
	}
	return nil
}

func NkCreateImage(img image.Image) (*GdipImage, error) {
	var i winapi.GpImage
	var err error
	var mem unsafe.Pointer
	switch t := img.(type) {
	case *image.NRGBA:
		mem = C.malloc(C.size_t(len(t.Pix)))
		s := t.Pix
		d := (*[^uint32(0) >> 1]byte)(mem)[:len(t.Pix)]
		for i := 0; i < len(t.Pix); i += 4 {
			d[i+0] = s[i+2]
			d[i+1] = s[i+1]
			d[i+2] = s[i+0]
			d[i+3] = s[i+3]
		}
		err = winapi.GdipCreateBitmapFromScan0(t.Rect.Dx(), t.Rect.Dy(), t.Rect.Dx()*4, winapi.GpPixelFormat32bppARGB, uintptr(mem), &i).Err()
	case *image.RGBA:
		mem = C.malloc(C.size_t(len(t.Pix)))
		s := t.Pix
		d := (*[^uint32(0) >> 1]byte)(mem)[:len(t.Pix)]
		for i := 0; i < len(t.Pix); i += 4 {
			if a := uint32(s[i+3]); a == 0 {
				d[i+0] = 0
				d[i+1] = 0
				d[i+2] = 0
				d[i+3] = 0
			} else if a == 255 {
				d[i+0] = s[i+2]
				d[i+1] = s[i+1]
				d[i+2] = s[i+0]
				d[i+3] = s[i+3]
			} else {
				d[i+0] = byte(uint32(s[i+2]) * 255 / a)
				d[i+1] = byte(uint32(s[i+1]) * 255 / a)
				d[i+2] = byte(uint32(s[i+0]) * 255 / a)
				d[i+3] = s[i+3]
			}
		}
		err = winapi.GdipCreateBitmapFromScan0(t.Rect.Dx(), t.Rect.Dy(), t.Rect.Dx()*4, winapi.GpPixelFormat32bppARGB, uintptr(mem), &i).Err()
	default:
		return nil, errors.New("unsupported image format")
	}
	if err != nil {
		return nil, err
	}
	rect := img.Bounds()
	r := &GdipImage{
		Width:  rect.Dx(),
		Height: rect.Dy(),
		img:    i,
		mem:    mem,
	}
	r.p = C.malloc(C.size_t(unsafe.Sizeof(winapi.GpImage(0))))
	*(*winapi.GpImage)(r.p) = i
	return r, nil
}

func keysPressed(win *w32.Window, keys ...int) int32 {
	for _, k := range keys {
		if _, ok := win.Keys[k]; ok {
			return 1
		}
	}
	return 0
}

func b2i(v bool) int32 {
	if v {
		return 1
	}
	return 0
}

func convColor(c *C.struct_nk_color) winapi.GpARGB {
	return winapi.GpARGB(uint32(c.a)<<24 | uint32(c.r)<<16 | uint32(c.g)<<8 | uint32(c.b))
}
func convColorref(c *C.struct_nk_color) winapi.COLORREF {
	return winapi.COLORREF(uint32(c.b)<<16 | uint32(c.g)<<8 | uint32(c.r))
}

func NkPlatformNewFrame() {
	win := state.win
	ctx := state.ctx
	w, h := win.GetSize()
	if w == 0 || h == 0 {
		if state.bmp != nil {
			state.bmp.Close()
			state.bmp = nil
		}
	} else if state.bmp == nil || state.bmp.Width != w || state.bmp.Height != h {
		bmp, err := internal.NewBitmap(w, h)
		if err != nil {
			panic(err)
		}
		if state.bmp != nil {
			state.bmp.Close()
		}
		state.bmp = bmp
	}

	NkInputBegin(ctx)
	NkInputKey(ctx, KeyDel, keysPressed(win, winapi.VK_DELETE))
	NkInputKey(ctx, KeyEnter, keysPressed(win, winapi.VK_RETURN))
	NkInputKey(ctx, KeyTab, keysPressed(win, winapi.VK_TAB))
	NkInputKey(ctx, KeyBackspace, keysPressed(win, winapi.VK_BACK))
	NkInputKey(ctx, KeyUp, keysPressed(win, winapi.VK_UP))
	NkInputKey(ctx, KeyDown, keysPressed(win, winapi.VK_DOWN))
	NkInputKey(ctx, KeyTextStart, keysPressed(win, winapi.VK_HOME))
	NkInputKey(ctx, KeyTextEnd, keysPressed(win, winapi.VK_END))
	NkInputKey(ctx, KeyScrollStart, keysPressed(win, winapi.VK_HOME))
	NkInputKey(ctx, KeyScrollEnd, keysPressed(win, winapi.VK_END))
	NkInputKey(ctx, KeyScrollDown, keysPressed(win, winapi.VK_NEXT))
	NkInputKey(ctx, KeyScrollUp, keysPressed(win, winapi.VK_PRIOR))
	NkInputKey(ctx, KeyShift, keysPressed(win, winapi.VK_SHIFT))
	ctrl := keysPressed(win, winapi.VK_CONTROL)
	NkInputKey(ctx, KeyCtrl, ctrl)
	if ctrl != 0 {
		NkInputKey(ctx, KeyLeft, 0)
		NkInputKey(ctx, KeyRight, 0)
		NkInputKey(ctx, KeyCopy, keysPressed(win, winapi.VK_C))
		NkInputKey(ctx, KeyPaste, keysPressed(win, winapi.VK_V))
		NkInputKey(ctx, KeyCut, keysPressed(win, winapi.VK_X))
		NkInputKey(ctx, KeyTextUndo, keysPressed(win, winapi.VK_Z))
		NkInputKey(ctx, KeyTextRedo, keysPressed(win, winapi.VK_R))
		NkInputKey(ctx, KeyTextWordLeft, keysPressed(win, winapi.VK_LEFT))
		NkInputKey(ctx, KeyTextWordRight, keysPressed(win, winapi.VK_RIGHT))
		NkInputKey(ctx, KeyTextLineStart, keysPressed(win, winapi.VK_B))
		NkInputKey(ctx, KeyTextLineEnd, keysPressed(win, winapi.VK_E))
	} else {
		NkInputKey(ctx, KeyLeft, keysPressed(win, winapi.VK_LEFT))
		NkInputKey(ctx, KeyRight, keysPressed(win, winapi.VK_RIGHT))
		NkInputKey(ctx, KeyCopy, 0)
		NkInputKey(ctx, KeyPaste, 0)
		NkInputKey(ctx, KeyCut, 0)
		NkInputKey(ctx, KeyTextUndo, 0)
		NkInputKey(ctx, KeyTextRedo, 0)
		NkInputKey(ctx, KeyTextWordLeft, 0)
		NkInputKey(ctx, KeyTextWordRight, 0)
		NkInputKey(ctx, KeyTextLineStart, 0)
		NkInputKey(ctx, KeyTextLineEnd, 0)
	}
	x, y := int32(win.Mouse.X), int32(win.Mouse.Y)
	NkInputMotion(ctx, x, y)
	if win.Mouse.Button&w32.MouseButtonLeftModified == 0 {
		NkInputButton(ctx, ButtonLeft, x, y, b2i(state.mouseL))
	} else {
		isDown := win.Mouse.Button&w32.MouseButtonLeft != 0
		NkInputButton(ctx, ButtonLeft, x, y, b2i(isDown))
		if state.mouseL == isDown {
			NkInputButton(ctx, ButtonLeft, x, y, b2i(isDown))
		}
		state.mouseL = isDown
	}
	if win.Mouse.Button&w32.MouseButtonRightModified == 0 {
		NkInputButton(ctx, ButtonRight, x, y, b2i(state.mouseR))
	} else {
		isDown := win.Mouse.Button&w32.MouseButtonRight != 0
		NkInputButton(ctx, ButtonRight, x, y, b2i(isDown))
		if state.mouseL == isDown {
			NkInputButton(ctx, ButtonRight, x, y, b2i(isDown))
		}
		state.mouseR = isDown
	}
	if win.Mouse.Button&w32.MouseButtonMiddleModified == 0 {
		NkInputButton(ctx, ButtonMiddle, x, y, b2i(state.mouseM))
	} else {
		isDown := win.Mouse.Button&w32.MouseButtonMiddle != 0
		NkInputButton(ctx, ButtonMiddle, x, y, b2i(isDown))
		if state.mouseL == isDown {
			NkInputButton(ctx, ButtonMiddle, x, y, b2i(isDown))
		}
		state.mouseM = isDown
	}
	win.Mouse.Button &^= w32.MouseButtonLeftModified | w32.MouseButtonRightModified | w32.MouseButtonMiddleModified
	// TODO: DOUBLE CLICK
	NkInputScroll(ctx, NkVec2(0, state.wheelY))
	state.wheelY = 0
	for _, r := range win.Chars {
		NkInputUnicode(ctx, Rune(r))
	}
	win.Chars = win.Chars[:0]
	NkInputEnd(ctx)
}

func NkPlatformRender(aa AntiAliasing, clearColor Color) {
	ctx := state.ctx
	if state.bmp == nil || !state.win.Visible() {
		NkClear(state.ctx)
		return
	}
	bmp := state.bmp

	if aa == AntiAliasingOn {
		gpMust(winapi.GdipSetSmoothingMode(bmp.Graphics, winapi.GpSmoothingModeHighQuality))
	} else {
		gpMust(winapi.GdipSetSmoothingMode(bmp.Graphics, winapi.GpSmoothingModeNone))
	}
	gpMust(winapi.GdipGraphicsClear(bmp.Graphics, convColor(clearColor.Ref())))

	for cmd := Nk_Begin(ctx); cmd != nil; cmd = Nk_Next(ctx, cmd) {
		switch cmd._type {
		case CommandTypeNop:
		case CommandTypeScissor:
			p := (*CommandScissor)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetClipRectI(bmp.Graphics, int(p.x), int(p.y), int(p.w+1), int(p.h+1), winapi.GpCombineModeReplace))
			hrgn, err := winapi.CreateRectRgn(int(p.x), int(p.y), int(p.x)+int(p.w)+1, int(p.y)+int(p.h)+1)
			if err != nil {
				panic(err)
			}
			if state.hrgn != 0 {
				winapi.DeleteObject(state.hrgn)
			}
			state.hrgn = hrgn

		case CommandTypeLine:
			p := (*CommandLine)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
			gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
			gpMust(winapi.GdipDrawLineI(bmp.Graphics, state.pen, int(p.begin.x), int(p.begin.y), int(p.end.x), int(p.end.y)))
		case CommandTypeRect:
			p := (*CommandRect)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
			gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
			x, y, w, h, r := int(p.x), int(p.y), int(p.w), int(p.h), int(p.rounding)
			if r == 0 {
				gpMust(winapi.GdipDrawRectangleI(bmp.Graphics, state.pen, x, y, w, h))
			} else {
				d := r * 2
				gpMust(winapi.GdipDrawArcI(bmp.Graphics, state.pen, x, y, d, d, 180, 90))
				gpMust(winapi.GdipDrawLineI(bmp.Graphics, state.pen, x+r, y, x+w-r, y))
				gpMust(winapi.GdipDrawArcI(bmp.Graphics, state.pen, x+w-d, y, d, d, 270, 90))
				gpMust(winapi.GdipDrawLineI(bmp.Graphics, state.pen, x+w, y+r, x+w, y+h-r))
				gpMust(winapi.GdipDrawArcI(bmp.Graphics, state.pen, x+w-d, y+h-d, d, d, 0, 90))
				gpMust(winapi.GdipDrawLineI(bmp.Graphics, state.pen, x, y+r, x, y+h-r))
				gpMust(winapi.GdipDrawArcI(bmp.Graphics, state.pen, x, y+h-d, d, d, 90, 90))
				gpMust(winapi.GdipDrawLineI(bmp.Graphics, state.pen, x+r, y+h, x+w-r, y+h))
			}
		case CommandTypeRectFilled:
			p := (*CommandRectFilled)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetSolidFillColor(state.brush, convColor(&p.color)))
			x, y, w, h, r := int(p.x), int(p.y), int(p.w), int(p.h), int(p.rounding)
			if r == 0 {
				gpMust(winapi.GdipFillRectangleI(bmp.Graphics, state.brush, x, y, w, h))
			} else {
				d := r * 2
				gpMust(winapi.GdipFillRectangleI(bmp.Graphics, state.brush, x+r-1, y, w-d+2, h))
				gpMust(winapi.GdipFillRectangleI(bmp.Graphics, state.brush, x, y+r-1, w, h-d+2))
				gpMust(winapi.GdipFillPieI(bmp.Graphics, state.brush, x, y, d, d, 180, 90))
				gpMust(winapi.GdipFillPieI(bmp.Graphics, state.brush, x+w-d, y, d, d, 270, 90))
				gpMust(winapi.GdipFillPieI(bmp.Graphics, state.brush, x+w-d, y+h-d, d, d, 0, 90))
				gpMust(winapi.GdipFillPieI(bmp.Graphics, state.brush, x, y+h-d, d, d, 90, 90))
			}
		case CommandTypeCircle:
			p := (*CommandCircle)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
			gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
			gpMust(winapi.GdipDrawEllipseI(bmp.Graphics, state.pen, int(p.x), int(p.y), int(p.w), int(p.h)))
		case CommandTypeCircleFilled:
			p := (*CommandCircleFilled)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetSolidFillColor(state.brush, convColor(&p.color)))
			gpMust(winapi.GdipFillEllipseI(bmp.Graphics, state.brush, int(p.x), int(p.y), int(p.w), int(p.h)))
		case CommandTypeTriangle:
			p := (*CommandTriangle)(unsafe.Pointer(cmd))
			points := [4]winapi.GpPoint{
				{int32(p.a.x), int32(p.a.y)},
				{int32(p.b.x), int32(p.b.y)},
				{int32(p.c.x), int32(p.c.y)},
				{int32(p.a.x), int32(p.a.y)},
			}
			gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
			gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
			gpMust(winapi.GdipDrawPolygonI(bmp.Graphics, state.pen, &points[0], len(points)))
		case CommandTypeTriangleFilled:
			p := (*CommandTriangleFilled)(unsafe.Pointer(cmd))
			points := [3]winapi.GpPoint{
				{int32(p.a.x), int32(p.a.y)},
				{int32(p.b.x), int32(p.b.y)},
				{int32(p.c.x), int32(p.c.y)},
			}
			gpMust(winapi.GdipSetSolidFillColor(state.brush, convColor(&p.color)))
			gpMust(winapi.GdipFillPolygonI(bmp.Graphics, state.brush, &points[0], len(points), winapi.GpFillModeAlternate))
		case CommandTypePolygon:
			p := (*CommandPolygon)(unsafe.Pointer(cmd))
			if p.point_count <= 0 {
				continue
			}
			vectors := (*[^uint32(0) >> 1]C.struct_nk_vec2i)(unsafe.Pointer(&p.points[0]))
			points := make([]winapi.GpPoint, p.point_count, p.point_count+1)
			for i := range points {
				points[i].X = int32(vectors[i].x)
				points[i].Y = int32(vectors[i].y)
			}
			points = append(points, points[0])
			gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
			gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
			gpMust(winapi.GdipDrawPolygonI(bmp.Graphics, state.pen, &points[0], len(points)))
		case CommandTypePolygonFilled:
			p := (*CommandPolygon)(unsafe.Pointer(cmd))
			if p.point_count <= 0 {
				continue
			}
			vectors := (*[^uint32(0) >> 1]C.struct_nk_vec2i)(unsafe.Pointer(&p.points[0]))
			points := make([]winapi.GpPoint, p.point_count)
			for i := range points {
				points[i].X = int32(vectors[i].x)
				points[i].Y = int32(vectors[i].y)
			}
			gpMust(winapi.GdipSetSolidFillColor(state.brush, convColor(&p.color)))
			gpMust(winapi.GdipFillPolygonI(bmp.Graphics, state.brush, &points[0], len(points), winapi.GpFillModeAlternate))
		case CommandTypePolyline:
			p := (*CommandPolyline)(unsafe.Pointer(cmd))
			if p.point_count > 0 {
				vectors := (*[^uint32(0) >> 1]C.struct_nk_vec2i)(unsafe.Pointer(&p.points[0]))
				points := make([]winapi.GpPoint, p.point_count)
				for i := range points {
					points[i].X = int32(vectors[i].x)
					points[i].Y = int32(vectors[i].y)
				}
				gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
				gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
				gpMust(winapi.GdipDrawPolygonI(bmp.Graphics, state.pen, &points[0], len(points)))
			}
		case CommandTypeText:
			p := (*CommandText)(unsafe.Pointer(cmd))
			s, err := syscall.UTF16FromString(string((*[^uint32(0) >> 1]byte)(unsafe.Pointer(&p.string[0]))[:p.length]))
			if err != nil {
				panic(err)
			}
			uf := (*C.struct_nk_user_font)(unsafe.Pointer(p.font))
			isFontChain := C.is_gdip_font_chain(uf) != 0

			if isFontChain {
				// Font chain: render with segmented approach
				if err = bmp.UseHDC(func(hdc syscall.Handle) error {
					if state.hrgn != 0 {
						winapi.SelectClipRgn(hdc, state.hrgn)
					}
					oldMode, err := winapi.SetBkMode(hdc, winapi.TRANSPARENT)
					if err != nil {
						return err
					}
					oldCol, err := winapi.SetTextColor(hdc, convColorref(&p.foreground))
					if err != nil {
						return err
					}

					// Get segments for rendering
					wsz := len(s) - 1 // Exclude null terminator
					if wsz > 0 {
						segments := make([]C.struct_nk_gdip_text_segment, wsz)
						segCount := C.gdip_font_chain_get_segments(uf, (*C.WCHAR)(unsafe.Pointer(&s[0])), C.int(wsz), &segments[0], C.int(wsz))

						for i := 0; i < int(segCount); i++ {
							seg := segments[i]
							font := C.gdip_font_chain_get(uf, seg.font_index)
							if font == nil {
								continue
							}

							old := winapi.SelectObject(hdc, syscall.Handle(unsafe.Pointer(font.hfont)))
							winapi.ExtTextOut(hdc, int(p.x)+int(seg.x_offset), int(p.y), 0, nil, &s[seg.start], int(seg.length), nil)
							winapi.SelectObject(hdc, old)
						}
					}

					winapi.SetBkMode(hdc, oldMode)
					winapi.SetTextColor(hdc, oldCol)
					return nil
				}); err != nil {
					panic(err)
				}
			} else {
				// Single font: original rendering
				font := *(**C.struct_nk_gdip_font)(unsafe.Pointer(&p.font.userdata))
				if err = bmp.UseHDC(func(hdc syscall.Handle) error {
					if state.hrgn != 0 {
						winapi.SelectClipRgn(hdc, state.hrgn)
					}
					old := winapi.SelectObject(hdc, syscall.Handle(unsafe.Pointer(font.hfont)))
					oldMode, err := winapi.SetBkMode(hdc, winapi.TRANSPARENT)
					if err != nil {
						return err
					}
					oldCol, err := winapi.SetTextColor(hdc, convColorref(&p.foreground))
					if err != nil {
						return err
					}

					err = winapi.ExtTextOut(hdc, int(p.x), int(p.y), 0, nil, &s[0], len(s)-1, nil)
					winapi.SelectObject(hdc, old)
					winapi.SetBkMode(hdc, oldMode)
					winapi.SetTextColor(hdc, oldCol)
					return err
				}); err != nil {
					panic(err)
				}
			}
		case CommandTypeCurve:
			p := (*CommandCurve)(unsafe.Pointer(cmd))
			gpMust(winapi.GdipSetPenWidth(state.pen, float32(p.line_thickness)))
			gpMust(winapi.GdipSetPenColor(state.pen, convColor(&p.color)))
			gpMust(winapi.GdipDrawBezierI(bmp.Graphics, state.pen,
				int(p.begin.x), int(p.begin.y),
				int(p.ctrl[0].x), int(p.ctrl[0].y),
				int(p.ctrl[1].x), int(p.ctrl[1].y),
				int(p.end.x), int(p.end.y),
			))
		case CommandTypeRectMultiColor:
		case CommandTypeImage:
			p := (*CommandImage)(unsafe.Pointer(cmd))
			img := **(**winapi.GpImage)(unsafe.Pointer(&p.img.handle))
			if p.img.w == 0 && p.img.h == 0 {
				gpMust(winapi.GdipDrawImageRectI(
					bmp.Graphics, img,
					int(p.x), int(p.y),
					int(p.w), int(p.h),
				))
			} else {
				gpMust(winapi.GdipDrawImageRectRectI(
					bmp.Graphics, img,
					int(p.x), int(p.y),
					int(p.w), int(p.h),
					int(p.img.region[0]), int(p.img.region[1]),
					int(p.img.region[2]), int(p.img.region[3]),
					winapi.GpUnitPixel,
					0, 0, nil,
				))
			}
		case CommandTypeArc:
		case CommandTypeArcFilled:
		default:
		}
	}
	if err := bmp.UseHDC(func(hdc syscall.Handle) error {
		return winapi.BitBlt(state.hdc, 0, 0, bmp.Width, bmp.Height, hdc, 0, 0, winapi.SRCCOPY)
	}); err != nil {
		panic(err)
	}
	NkClear(state.ctx)
}

type platformState struct {
	win                    *w32.Window
	wheelY                 float32
	mouseL, mouseR, mouseM bool
	hdc                    syscall.Handle
	pen                    winapi.GpPen
	brush                  winapi.GpBrush
	hrgn                   syscall.Handle

	bmp *internal.Bitmap
	ctx *Context
}

var state = &platformState{}
