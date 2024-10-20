package tools

// BoolFormat 三目运算
func BoolFormat(expr bool, a, b interface{}) interface{} {
	if expr {
		return a
	}
	return b
}

func BoolFormatInt32(expr bool, a, b int32) int32 {
	if expr {
		return a
	}
	return b
}

func BoolFormatStr(expr bool, a, b string) string {
	if expr {
		return a
	}
	return b
}
