package config

import (
	"time"

	"github.com/ajanata/pictureframe/modules/walltaker"
)

type Config struct {
	Fullscreen bool
	FadeDelay  time.Duration
	FadeSpeed  uint

	Walltaker walltaker.Config
}
