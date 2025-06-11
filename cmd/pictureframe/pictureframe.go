package main

import (
	"encoding/json"
	"image"
	"log"
	"net/http"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/BurntSushi/toml"

	"github.com/ajanata/pictureframe/internal/config"
)

var c config.Config
var img image.Image

func main() {
	_, err := toml.DecodeFile("pictureframe.toml", &c)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "https://walltaker.joi.how/api/links/42586.json", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	decode := json.NewDecoder(resp.Body)
	link := map[string]interface{}{}
	err = decode.Decode(&link)
	if err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()

	req, err = http.NewRequest(http.MethodGet, link["post_url"].(string), nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err = image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	_ = resp.Body.Close()

	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imgWidget.Src = imageOp
	imgWidget.Fit = widget.Contain

	go func() {
		window := new(app.Window)
		window.Option(app.Title("PictureFrame"))
		window.Option(app.MinSize(300, 200))
		if c.Fullscreen {
			window.Option(app.Fullscreen.Option())
		}
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

var imgWidget widget.Image

func run(window *app.Window) error {
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// widget.Image{}

			imgWidget.Layout(gtx)

			// drawImage(gtx.Ops, img)
			gtx.Execute(op.InvalidateCmd{At: time.Now().Add(time.Second / 60)})
			e.Frame(gtx.Ops)
		}
	}
}

func drawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Filter = paint.FilterNearest
	imageOp.Add(ops)
	// op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4))).Add(ops)
	paint.PaintOp{}.Add(ops)
}
