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
	var filePath string
	if len(args) == 0 {
		//
		filePath = os.Getenv("CONFIG_PATH")
		if filePath == "" {
			filePath = "./config.yaml"
		}
	} else if path, ok := args[0].(string); ok {
		filePath = path
	}
	var reader io.Reader
	var err error
	if filePath != "" {
		reader, err = os.Open(filePath)
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
