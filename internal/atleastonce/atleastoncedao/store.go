package atleastoncedao

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"taxes-be/internal/atleastonce"
	"taxes-be/internal/daoutil"
	"taxes-be/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Save(ctx context.Context, task atleastonce.Task) error {
	t := models.AtLeastOnceTask{
		Key:  task.Key,
		ID:   task.ID.String(),
		Done: task.Done,
		Time: task.Time,
	}
	err := daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		return t.Upsert(
			ctx,
			tx,
			true,
			[]string{
				models.AtLeastOnceTaskColumns.Key,
				models.AtLeastOnceTaskColumns.ID,
			},
			boil.Infer(),
			boil.Infer(),
		)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Find(ctx context.Context, key string, id uuid.UUID) (task *atleastonce.Task, err error) {
	err = daoutil.EnsureTransaction(ctx, s.db, func(tx *sql.Tx) error {
		var t *models.AtLeastOnceTask
		t, err = models.FindAtLeastOnceTask(ctx, tx, key, id.String())
		if err != nil {
			return err
		}
		task, err = toDomainObject(t)
		return err
	})
	if err != nil {
		return nil, err
	}

	return task, err
}

func (s *Store) FindAllToDo(ctx context.Context, limit int) ([]atleastonce.Task, error) {
	tasks, err := models.AtLeastOnceTasks(
		qm.Where("done=?", false),
		qm.Limit(limit),
	).All(ctx, s.db)
	if err != nil {
		return nil, err
	}

	result := make([]atleastonce.Task, len(tasks))
	for i, t := range tasks {
		dt, err := toDomainObject(t)
		if err != nil {
			return nil, err
		}
		result[i] = *dt
	}

	return result, nil
}

func (s *Store) InTransaction(ctx context.Context, retryableFn func(ctx context.Context) error) error {
	return daoutil.InTransaction(ctx, s.db, retryableFn)
}

func toDomainObject(t *models.AtLeastOnceTask) (*atleastonce.Task, error) {
	readID, err := uuid.Parse(t.ID)
	if err != nil {
		return nil, err
	}
	return &atleastonce.Task{
		Key:  t.Key,
		ID:   readID,
		Time: t.Time,
		Done: t.Done,
	}, nil
}
