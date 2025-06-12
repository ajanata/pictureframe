package config

import (
	"github.com/ajanata/pictureframe/modules/walltaker"
)

type Config struct {
	Fullscreen bool

	Walltaker walltaker.Config
}
