package walltaker

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"log"
	"net/http"
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/ajanata/pictureframe"
)

type Walltaker struct {
	c Config

	tick     *time.Ticker
	img      widget.Image
	newImage image.Image
	imgLock  sync.Mutex
}

var _ pictureframe.Module = (*Walltaker)(nil)

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

	w.img.Layout(gtx)

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
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://walltaker.joi.how/api/links/%d.json", w.c.LinkID),
		nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	decode := json.NewDecoder(resp.Body)
	link := map[string]interface{}{}
	err = decode.Decode(&link)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	req, err = http.NewRequest(http.MethodGet, link["post_url"].(string), nil)
	if err != nil {
		return nil, err
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	_ = resp.Body.Close()

	return img, nil
}
