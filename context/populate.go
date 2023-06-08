package composition

import "context"

func Populate(source, dest context.Context) context.Context {
	for _, k := range ctxDefaultKeys {
		if v := source.Value(k); v != nil {
			dest = context.WithValue(dest, k, v)
		}
	}
	return dest
}
