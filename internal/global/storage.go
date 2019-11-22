package global

import (
	// internal golang package

	"context"
	"database/sql"
	"strconv"

	// internal package
	"phonebook/config"

	// thirdparty package
	db "phonebook/pkg/queryable"

	"github.com/jmoiron/sqlx"
)

// DB [singleton] global instance of sqlx
func DB() *sqlx.DB {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}
	con, err := sql.Open(cfg.DBType, cfg.DBConnectionString)
	if err != nil {
		panic(err)
	}

	openConn, err := strconv.Atoi(cfg.MaxOpenCon)
	if err != nil {
		panic(err)
	}

	idleConn, err := strconv.Atoi(cfg.MaxIdleCon)
	if err != nil {
		panic(err)
	}

	con.SetMaxIdleConns(idleConn)
	con.SetMaxOpenConns(openConn)

	sqlxDB := sqlx.NewDb(con, cfg.DBType)
	return sqlxDB
}

func GetQuery(ctx context.Context) *db.Queryable {
	q, ok := db.QueryableFromContext(ctx)
	if !ok {
		qctx := db.NewQueryable(DB())
		return &qctx
	}

	return &q
}
