package tools

import "strings"

func ConvImagePath(path string) string {

	if !strings.HasPrefix(path, "http") {
		path = "https://extest123.sukeeper.com/obs/" + path
	}

	return path
}
