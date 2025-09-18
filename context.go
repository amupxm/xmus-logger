package xmuslogger

type Context struct {
	logger *Logger
}

func (c *Context) Str(key, val string) *Context {
	c.logger.context = appendString(c.logger.context, key, val)
	return c
}

func (c *Context) Int(key string, val int) *Context {
	c.logger.context = appendInt(c.logger.context, key, val)
	return c
}

func (c *Context) Bool(key string, val bool) *Context {
	c.logger.context = appendBool(c.logger.context, key, val)
	return c
}
func (c *Context) Logger() *Logger {
	return c.logger
}
