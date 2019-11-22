package queryable

import (
	// internal golang package
	"context"
	"database/sql"
	"fmt"

	// thirdparty package
	"github.com/jmoiron/sqlx"
)

// Q ...
type Q interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	MustExec(query string, args ...interface{}) sql.Result
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Queryable struct {
	q  Q
	db *sqlx.DB
	tx *sqlx.Tx
}

type key int

const queryableKey key = 0

//NewQueryableContext ...
func NewQueryableContext(ctx context.Context, q Queryable) context.Context {
	return context.WithValue(ctx, queryableKey, q)
}

// QueryableFromContext ...
func QueryableFromContext(ctx context.Context) (Queryable, bool) {
	q, ok := ctx.Value(queryableKey).(Queryable)
	return q, ok
}

func NewQueryable(db interface{}) Queryable {
	var newQueryable Q
	var sqlxTx *sqlx.Tx
	var sqlxDB *sqlx.DB
	sqlxTx, ok := db.(*sqlx.Tx)
	if !ok {
		sqlxDB = db.(*sqlx.DB)
		newQueryable = Q(sqlxDB)
	} else {
		newQueryable = Q(sqlxTx)
	}
	return Queryable{
		q:  newQueryable,
		db: sqlxDB,
		tx: sqlxTx,
	}
}

func RunInTransaction(ctx context.Context, db *sqlx.DB, fn func(ctx context.Context) error) error {
	tx, err := db.Beginx()
	if nil != err {
		return err
	}

	ctx = NewQueryableContext(ctx, NewQueryable(tx))
	err = fn(ctx)

	if nil != err {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error when committing transaction: %v", err)
	}

	return nil
}
