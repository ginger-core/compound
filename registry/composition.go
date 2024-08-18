package registry

import (
	"context"
	"os"

	"github.com/ginger-core/errors"
)

func New(ctx context.Context, args ...any) (Registry, errors.Error) {
	format := os.Getenv("CONFIG_FORMAT")
	if format == "" {
		format = "yaml"
	}
	configTypeStr := os.Getenv("CONFIG_TYPE")
	var configType Type
	switch configTypeStr {
	case "FILE", "":
		configType = TypeFile
	case "GIT", "REMOTE":
		configType = TypeRemote
	default:
		panic("invalid config type")
	}
	//
	switch configType {
	case TypeFile:
		c, err := newFile(ctx, format, args...)
		if err != nil {
			return nil, err.
				WithTrace("newFile")
		}
		return c, nil
	case TypeRemote:
		c, err := newGit(ctx, format)
		if err != nil {
			return nil, err.
				WithTrace("newGit")
		}
		return c, nil
	}
	return nil, errors.Internal().
		WithContext(ctx).
		WithId("unknown.type").
		WithMessage("unknown type")
}
