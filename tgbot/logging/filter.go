package logging

import "go.uber.org/zap/zapcore"

// from https://github.com/moul/zapfilter/

type FilterFunc = func(zapcore.Entry, []zapcore.Field) bool

func NewFilteringCore(next zapcore.Core, filter FilterFunc) zapcore.Core {
	if filter == nil {
		filter = alwaysFalseFilter
	}
	return &filteringCore{next, filter}
}

type filteringCore struct {
	next   zapcore.Core
	filter FilterFunc
}

func alwaysFalseFilter(_ zapcore.Entry, _ []zapcore.Field) bool {
	return false
}

func alwaysTrueFilter(_ zapcore.Entry, _ []zapcore.Field) bool {
	return true
}

func (c *filteringCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.filter(ent, nil) {
		ce = ce.AddCore(ent, c)
	}

	return ce
}

func (c *filteringCore) With(fields []zapcore.Field) zapcore.Core {
	return &filteringCore{c.next.With(fields), c.filter}
}

func (c *filteringCore) Enabled(level zapcore.Level) bool {
	return c.next.Enabled(level)
}

func (c *filteringCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if c.filter(ent, fields) {
		return c.next.Write(ent, fields)
	}
	return nil
}

func (c *filteringCore) Sync() error {
	return c.next.Sync()
}
