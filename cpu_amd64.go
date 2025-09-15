//go:build amd64 && !purego

package uint256

import (
	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2
)
