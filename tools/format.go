package tools

import (
	"strconv"
	"strings"
)

// BoolFormat 根据条件返回两个值中的一个，类似三目运算符
//
// 如果 expr 为 true，则返回 a，否则返回 b
//
// 参数:
//   - expr: 条件表达式
//   - a: 条件为真时返回的值
//   - b: 条件为假时返回的值
//
// 返回:
//   - 根据条件选择的值
func BoolFormat(expr bool, a, b interface{}) interface{} {
	if expr {
		return a
	}
	return b
}

// BoolFormatInt32 根据条件返回两个int32值中的一个
//
// 如果 expr 为 true，则返回 a，否则返回 b
//
// 参数:
//   - expr: 条件表达式
//   - a: 条件为真时返回的int32值
//   - b: 条件为假时返回的int32值
//
// 返回:
//   - 根据条件选择的int32值
func BoolFormatInt32(expr bool, a, b int32) int32 {
	if expr {
		return a
	}
	return b
}

// BoolFormatStr 根据条件返回两个字符串中的一个
//
// 如果 expr 为 true，则返回 a，否则返回 b
//
// 参数:
//   - expr: 条件表达式
//   - a: 条件为真时返回的字符串
//   - b: 条件为假时返回的字符串
//
// 返回:
//   - 根据条件选择的字符串
func BoolFormatStr(expr bool, a, b string) string {
	if expr {
		return a
	}
	return b
}

// Int64SliceToString 将int64切片转换为逗号分隔的字符串
//
// 参数:
//   - ids: 需要转换的int64切片
//
// 返回:
//   - 逗号分隔的字符串
func Int64SliceToString(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}

	strSlice := make([]string, len(ids))
	for i, id := range ids {
		strSlice[i] = strconv.FormatInt(id, 10)
	}
	return strings.Join(strSlice, ",")
}
