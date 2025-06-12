package walltaker

import (
	"time"
)

type Config struct {
	Enabled bool

	LinkID        int
	APIKey        string
	CheckInterval time.Duration
}
