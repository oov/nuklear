package internal

import (
	"syscall"
	"unsafe"

	"github.com/golang-ui/nuklear/nk/internal/winapi"
)

type Bitmap struct {
	hdc      syscall.Handle
	oldBMP   syscall.Handle
	bmp      syscall.Handle
	Graphics winapi.GpGraphics
	Width    int
	Height   int
	P        *byte
}

func (b *Bitmap) Close() error {
	b.Graphics.Close()
	winapi.SelectObject(b.hdc, b.oldBMP)
	winapi.DeleteObject(b.bmp)
	winapi.DeleteDC(b.hdc)
	return nil
}

func (b *Bitmap) UseHDC(f func(hdc syscall.Handle) error) error {
	var hdc syscall.Handle
	if err := winapi.GdipGetDC(b.Graphics, &hdc).Err(); err != nil {
		return err
	}
	defer winapi.GdipReleaseDC(b.Graphics, hdc)
	return f(hdc)
}

func NewBitmap(width, height int) (*Bitmap, error) {
	desktop := winapi.GetDesktopWindow()
	deskdc, err := winapi.GetDC(desktop)
	if err != nil {
		return nil, err
	}
	defer winapi.ReleaseDC(desktop, deskdc)

	hdc, err := winapi.CreateCompatibleDC(deskdc)
	if err != nil {
		return nil, err
	}
	var p *byte
	bmp, err := winapi.CreateDIBSection(deskdc, &winapi.BITMAPINFO{
		Header: winapi.BITMAPINFOHEADER{
			Size:        uint32(unsafe.Sizeof(winapi.BITMAPINFOHEADER{})),
			Width:       int32(width),
			Height:      int32(height),
			Planes:      1,
			BitCount:    32,
			Compression: winapi.BI_RGB,
		},
	}, winapi.DIB_RGB_COLORS, &p, 0, 0)
	if err != nil {
		winapi.DeleteDC(hdc)
		return nil, err
	}
	oldBMP := winapi.SelectObject(hdc, bmp)

	var g winapi.GpGraphics
	if err = winapi.GdipCreateFromHDC(hdc, &g).Err(); err != nil {
		winapi.SelectObject(hdc, oldBMP)
		winapi.DeleteObject(bmp)
		winapi.DeleteDC(hdc)
		return nil, err
	}
	return &Bitmap{
		hdc:      hdc,
		oldBMP:   oldBMP,
		bmp:      bmp,
		Graphics: g,
		Width:    width,
		Height:   height,
		P:        p,
	}, nil
}
