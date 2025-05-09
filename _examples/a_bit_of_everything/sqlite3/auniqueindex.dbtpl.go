package sqlite3

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// AUniqueIndex represents a row from 'a_unique_index'.
type AUniqueIndex struct {
	AKey sql.NullInt64 `json:"a_key"` // a_key
}

// AUniqueIndexByAKey retrieves a row from 'a_unique_index' as a [AUniqueIndex].
//
// Generated from index 'sqlite_autoindex_a_unique_index_1'.
func AUniqueIndexByAKey(ctx context.Context, db DB, aKey sql.NullInt64) (*AUniqueIndex, error) {
	// query
	const sqlstr = `SELECT ` +
		`a_key ` +
		`FROM a_unique_index ` +
		`WHERE a_key = $1`
	// run
	logf(sqlstr, aKey)
	aui := AUniqueIndex{}
	if err := db.QueryRowContext(ctx, sqlstr, aKey).Scan(&aui.AKey); err != nil {
		return nil, logerror(err)
	}
	return &aui, nil
}
