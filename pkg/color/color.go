package color

import "runtime"

var (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

func init() {
	if runtime.GOOS == "windows" {
		green = ""
		yellow = ""
	}
}

func Green(text string) string {
	return green + text + reset
}

func Yellow(text string) string {
	return yellow + text + reset
}
