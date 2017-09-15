package internal

import (
	"errors"
	"syscall"

	"github.com/golang-ui/nuklear/nk/internal/winapi"
)

type PrivateFont struct {
	name string
	h    syscall.Handle
}

func (pf *PrivateFont) Name() string {
	return pf.name
}

func (pf *PrivateFont) Close() error {
	winapi.RemoveFontMemResourceEx(pf.h)
	return nil
}

func (pf *PrivateFont) NewFont(size int) (*Font, error) {
	return NewFont(pf.name, size)
}

func NewPrivateFont(name string, data []byte) (*PrivateFont, error) {
	var fonts uint32
	h := winapi.AddFontMemResourceEx(&data[0], uint32(len(data)), nil, &fonts)
	if h == 0 {
		return nil, errors.New("w32gdi: could not add private font")
	}
	return &PrivateFont{name, h}, nil
}

type Font struct {
	Handle  syscall.Handle
	DC      syscall.Handle
	Height  int
	oldFont syscall.Handle
}

func (f *Font) Width(str string) (int, error) {
	s, err := syscall.UTF16FromString(str)
	if err != nil {
		return 0, err
	}
	var sz winapi.SIZE
	if err = winapi.GetTextExtentPoint32(f.DC, &s[0], len(s)-1, &sz); err != nil {
		return 0, err
	}
	return int(sz.CX), nil
}

func (f *Font) Close() error {
	h := winapi.SelectObject(f.DC, f.oldFont)
	winapi.DeleteDC(f.DC)
	winapi.DeleteObject(h)
	f.oldFont = 0
	f.DC = 0
	f.Height = 0
	return nil
}

func NewFont(name string, size int) (*Font, error) {
	s, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	dc, err := winapi.CreateCompatibleDC(0)
	if err != nil {
		return nil, err
	}
	h, err := winapi.CreateFont(
		int32(size),
		0,
		0,
		0,
		winapi.FW_NORMAL,
		0,
		0,
		0,
		winapi.DEFAULT_CHARSET,
		winapi.OUT_DEFAULT_PRECIS,
		winapi.CLIP_DEFAULT_PRECIS,
		winapi.DEFAULT_QUALITY,
		winapi.DEFAULT_PITCH|winapi.FF_DONTCARE,
		s,
	)
	if err != nil {
		winapi.DeleteDC(dc)
		return nil, err
	}
	oldFont := winapi.SelectObject(dc, h)
	var tm winapi.TEXTMETRIC
	if err = winapi.GetTextMetrics(dc, &tm); err != nil {
		winapi.SelectObject(dc, oldFont)
		winapi.DeleteDC(dc)
		winapi.DeleteObject(h)
		return nil, err
	}
	return &Font{
		Handle:  h,
		DC:      dc,
		Height:  int(tm.Height),
		oldFont: oldFont,
	}, nil
}
