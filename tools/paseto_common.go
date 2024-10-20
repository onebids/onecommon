package tools

import (
	pt "aidanwoods.dev/go-paseto"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/paseto"
	"github.com/onebids/onecommon/consts"
	"github.com/onebids/onecommon/errno"
	"github.com/onebids/onecommon/model"
	"net/http"
)

func CommonMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		// use cors mw
		//middleware.Cors(),
		// use recovery mw
		//middleware.Recovery(),
		// use gzip mw
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
	}
}

func PasetoAuth(audience string, pi model.PasetoConfig) app.HandlerFunc {
	pf, err := paseto.NewV4PublicParseFunc(pi.PubKey, []byte(pi.Implicit), paseto.WithAudience(audience), paseto.WithNotBefore())
	if err != nil {
		hlog.Fatal(err)
	}
	sh := func(ctx context.Context, c *app.RequestContext, token *pt.Token) {
		aid, err := token.GetString("id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, BuildBaseResp(errno.BadRequest.WithMessage("missing accountID in token")))
			c.Abort()
			return
		}
		c.Set(consts.AccountID, aid)
	}

	eh := func(ctx context.Context, c *app.RequestContext) {
		//c.JSON(http.StatusUnauthorized, tools.BuildBaseResp(errno.thrift.BadRequest.WithMessage("invalid token")))
		c.Abort()
	}
	return paseto.New(paseto.WithTokenPrefix("Bearer "), paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(eh))
}
