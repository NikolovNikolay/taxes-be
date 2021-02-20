package atleastonce

import (
	"context"
	"errors"
	"sync"
	"taxes-be/internal/core"
	util "taxes-be/utils"
	"taxes-be/utils/asynctx"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const defaultAloLimit = 15

type Task struct {
	Key  string
	ID   uuid.UUID
	Time time.Time
	Done bool
}

type Store interface {
	Save(ctx context.Context, task Task) error
	Find(ctx context.Context, key string, id uuid.UUID) (*Task, error)
	FindAllToDo(ctx context.Context, limit int) ([]Task, error)
	InTransaction(ctx context.Context, retryableFn func(ctx context.Context) error) error
}

func New(store Store) *Doer {
	return &Doer{
		store:    store,
		handlers: make(map[string]func(ctx context.Context, id uuid.UUID) error),
		delayed:  util.NewDelayedExecutor(20),
	}
}

type Doer struct {
	store    Store
	handlers map[string]func(ctx context.Context, id uuid.UUID) error
	delayed  *util.DelayedExecutor
}

func (d *Doer) Close() {
	d.delayed.Stop()
}

func (d *Doer) RegisterHandler(key string, handler func(ctx context.Context, id uuid.UUID) error) {
	d.handlers[key] = handler
}

func (d *Doer) ToDo(ctx context.Context, key string, id uuid.UUID, lastUpdated time.Time) error {
	return d.store.Save(ctx, Task{
		Key:  key,
		ID:   id,
		Time: lastUpdated,
		Done: false,
	})
}

var ErrHandlerNotRegistered = errors.New("at-least-once handler not registered for the given key")

func (d *Doer) Try(ctx context.Context, key string, id uuid.UUID) error {
	return d.store.InTransaction(ctx, func(ctx context.Context) error {
		task, err := d.store.Find(ctx, key, id)
		if err != nil {
			return err
		}
		if task.Done {
			return nil
		}
		h, ok := d.handlers[task.Key]
		if !ok {
			return ErrHandlerNotRegistered
		}
		err = h(ctx, task.ID)
		if err != nil {
			return err
		}
		task.Done = true
		return d.store.Save(ctx, *task)
	})
}

func (d *Doer) tryAndLog(ctx context.Context, key string, id uuid.UUID) {
	err := d.Try(ctx, key, id)
	if err != nil {
		if !core.IsNoLogError(err) {
			logrus.WithContext(ctx).
				WithField("key", key).
				WithField("id", id).
				WithError(err).
				Errorf("at-least-once error, key: %s, id: %v", key, id)
		}
	}
}

func (d *Doer) TryAsync(ctx context.Context, key string, id uuid.UUID) {
	d.delayed.Run(func() {
		// wait for the initial operation to finish to avoid conflicts over the same record in the DB
		<-ctx.Done()
		ctx = asynctx.Linked(ctx, "atleastonce.Doer.TryAsync")

		d.tryAndLog(ctx, key, id)
	})
}

func (d *Doer) TryAll(ctx context.Context, limit int) {
	if limit <= 0 {
		limit = defaultAloLimit
	}
	all, err := d.store.FindAllToDo(ctx, limit)
	if err != nil {
		logrus.WithContext(ctx).
			WithError(err).
			Error("failed finding items to try with")
	}
	wg := &sync.WaitGroup{}
	for _, task := range all {
		task := task
		wg.Add(1)
		d.delayed.Run(func() {
			defer wg.Done()
			d.tryAndLog(ctx, task.Key, task.ID)
		})
	}
	wg.Wait()
}
