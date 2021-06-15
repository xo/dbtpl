// Command xo generates code from database schemas and custom queries. Works
// with PostgreSQL, MySQL, Microsoft SQL Server, Oracle Database, and SQLite3.
package main

//go:generate ./gen.sh models

import (
	"context"
	"fmt"
	"os"

	// drivers
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/sijms/go-ora"

	// templates
	_ "github.com/mmmcorp/xo/templates/gotpl"
	_ "github.com/mmmcorp/xo/templates/graphviztpl"

	"github.com/mmmcorp/xo/cmd"
)

// version is the app version.
var version = "0.0.0-dev"

func main() {
	if err := cmd.Run(context.Background(), "xo", version); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
