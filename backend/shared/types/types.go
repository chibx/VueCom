package types

import (
	"go.uber.org/zap/zapcore"
)

type zapPrefixCore struct {
	zapcore.Core
	prefix string
}

func NewZapPrefix(core zapcore.Core, prefix string) *zapPrefixCore {
	return &zapPrefixCore{
		Core:   core,
		prefix: prefix,
	}
}

func (c *zapPrefixCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	ent.Message = c.prefix + ent.Message // ← this is the magic
	return c.Core.Write(ent, fields)
}

func (c *zapPrefixCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		// Register *this* wrapper as the core that will receive Write
		return ce.AddCore(ent, c)
	}
	return ce
}

type Pagination struct {
	Page     int
	PageSize int
}
