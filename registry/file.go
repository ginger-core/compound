package registry

import (
	"context"
	"io"
	"os"

	"github.com/ginger-core/errors"
)

type file struct {
	*base
}

func newFile(ctx context.Context, format string, args ...interface{}) (Registry, errors.Error) {
	if len(args) == 0 {
		return nil, errors.Internal().
			WithContext(ctx).
			WithTrace("noPath").
			WithMessage("file path not given")
	}

	var reader io.Reader
	var err error

	if path, ok := args[0].(string); ok {
		reader, err = os.Open(path)
		if err != nil {
			return nil, errors.Internal(err).
				WithContext(ctx).
				WithTrace("os.open").
				WithMessage("cannot read file")
		}
	} else if r, ok := args[0].(io.Reader); ok {
		reader = r
	} else {
		return nil, errors.Internal().
			WithContext(ctx).
			WithTrace("arg0.cast").
			WithMessage("invalid arg")
	}

	c := &file{
		base: _new(format),
	}
	if err = c.viper.ReadConfig(reader); err != nil {
		return nil, errors.Internal(err).
			WithContext(ctx).
			WithTrace("viper.ReadConfig").
			WithMessage("cannot read file")
	}
	return c, nil
}
