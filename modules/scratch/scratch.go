package scratch

import (
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/outlay"

	"github.com/ajanata/pictureframe"
)

const buttonAlpha = 0x40

type Scratch struct {
	c Config

	theme     *material.Theme
	btnQuit   widget.Clickable
	btnFull   widget.Clickable
	btnMax    widget.Clickable
	btnWindow widget.Clickable
}

var _ pictureframe.Module = (*Scratch)(nil)

func New(c Config) *Scratch {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return &Scratch{
		c:     c,
		theme: theme,
	}
}

func (*Scratch) Init() error { return nil }

func (s *Scratch) Render(gtx layout.Context, alpha byte, window *app.Window) layout.Dimensions {
	if s == nil {
		return layout.Dimensions{}
	}

	if s.btnQuit.Clicked(gtx) {
		os.Exit(0)
	}
	if s.btnFull.Clicked(gtx) {
		window.Option(app.Fullscreen.Option())
	}
	if s.btnMax.Clicked(gtx) {
		window.Option(app.Maximized.Option())
	}
	if s.btnWindow.Clicked(gtx) {
		window.Option(app.Windowed.Option())
	}

	return layout.NE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return s.renderButtons(gtx, alpha)
	})
}

func (s *Scratch) renderButtons(gtx layout.Context, alpha byte) layout.Dimensions {
	dim := func(axis layout.Axis, index, constraint int) int {
		switch axis {
		case layout.Vertical:
			return 50
		case layout.Horizontal:
			return 85
		default:
			return 0
		}
	}

	grid := &outlay.Grid{
		Horizontal: outlay.AxisPosition{
			First:     0,
			Last:      0,
			Offset:    0,
			OffsetAbs: 0,
			Length:    50,
		},
		Vertical: outlay.AxisPosition{
			First:     0,
			Last:      3,
			Offset:    0,
			OffsetAbs: 0,
			Length:    25,
		},
	}

	ba := minByte(buttonAlpha, alpha)

	return grid.Layout(gtx, 4, 1, dim, func(gtx layout.Context, row, _ int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(3)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch row {
			case 0:
				a := ba
				if s.btnQuit.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return s.button(&s.btnQuit, "quit", a).Layout(gtx)
			case 1:
				a := ba
				if s.btnFull.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return s.button(&s.btnFull, "full", a).Layout(gtx)
			case 2:
				a := ba
				if s.btnMax.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return s.button(&s.btnMax, "max", a).Layout(gtx)
			case 3:
				a := ba
				if s.btnWindow.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return s.button(&s.btnWindow, "window", a).Layout(gtx)
			}
			return layout.Dimensions{}
		})
	})
}

func minByte(a, b byte) byte {
	if a < b {
		return a
	}
	return b
}

func (s *Scratch) button(btn *widget.Clickable, label string, alpha byte) material.ButtonStyle {
	b := material.Button(s.theme, btn, label)
	b.Background.A = alpha
	b.Color.A = alpha
	return b
}
