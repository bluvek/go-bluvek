package bvmiddleware

import (
	"strings"

	"github.com/bluvek/go-bluvek/core"
	"github.com/bluvek/go-bluvek/pkg/bvauth"
	"github.com/bluvek/go-bluvek/pkg/bverror"
	"github.com/gin-gonic/gin"
)

func Jwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			core.Fail(ctx, bverror.NeedLogin)
			ctx.Abort()
			return
		}

		claims, err := bvauth.ParseJwtToken(tokenString[7:])
		if err != nil {
			core.Fail(ctx, bverror.NeedLogin)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
