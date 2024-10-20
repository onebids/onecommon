package tools

// BoolFormat 三目运算
func BoolFormat(expr bool, a, b interface{}) interface{} {
	if expr {
		return a
	}
	return b
}
