package statsig

import (
	"context"
	"fmt"

	statsigsdk "github.com/statsig-io/go-sdk"
)

type ctxStatsigUserMarker struct{}

type ctxStatsig struct {
	user statsigsdk.User
}

var (
	ctxStatsigUserKey = &ctxStatsigUserMarker{}
)

// cmnt me
func Extract(ctx context.Context) (statsigsdk.User, error) {
	r, ok := ctx.Value(ctxStatsigUserKey).(*ctxStatsig)
	if !ok || r == nil {
		return statsigsdk.User{}, fmt.Errorf("unable to get the statsig user")
	}

	return r.user, nil
}

// cmnt me
func ToContext(ctx context.Context, u statsigsdk.User) context.Context {
	r := &ctxStatsig{
		user: u,
	}
	return context.WithValue(ctx, ctxStatsigUserKey, r)
}
