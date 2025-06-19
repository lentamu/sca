package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const timeout = 10 * time.Second

func Connect(username, password, host, database string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return sqlx.ConnectContext(ctx, "mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, database))
}

func IsDuplicate(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
