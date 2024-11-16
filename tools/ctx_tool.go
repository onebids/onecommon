package tools

import (
	"context"
	"github.com/bytedance/gopkg/cloud/metainfo"
)

func GetCtxValue(ctx context.Context, key string, defaultValue string) string {
	lang, ok1 := metainfo.GetValue(ctx, key)
	if ok1 {
		return lang
	} else {
		return defaultValue
	}
}

func SetCtxValue(ctx context.Context, key string, value string) context.Context {
	ctx = metainfo.WithValue(ctx, key, value)
	return ctx
}
