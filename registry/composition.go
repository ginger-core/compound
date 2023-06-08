package registry

import (
	"context"

	"github.com/ginger-core/errors"
)

func New(ctx context.Context, t Type, format string, args ...interface{}) (Registry, errors.Error) {
	switch t {
	case TypeFile:
		c, err := newFile(ctx, format, args...)
		if err != nil {
			return nil, err.
				WithTrace("newFile")
		}
		return c, nil
	case TypeGitAPI:
		c, err := newGit(ctx, format, args...)
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
