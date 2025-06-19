package pictureframe

import (
	"gioui.org/app"
	"gioui.org/layout"
)

type Module interface {
	Init() error
	Render(gtx layout.Context, alpha byte, window *app.Window) layout.Dimensions
}
