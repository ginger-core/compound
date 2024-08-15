package registry

import (
	"context"
	"os"
	"strings"

	"github.com/ginger-core/errors"
)

var readerMap = map[string]func(base *base) reader{
	"gitlab": newGitlabReader,
	"github": newGithubReader,
}

type git struct {
	*base
	reader
}

type reader interface {
	read() errors.Error
}

func newGit(_ context.Context, format string) (Registry, errors.Error) {
	c := &git{
		base: _new(format),
	}
	url := os.Getenv("CONFIG_URL")
	if url == "" {
		return nil, errors.Validation().
			WithTrace("InvalidConfigUrl").
			WithMessage("Config url environment is empty." +
				" Set environment vartiale of key `CONFIG_URL` to your config url path")
	}
	remoteType := os.Getenv("CONFIG_REMOTE_TYPE")
	switch strings.ToUpper(remoteType) {
	case "GITHUB":
		c.reader = newGithubReader(c.base)
	case "GITLAB":
		c.reader = newGitlabReader(c.base)
	default:
		switch {
		case strings.Contains(url, "github"):
			c.reader = newGithubReader(c.base)
		case strings.Contains(url, "gitlab"):
			c.reader = newGitlabReader(c.base)
		default:
			return nil, errors.Validation().
				WithTrace("RemoteTypeNotFound").
				WithMessage("config remote type not found")
		}
	}
	if err := c.read(); err != nil {
		return nil, err
	}
	return c, nil
}
