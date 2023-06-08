package composition

import (
	"context"
)

func GetContextData(ctx context.Context) map[Key]interface{} {
	r := make(map[Key]interface{})
	for _, k := range ctxDefaultKeys {
		if v := ctx.Value(k); v != nil {
			r[k] = v
		}
	}
	return r
}

func NewBackgroundContext(ctx context.Context) context.Context {
	r := context.Background()
	return Populate(ctx, r)
}
