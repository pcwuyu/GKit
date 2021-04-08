package bbr

import (
	"github.com/Songzhibin/GKit/middleware"
	"github.com/Songzhibin/GKit/options"
	"github.com/Songzhibin/GKit/overload"
	"context"
)

const (
	LimitKey = "LimitKey"
	LimitOp  = "LimitLoad"
)

func NewLimiter(options ...options.Option) middleware.MiddleWare {
	g := NewGroup(options...)
	return func(next middleware.Endpoint) middleware.Endpoint {
		return func(ctx context.Context, i interface{}) (interface{}, error) {
			// 通过ctx 获取 g中的限制器
			defaultKey := "default"
			defaultOp := overload.Success
			if v := ctx.Value(LimitKey); v != nil {
				defaultKey = v.(string)
			}
			if v := ctx.Value(LimitOp); v != nil {
				defaultOp = v.(overload.Op)
			}
			limiter := g.Get(defaultKey)
			if f, err := limiter.Allow(ctx); err != nil {
				return nil, err
			} else {
				f(overload.DoneInfo{Op: defaultOp})
				return next(ctx, i)
			}
		}
	}
}
