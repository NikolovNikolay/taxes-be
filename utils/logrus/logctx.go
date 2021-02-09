package logrus

// Package logctx provides methods to add and retrieve logging information
// from a context.
//
// It's main use is to store logging data between method invocations in a
// given request and pass it to the logger after the request's response is
// provided. This is achieved in conjunction with the
// git.ridewithvia.dev/golib/logrusctx package.

import "context"

type ctxkey struct{}

// nolint:gochecknoglobals // keyInst is a key for context storage. It is a constant key, unreproducible in other pkg.
var keyInst ctxkey

// LogItem is the value written in the log.
type LogItem struct {
	Key   string
	Value interface{}
}

// With is a shorthand for the creation of LogItem.
// To be used in combination with Extend.
func With(key string, value interface{}) LogItem {
	return LogItem{
		Key:   key,
		Value: value,
	}
}

// Extend is used to add logging information to the context.
// Usage:
// 	ctx = logctx.Extend(ctx,
// 	 	logctx.With("parent", "John"),
// 	 	logctx.With("child", "George"),
// 	 	logctx.With("age", 11),
// 		logctx.With("httpRequest", req),
// 	)
func Extend(ctx context.Context, items ...LogItem) context.Context {
	l := list{
		items: items,
	}
	v := ctx.Value(keyInst)
	l.prev, _ = v.(*list)
	return context.WithValue(ctx, keyInst, &l)
}

// Get retrieves all logctx items from the given context.
func Get(ctx context.Context) []LogItem {
	v := ctx.Value(keyInst)
	l, ok := v.(*list)
	if !ok {
		return nil
	}
	var res []LogItem
	for ; l != nil; l = l.prev {
		res = append(res, l.items...)
	}
	return res
}

// list represents a linked list of log items.
type list struct {
	prev  *list
	items []LogItem
}
