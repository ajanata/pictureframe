package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
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
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			for _, m := range modules {
				err := m.Render(gtx)
				if err != nil {
					return err
				}
			}

			gtx.Execute(op.InvalidateCmd{At: time.Now().Add(time.Second / 60)})
			e.Frame(gtx.Ops)
		}
	}
}
