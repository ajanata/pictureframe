package control

import (
	"image/color"
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

var white = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

type Control struct {
	c Config

	theme *material.Theme

	chkShow   widget.Bool
	btnQuit   widget.Clickable
	btnFull   widget.Clickable
	btnMax    widget.Clickable
	btnWindow widget.Clickable
}

var _ pictureframe.Module = (*Control)(nil)

func New(c Config) *Control {
	if !c.Enabled {
		return nil
	}

	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return &Control{
		c:     c,
		theme: theme,
	}
}

func (*Control) Init() error { return nil }

func (c *Control) Render(gtx layout.Context, alpha byte, window *app.Window) layout.Dimensions {
	if c == nil {
		return layout.Dimensions{}
	}

	_ = c.chkShow.Update(gtx)

	if c.btnQuit.Clicked(gtx) {
		os.Exit(0)
	}
	if c.btnFull.Clicked(gtx) {
		window.Option(app.Fullscreen.Option())
	}
	if c.btnMax.Clicked(gtx) {
		window.Option(app.Maximized.Option())
	}
	if c.btnWindow.Clicked(gtx) {
		window.Option(app.Windowed.Option())
	}

	return layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return c.renderButtons(gtx, alpha)
	})
}

func (c *Control) renderButtons(gtx layout.Context, alpha byte) layout.Dimensions {
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
			Last:      4,
			Offset:    0,
			OffsetAbs: 0,
			Length:    25,
		},
	}

	ba := minByte(buttonAlpha, alpha)

	return grid.Layout(gtx, 5, 1, dim, func(gtx layout.Context, row, _ int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(3)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch row {
			case 0:
				a := ba
				if c.chkShow.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				chk := material.CheckBox(c.theme, &c.chkShow, "show")
				chk.Color = white
				chk.Color.A = a
				chk.IconColor.A = a
				chk.TextSize = unit.Sp(17)
				return chk.Layout(gtx)
			case 1:
				if !c.chkShow.Value {
					break
				}
				a := ba
				if c.btnQuit.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return c.button(&c.btnQuit, "quit", a).Layout(gtx)
			case 2:
				if !c.chkShow.Value {
					break
				}
				a := ba
				if c.btnFull.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return c.button(&c.btnFull, "full", a).Layout(gtx)
			case 3:
				if !c.chkShow.Value {
					break
				}
				a := ba
				if c.btnMax.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return c.button(&c.btnMax, "max", a).Layout(gtx)
			case 4:
				if !c.chkShow.Value {
					break
				}
				a := ba
				if c.btnWindow.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return c.button(&c.btnWindow, "window", a).Layout(gtx)
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

func (c *Control) button(btn *widget.Clickable, label string, alpha byte) material.ButtonStyle {
	b := material.Button(c.theme, btn, label)
	b.Background.A = alpha
	b.Color.A = alpha
	return b
}
