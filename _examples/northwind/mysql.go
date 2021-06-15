package main

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/mmmcorp/xo/_examples/northwind/mysql"
)

func runMysql(ctx context.Context, db *sql.DB) error {
	p, err := models.ProductByProductID(ctx, db, 16)
	if err != nil {
		return err
	}
	fmt.Printf("product %d: %q\n", p.ProductID, p.ProductName)
	return nil
}
