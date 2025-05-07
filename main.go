// Command dbtpl generates code from database schemas and custom queries. Works
// with PostgreSQL, MySQL, Microsoft SQL Server, Oracle Database, and SQLite3.
//
//go:debug x509negativeserial=1
package main

//go:generate ./gen.sh models
//go:generate go generate ./internal

import (
	"context"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"

	"github.com/xo/dbtpl/cmd"
)

func main() {
	cmd.Run(context.Background(), "dbtpl")
}
