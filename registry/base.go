package registry

import (
	"unsafe"

	"github.com/ginger-core/errors"
	"github.com/spf13/viper"
)

type base struct {
	viper *viper.Viper
}

func _new(format string) *base {
	c := &base{
		viper: viper.NewWithOptions(viper.KeyDelimiter(".")),
	}
	c.viper.SetConfigType(format)
	return c
}

func (c *base) WithDelimiter(delimiter string) Registry {
	keyDelim := (*string)(unsafe.Pointer(c.viper))
	*keyDelim = delimiter
	return c
}

func (c *base) ValueOf(key string) Registry {
	v := c.viper.Sub(key)
	if v == nil {
		return nil
	}
	return &base{
		viper: v,
	}
}

func (c *base) Unmarshal(ref interface{}, opts ...viper.DecoderConfigOption) errors.Error {
	if err := c.viper.Unmarshal(ref, opts...); err != nil {
		return errors.New(err)
	}
	return nil
}
