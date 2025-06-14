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

	l := link{}
	err := req.GetJson(ctx, fmt.Sprintf("https://walltaker.joi.how/api/links/%d.json", w.c.LinkID), &l)
	if err != nil {
		return nil, err
	}

	reader, err := req.GetRaw(ctx, l.PostUrl)
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
