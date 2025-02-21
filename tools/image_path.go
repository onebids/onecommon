package tools

import "strings"

func ConvImagePath(path string, baseUrl string) string {

	if !strings.HasPrefix(path, "http") {
		path = baseUrl + path
	}

	return path
}
