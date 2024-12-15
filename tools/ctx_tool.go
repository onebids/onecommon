package tools

import (
	"context"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/onebids/onecommon/consts"
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

func GetAccountId(ctx context.Context) string {
	return GetCtxValue(ctx, consts.AccountID, "")
}
func GetLanguage(ctx context.Context) string {
	return GetCtxValue(ctx, consts.Language, consts.DefaultLanguage)
}
