package utils

import (
	"runtime"
)

func IsWindowsMachine() bool {
	return runtime.GOOS == "windows"
}
