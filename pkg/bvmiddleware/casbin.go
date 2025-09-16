package bvmiddleware

import (
	"strings"

	"github.com/bluvek/go-bluvek/console"
	"github.com/bluvek/go-bluvek/core"
	"github.com/bluvek/go-bluvek/pkg/bvauth"
	"github.com/bluvek/go-bluvek/pkg/bverror"
	"github.com/bluvek/go-bluvek/pkg/bvutils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func Casbin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(ctx.Request.Method) == "OPTIONS" {
			ctx.Next()
			return
		}

		roleId := bvauth.GetTokenValue[int64](ctx, "role_id")
		if roleId == 0 {
			console.Echo.Info("ℹ️ 提示: 无法使用 `Casbin` 权限校验, 请确保 `Token` 中包含了字段 `role_id`")
			ctx.Next()
			return
		}

		path := bvutils.ConvertToRestfulURL(strings.TrimPrefix(ctx.Request.URL.Path, viper.GetString("App.RouterPrefix")))
		success, _ := core.Casbin.Enforce(cast.ToString(roleId), path, ctx.Request.Method)
		if !success {
			core.Fail(ctx, bverror.NoAuth)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
