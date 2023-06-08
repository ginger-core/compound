package registry

import (
	"github.com/ginger-core/errors"
	"github.com/spf13/viper"
)

type Registry interface {
	SetDefault(key string, value any)
	WithDelimiter(delimiter string) Registry
	ValueOf(key string) Registry
	Unmarshal(ref interface{}, opts ...viper.DecoderConfigOption) errors.Error
}
