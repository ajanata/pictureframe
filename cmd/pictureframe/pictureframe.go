package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/BurntSushi/toml"

	"github.com/ajanata/pictureframe"
	"github.com/ajanata/pictureframe/internal/config"
	"github.com/ajanata/pictureframe/modules/scratch"
	"github.com/ajanata/pictureframe/modules/walltaker"
)

var c config.Config
var modules []pictureframe.Module

var black = color.NRGBA{A: 0xFF}

func main() {
	_, err := toml.DecodeFile("pictureframe.toml", &c)
	if err != nil {
		log.Fatal(err)
	}

	modules = append(modules, walltaker.New(c.Walltaker))
	modules = append(modules, scratch.New(c.Scratch))

	for _, m := range modules {
		if err := m.Init(); err != nil {
			log.Fatal(err)
		}
	}

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

func run(window *app.Window) error {
	var ops op.Ops
	var alpha byte = 0xFF

	lastMove := time.Now()
	var lastPos f32.Point
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			var tag bool
			event.Op(&ops, tag)
			ev, ok := gtx.Event(pointer.Filter{
				Target: tag,
				Kinds:  pointer.Move,
			})
			if ok {
				ev := ev.(pointer.Event)
				// we get spurious events even if the position hasn't changed, so explicitly check it
				if ev.Position != lastPos {
					lastMove = time.Now()
					lastPos = ev.Position
					alpha = 0xFF
				}
			}

			if time.Since(lastMove) > c.FadeDelay && alpha > 0 {
				for i := c.FadeSpeed; i > 0 && alpha > 0; i-- {
					alpha--
				}
			}

			layout.Background{}.Layout(gtx,
				renderBackground(),
				renderModule(0, alpha, window),
			)

			gtx.Execute(op.InvalidateCmd{At: time.Now().Add(time.Second / 60)})
			e.Frame(gtx.Ops)
		}
	}
}

func renderBackground() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, black)
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
}

func renderModule(i int, alpha byte, window *app.Window) layout.Widget {
	if i+1 < len(modules) {
		// there's another module to render
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return modules[i].Render(gtx, alpha, window)
				},
				renderModule(i+1, alpha, window),
			)
		}
	} else {
		// this is the last module
		return func(gtx layout.Context) layout.Dimensions {
			return modules[i].Render(gtx, alpha, window)
		}
	}
}
