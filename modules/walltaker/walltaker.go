package walltaker

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/outlay"
	"github.com/nfnt/resize"

	"github.com/ajanata/pictureframe"
	"github.com/ajanata/pictureframe/internal/req"
)

const buttonAlpha = 0x40

type Walltaker struct {
	c Config

	tick        *time.Ticker
	img         widget.Image
	newImage    image.Image
	currentLink link
	imgLock     sync.Mutex

	theme     *material.Theme
	btnLoveIt widget.Clickable
	btnHateIt widget.Clickable
	btnCame   widget.Clickable
}

type link struct {
	Id               int         `json:"id"`
	Expires          interface{} `json:"expires"`
	Terms            string      `json:"terms"`
	Blacklist        string      `json:"blacklist"`
	PostUrl          string      `json:"post_url"`
	PostThumbnailUrl string      `json:"post_thumbnail_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	ResponseType     string      `json:"response_type"`
	ResponseText     string      `json:"response_text"`
	Username         string      `json:"username"`
	SetBy            string      `json:"set_by"`
	Online           bool        `json:"online"`
	PostDescription  string      `json:"post_description"`
	Url              string      `json:"url"`
}

type response struct {
	ApiKey string `json:"api_key"`
	Type   string `json:"type"`
	Text   string `json:"text"`
}

type reaction string

const loveIt reaction = "horny"
const hateIt reaction = "disgust"
const came reaction = "came"

var _ pictureframe.Module = (*Walltaker)(nil)

var white = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

func New(c Config) *Walltaker {
	if !c.Enabled {
		return nil
	}

	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	return &Walltaker{
		c:     c,
		theme: theme,
	}
}

func (w *Walltaker) Init() error {
	if w == nil {
		return nil
	}

	if w.c.CheckInterval < time.Second {
		return errors.New("interval too small")
	}

	w.tick = time.NewTicker(w.c.CheckInterval)
	go func() {
		for range w.tick.C {
			w.update()
		}
	}()

	w.img.Fit = widget.Contain
	w.update()

	return nil
}

func (w *Walltaker) Render(gtx layout.Context, alpha byte, _ *app.Window) layout.Dimensions {
	if w == nil {
		return layout.Dimensions{}
	}

	if w.btnLoveIt.Clicked(gtx) {
		w.buttonPressed(loveIt)
	}
	if w.btnHateIt.Clicked(gtx) {
		w.buttonPressed(hateIt)
	}
	if w.btnCame.Clicked(gtx) {
		w.buttonPressed(came)
	}

	w.imgLock.Lock()
	if w.newImage != nil {
		imgOp := paint.NewImageOp(w.newImage)
		imgOp.Filter = paint.FilterNearest
		w.img.Src = imgOp
		w.newImage = nil
	}
	w.imgLock.Unlock()

	return layout.Background{}.Layout(gtx,
		w.renderImage(),
		func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return maybe(alpha, w.renderCaption(gtx, alpha))
				},
				func(gtx layout.Context) layout.Dimensions {
					return layout.SW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return maybe(alpha, w.renderButtons(gtx, alpha))
					})
				},
			)
		},
	)
}

func (w *Walltaker) renderCaption(gtx layout.Context, alpha byte) layout.Dimensions {
	// updated at is definitely not working right
	// ago := time.Now().Sub(w.currentLink.UpdatedAt).Round(time.Minute)
	// lbl := material.H5(w.theme, fmt.Sprintf("set by %s %s ago", w.currentLink.SetBy, ago))
	lbl := material.H5(w.theme, fmt.Sprintf("set by %s", w.currentLink.SetBy))
	lbl.Color = white
	lbl.Color.A = alpha
	lbl.Alignment = text.Middle
	return lbl.Layout(gtx)
}

func (w *Walltaker) renderButtons(gtx layout.Context, alpha byte) layout.Dimensions {
	dim := func(axis layout.Axis, index, constraint int) int {
		switch axis {
		case layout.Vertical:
			return 50
		case layout.Horizontal:
			return 75
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
			Last:      2,
			Offset:    0,
			OffsetAbs: 0,
			Length:    25,
		},
	}

	ba := minByte(buttonAlpha, alpha)

	return grid.Layout(gtx, 3, 1, dim, func(gtx layout.Context, row, _ int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(3)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch row {
			case 0:
				a := ba
				if w.btnLoveIt.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return w.button(&w.btnLoveIt, "love", a).Layout(gtx)
			case 1:
				a := ba
				if w.btnHateIt.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return w.button(&w.btnHateIt, "hate", a).Layout(gtx)
			case 2:
				a := ba
				if w.btnCame.Hovered() {
					a = minByte(alpha, 0xFF)
				}
				return w.button(&w.btnCame, "came", a).Layout(gtx)
			}
			return layout.Dimensions{}
		})
	})
}

func (w *Walltaker) renderImage() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return w.img.Layout(gtx)
		})
	}
}

func minByte(a, b byte) byte {
	if a < b {
		return a
	}
	return b
}

func (w *Walltaker) button(btn *widget.Clickable, label string, alpha byte) material.ButtonStyle {
	b := material.Button(w.theme, btn, label)
	b.Background.A = alpha
	b.Color.A = alpha
	return b
}

func (w *Walltaker) update() {
	l, img, err := w.getNewImage()
	if err != nil {
		log.Println(err)
		return
	}

	w.imgLock.Lock()
	w.newImage = img
	w.currentLink = l
	w.imgLock.Unlock()
}

func (w *Walltaker) getNewImage() (link, image.Image, error) {
	ctx := context.Background()

	l := link{}
	err := req.GetJson(ctx, fmt.Sprintf("https://walltaker.joi.how/api/links/%d.json", w.c.LinkID), &l)
	if err != nil {
		return l, nil, err
	}

	reader, err := req.GetRaw(ctx, l.PostUrl)
	if err != nil {
		return l, nil, err
	}

	img, _, err := image.Decode(reader)
	_ = reader.Close()
	if err != nil {
		return l, nil, err
	}

	img = resize.Thumbnail(w.c.MaxSize.Width, w.c.MaxSize.Height, img, resize.Lanczos3)

	return l, img, nil
}

func (w *Walltaker) buttonPressed(reaction reaction) {
	r := response{
		ApiKey: w.c.APIKey,
		Type:   string(reaction),
	}

	err := req.Post(context.Background(), fmt.Sprintf("https://walltaker.joi.how/api/links/%d/response.json", w.c.LinkID), r)
	if err != nil {
		fmt.Println(err)
		return
	}

	if reaction == hateIt {
		w.update()
	}
}

func maybe(alpha byte, dim layout.Dimensions) layout.Dimensions {
	if alpha > 0 {
		return dim
	}
	return layout.Dimensions{}
}
