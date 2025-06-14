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

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/ajanata/pictureframe"
	"github.com/ajanata/pictureframe/internal/req"
)

type Walltaker struct {
	c Config

	tick     *time.Ticker
	img      widget.Image
	newImage image.Image
	imgLock  sync.Mutex
}

var _ pictureframe.Module = (*Walltaker)(nil)

var black = color.NRGBA{A: 0xFF}

func New(c Config) *Walltaker {
	if !c.Enabled {
		return nil
	}

	return &Walltaker{c: c}
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

func (w *Walltaker) Render(gtx layout.Context) error {
	if w == nil {
		return nil
	}

	w.imgLock.Lock()
	if w.newImage != nil {
		imgOp := paint.NewImageOp(w.newImage)
		imgOp.Filter = paint.FilterNearest
		w.img.Src = imgOp
		w.newImage = nil
	}
	w.imgLock.Unlock()

	layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
			paint.Fill(gtx.Ops, black)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}, func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return w.img.Layout(gtx)
			})
		})

	return nil
}

func (w *Walltaker) update() {
	img, err := w.getNewImage()
	if err != nil {
		log.Println(err)
		return
	}

	w.imgLock.Lock()
	w.newImage = img
	w.imgLock.Unlock()
}

func (w *Walltaker) getNewImage() (image.Image, error) {
	ctx := context.Background()

	link := map[string]interface{}{}
	err := req.GetJson(ctx, fmt.Sprintf("https://walltaker.joi.how/api/links/%d.json", w.c.LinkID), &link)
	if err != nil {
		return nil, err
	}

	reader, err := req.GetRaw(ctx, link["post_url"].(string))
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(reader)
	_ = reader.Close()
	if err != nil {
		return nil, err
	}

	return img, nil
}
