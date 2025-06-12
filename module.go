package pictureframe

import (
	"gioui.org/layout"
)

type Module interface {
	Init() error
	Render(layout.Context) error
}
