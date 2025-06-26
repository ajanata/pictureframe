package clock

import (
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"

	"github.com/ajanata/pictureframe"
)

var white = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

type Clock struct {
	c Config

	theme *material.Theme
}

var _ pictureframe.Module = (*Clock)(nil)

func New(c Config) *Clock {
	if !c.Enabled {
		return nil
	}

	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	return &Clock{
		c:     c,
		theme: theme,
	}
}

func (*Clock) Init() error { return nil }

func (c *Clock) Render(gtx layout.Context, _ byte, _ *app.Window) layout.Dimensions {
	if c == nil {
		return layout.Dimensions{}
	}

	return layout.SE.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// lbl := material.H5(c.theme, time.Now().Format(time.Kitchen))
		lbl := material.H5(c.theme, time.Now().Format(time.StampMilli))
		lbl.Color = white
		lbl.Alignment = text.End
		return lbl.Layout(gtx)
	})
}
