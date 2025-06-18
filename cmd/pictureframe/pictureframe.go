package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"github.com/BurntSushi/toml"

	"github.com/ajanata/pictureframe"
	"github.com/ajanata/pictureframe/internal/config"
	"github.com/ajanata/pictureframe/modules/walltaker"
)

var c config.Config
var modules []pictureframe.Module

func main() {
	_, err := toml.DecodeFile("pictureframe.toml", &c)
	if err != nil {
		log.Fatal(err)
	}

	modules = append(modules, walltaker.New(c.Walltaker))

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

			for _, m := range modules {
				err := m.Render(gtx, alpha)
				if err != nil {
					return err
				}
			}

			gtx.Execute(op.InvalidateCmd{At: time.Now().Add(time.Second / 60)})
			e.Frame(gtx.Ops)
		}
	}
}
