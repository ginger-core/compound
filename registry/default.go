package registry

func (c *base) SetDefault(key string, value any) {
	c.viper.SetDefault(key, value)
}
