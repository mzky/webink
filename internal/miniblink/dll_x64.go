//go:build amd64

package miniblink

import (
	"embed"
)

const (
	ARCH    = "x64"
	VERSION = "4975"
)

//go:embed release/x64
var res embed.FS
