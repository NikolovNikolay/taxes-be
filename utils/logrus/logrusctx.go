package logrus

import "github.com/sirupsen/logrus"

// Package logrusctx provides a logrus hook which extracts the logging data
// from the context which has been accumulated via logctx.

// LogCtxHook reads logctx data from the context and adds it to the entry.
type LogCtxHook struct{}

// Levels returns all levels.
func (LogCtxHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var _ logrus.Hook = (*LogCtxHook)(nil)

// Fire is invoked for every entry and expands the logctx into the Data map.
func (LogCtxHook) Fire(e *logrus.Entry) error {
	ctx := e.Context
	if e.Context == nil {
		return nil
	}

	for _, v := range Get(ctx) {
		// avoid key collision, existing entries take precedence
		k := v.Key
		for {
			_, has := e.Data[k]
			if !has {
				break
			}
			k = "_" + k
		}
		// set the context one
		e.Data[k] = v.Value
	}

	return nil
}
