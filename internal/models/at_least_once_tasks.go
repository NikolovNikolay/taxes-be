// Code generated by SQLBoiler 4.4.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// AtLeastOnceTask is an object representing the database table.
type AtLeastOnceTask struct {
	Key  string    `boil:"key" json:"key" toml:"key" yaml:"key"`
	ID   string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Done bool      `boil:"done" json:"done" toml:"done" yaml:"done"`
	Time time.Time `boil:"time" json:"time" toml:"time" yaml:"time"`

	R *atLeastOnceTaskR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L atLeastOnceTaskL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var AtLeastOnceTaskColumns = struct {
	Key  string
	ID   string
	Done string
	Time string
}{
	Key:  "key",
	ID:   "id",
	Done: "done",
	Time: "time",
}

// Generated where

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var AtLeastOnceTaskWhere = struct {
	Key  whereHelperstring
	ID   whereHelperstring
	Done whereHelperbool
	Time whereHelpertime_Time
}{
	Key:  whereHelperstring{field: "\"at_least_once_tasks\".\"key\""},
	ID:   whereHelperstring{field: "\"at_least_once_tasks\".\"id\""},
	Done: whereHelperbool{field: "\"at_least_once_tasks\".\"done\""},
	Time: whereHelpertime_Time{field: "\"at_least_once_tasks\".\"time\""},
}

// AtLeastOnceTaskRels is where relationship names are stored.
var AtLeastOnceTaskRels = struct {
}{}

// atLeastOnceTaskR is where relationships are stored.
type atLeastOnceTaskR struct {
}

// NewStruct creates a new relationship struct
func (*atLeastOnceTaskR) NewStruct() *atLeastOnceTaskR {
	return &atLeastOnceTaskR{}
}

// atLeastOnceTaskL is where Load methods for each relationship are stored.
type atLeastOnceTaskL struct{}

var (
	atLeastOnceTaskAllColumns            = []string{"key", "id", "done", "time"}
	atLeastOnceTaskColumnsWithoutDefault = []string{"key", "id"}
	atLeastOnceTaskColumnsWithDefault    = []string{"done", "time"}
	atLeastOnceTaskPrimaryKeyColumns     = []string{"key", "id"}
)

type (
	// AtLeastOnceTaskSlice is an alias for a slice of pointers to AtLeastOnceTask.
	// This should generally be used opposed to []AtLeastOnceTask.
	AtLeastOnceTaskSlice []*AtLeastOnceTask
	// AtLeastOnceTaskHook is the signature for custom AtLeastOnceTask hook methods
	AtLeastOnceTaskHook func(context.Context, boil.ContextExecutor, *AtLeastOnceTask) error

	atLeastOnceTaskQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	atLeastOnceTaskType                 = reflect.TypeOf(&AtLeastOnceTask{})
	atLeastOnceTaskMapping              = queries.MakeStructMapping(atLeastOnceTaskType)
	atLeastOnceTaskPrimaryKeyMapping, _ = queries.BindMapping(atLeastOnceTaskType, atLeastOnceTaskMapping, atLeastOnceTaskPrimaryKeyColumns)
	atLeastOnceTaskInsertCacheMut       sync.RWMutex
	atLeastOnceTaskInsertCache          = make(map[string]insertCache)
	atLeastOnceTaskUpdateCacheMut       sync.RWMutex
	atLeastOnceTaskUpdateCache          = make(map[string]updateCache)
	atLeastOnceTaskUpsertCacheMut       sync.RWMutex
	atLeastOnceTaskUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var atLeastOnceTaskBeforeInsertHooks []AtLeastOnceTaskHook
var atLeastOnceTaskBeforeUpdateHooks []AtLeastOnceTaskHook
var atLeastOnceTaskBeforeDeleteHooks []AtLeastOnceTaskHook
var atLeastOnceTaskBeforeUpsertHooks []AtLeastOnceTaskHook

var atLeastOnceTaskAfterInsertHooks []AtLeastOnceTaskHook
var atLeastOnceTaskAfterSelectHooks []AtLeastOnceTaskHook
var atLeastOnceTaskAfterUpdateHooks []AtLeastOnceTaskHook
var atLeastOnceTaskAfterDeleteHooks []AtLeastOnceTaskHook
var atLeastOnceTaskAfterUpsertHooks []AtLeastOnceTaskHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *AtLeastOnceTask) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *AtLeastOnceTask) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *AtLeastOnceTask) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *AtLeastOnceTask) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *AtLeastOnceTask) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *AtLeastOnceTask) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *AtLeastOnceTask) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *AtLeastOnceTask) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *AtLeastOnceTask) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range atLeastOnceTaskAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddAtLeastOnceTaskHook registers your hook function for all future operations.
func AddAtLeastOnceTaskHook(hookPoint boil.HookPoint, atLeastOnceTaskHook AtLeastOnceTaskHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		atLeastOnceTaskBeforeInsertHooks = append(atLeastOnceTaskBeforeInsertHooks, atLeastOnceTaskHook)
	case boil.BeforeUpdateHook:
		atLeastOnceTaskBeforeUpdateHooks = append(atLeastOnceTaskBeforeUpdateHooks, atLeastOnceTaskHook)
	case boil.BeforeDeleteHook:
		atLeastOnceTaskBeforeDeleteHooks = append(atLeastOnceTaskBeforeDeleteHooks, atLeastOnceTaskHook)
	case boil.BeforeUpsertHook:
		atLeastOnceTaskBeforeUpsertHooks = append(atLeastOnceTaskBeforeUpsertHooks, atLeastOnceTaskHook)
	case boil.AfterInsertHook:
		atLeastOnceTaskAfterInsertHooks = append(atLeastOnceTaskAfterInsertHooks, atLeastOnceTaskHook)
	case boil.AfterSelectHook:
		atLeastOnceTaskAfterSelectHooks = append(atLeastOnceTaskAfterSelectHooks, atLeastOnceTaskHook)
	case boil.AfterUpdateHook:
		atLeastOnceTaskAfterUpdateHooks = append(atLeastOnceTaskAfterUpdateHooks, atLeastOnceTaskHook)
	case boil.AfterDeleteHook:
		atLeastOnceTaskAfterDeleteHooks = append(atLeastOnceTaskAfterDeleteHooks, atLeastOnceTaskHook)
	case boil.AfterUpsertHook:
		atLeastOnceTaskAfterUpsertHooks = append(atLeastOnceTaskAfterUpsertHooks, atLeastOnceTaskHook)
	}
}

// One returns a single atLeastOnceTask record from the query.
func (q atLeastOnceTaskQuery) One(ctx context.Context, exec boil.ContextExecutor) (*AtLeastOnceTask, error) {
	o := &AtLeastOnceTask{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for at_least_once_tasks")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all AtLeastOnceTask records from the query.
func (q atLeastOnceTaskQuery) All(ctx context.Context, exec boil.ContextExecutor) (AtLeastOnceTaskSlice, error) {
	var o []*AtLeastOnceTask

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to AtLeastOnceTask slice")
	}

	if len(atLeastOnceTaskAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all AtLeastOnceTask records in the query.
func (q atLeastOnceTaskQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count at_least_once_tasks rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q atLeastOnceTaskQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if at_least_once_tasks exists")
	}

	return count > 0, nil
}

// AtLeastOnceTasks retrieves all the records using an executor.
func AtLeastOnceTasks(mods ...qm.QueryMod) atLeastOnceTaskQuery {
	mods = append(mods, qm.From("\"at_least_once_tasks\""))
	return atLeastOnceTaskQuery{NewQuery(mods...)}
}

// FindAtLeastOnceTask retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindAtLeastOnceTask(ctx context.Context, exec boil.ContextExecutor, key string, iD string, selectCols ...string) (*AtLeastOnceTask, error) {
	atLeastOnceTaskObj := &AtLeastOnceTask{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"at_least_once_tasks\" where \"key\"=$1 AND \"id\"=$2", sel,
	)

	q := queries.Raw(query, key, iD)

	err := q.Bind(ctx, exec, atLeastOnceTaskObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from at_least_once_tasks")
	}

	return atLeastOnceTaskObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *AtLeastOnceTask) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no at_least_once_tasks provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(atLeastOnceTaskColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	atLeastOnceTaskInsertCacheMut.RLock()
	cache, cached := atLeastOnceTaskInsertCache[key]
	atLeastOnceTaskInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			atLeastOnceTaskAllColumns,
			atLeastOnceTaskColumnsWithDefault,
			atLeastOnceTaskColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(atLeastOnceTaskType, atLeastOnceTaskMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(atLeastOnceTaskType, atLeastOnceTaskMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"at_least_once_tasks\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"at_least_once_tasks\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into at_least_once_tasks")
	}

	if !cached {
		atLeastOnceTaskInsertCacheMut.Lock()
		atLeastOnceTaskInsertCache[key] = cache
		atLeastOnceTaskInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the AtLeastOnceTask.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *AtLeastOnceTask) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	atLeastOnceTaskUpdateCacheMut.RLock()
	cache, cached := atLeastOnceTaskUpdateCache[key]
	atLeastOnceTaskUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			atLeastOnceTaskAllColumns,
			atLeastOnceTaskPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update at_least_once_tasks, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"at_least_once_tasks\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, atLeastOnceTaskPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(atLeastOnceTaskType, atLeastOnceTaskMapping, append(wl, atLeastOnceTaskPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update at_least_once_tasks row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for at_least_once_tasks")
	}

	if !cached {
		atLeastOnceTaskUpdateCacheMut.Lock()
		atLeastOnceTaskUpdateCache[key] = cache
		atLeastOnceTaskUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q atLeastOnceTaskQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for at_least_once_tasks")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for at_least_once_tasks")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o AtLeastOnceTaskSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), atLeastOnceTaskPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"at_least_once_tasks\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, atLeastOnceTaskPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in atLeastOnceTask slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all atLeastOnceTask")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *AtLeastOnceTask) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no at_least_once_tasks provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(atLeastOnceTaskColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	atLeastOnceTaskUpsertCacheMut.RLock()
	cache, cached := atLeastOnceTaskUpsertCache[key]
	atLeastOnceTaskUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			atLeastOnceTaskAllColumns,
			atLeastOnceTaskColumnsWithDefault,
			atLeastOnceTaskColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			atLeastOnceTaskAllColumns,
			atLeastOnceTaskPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert at_least_once_tasks, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(atLeastOnceTaskPrimaryKeyColumns))
			copy(conflict, atLeastOnceTaskPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"at_least_once_tasks\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(atLeastOnceTaskType, atLeastOnceTaskMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(atLeastOnceTaskType, atLeastOnceTaskMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert at_least_once_tasks")
	}

	if !cached {
		atLeastOnceTaskUpsertCacheMut.Lock()
		atLeastOnceTaskUpsertCache[key] = cache
		atLeastOnceTaskUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single AtLeastOnceTask record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *AtLeastOnceTask) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no AtLeastOnceTask provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), atLeastOnceTaskPrimaryKeyMapping)
	sql := "DELETE FROM \"at_least_once_tasks\" WHERE \"key\"=$1 AND \"id\"=$2"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from at_least_once_tasks")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for at_least_once_tasks")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q atLeastOnceTaskQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no atLeastOnceTaskQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from at_least_once_tasks")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for at_least_once_tasks")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o AtLeastOnceTaskSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(atLeastOnceTaskBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), atLeastOnceTaskPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"at_least_once_tasks\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, atLeastOnceTaskPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from atLeastOnceTask slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for at_least_once_tasks")
	}

	if len(atLeastOnceTaskAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *AtLeastOnceTask) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindAtLeastOnceTask(ctx, exec, o.Key, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *AtLeastOnceTaskSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := AtLeastOnceTaskSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), atLeastOnceTaskPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"at_least_once_tasks\".* FROM \"at_least_once_tasks\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, atLeastOnceTaskPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in AtLeastOnceTaskSlice")
	}

	*o = slice

	return nil
}

// AtLeastOnceTaskExists checks if the AtLeastOnceTask row exists.
func AtLeastOnceTaskExists(ctx context.Context, exec boil.ContextExecutor, key string, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"at_least_once_tasks\" where \"key\"=$1 AND \"id\"=$2 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, key, iD)
	}
	row := exec.QueryRowContext(ctx, sql, key, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if at_least_once_tasks exists")
	}

	return exists, nil
}
