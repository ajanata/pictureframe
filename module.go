package pictureframe

import (
	"gioui.org/layout"
)

type Module interface {
	Init() error
	Render(gtx layout.Context, alpha byte) error
}
