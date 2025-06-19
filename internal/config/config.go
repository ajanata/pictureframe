package config

import (
	"time"

	"github.com/ajanata/pictureframe/modules/scratch"
	"github.com/ajanata/pictureframe/modules/walltaker"
)

type Config struct {
	Fullscreen      bool
	FullscreenDelay time.Duration
	FadeDelay       time.Duration
	FadeSpeed       uint

	Walltaker walltaker.Config
	Scratch   scratch.Config
}
