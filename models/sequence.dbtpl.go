package models

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
)

// Sequence is a sequence.
type Sequence struct {
	ColumnName string `json:"column_name"` // column_name
}

// PostgresTableSequences runs a custom query, returning results as [Sequence].
func PostgresTableSequences(ctx context.Context, db DB, schema, table string) ([]*Sequence, error) {
	// query
	const sqlstr = `SELECT ` +
		`a.attname ` + // ::varchar as column_name
		`FROM pg_class s ` +
		`JOIN pg_depend d ON d.objid = s.oid ` +
		`JOIN pg_class t ON d.objid = s.oid AND d.refobjid = t.oid ` +
		`JOIN pg_attribute a ON (d.refobjid, d.refobjsubid) = (a.attrelid, a.attnum) ` +
		`JOIN pg_namespace n ON n.oid = s.relnamespace ` +
		`WHERE s.relkind = 'S' ` +
		`AND n.nspname = $1 ` +
		`AND t.relname = $2`
	// run
	logf(sqlstr, schema, table)
	rows, err := db.QueryContext(ctx, sqlstr, schema, table)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*Sequence
	for rows.Next() {
		var s Sequence
		// scan
		if err := rows.Scan(&s.ColumnName); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// MysqlTableSequences runs a custom query, returning results as [Sequence].
func MysqlTableSequences(ctx context.Context, db DB, schema, table string) ([]*Sequence, error) {
	// query
	const sqlstr = `SELECT ` +
		`column_name ` +
		`FROM information_schema.columns c ` +
		`WHERE c.extra = 'auto_increment' ` +
		`AND c.table_schema = ? ` +
		`AND c.table_name = ?`
	// run
	logf(sqlstr, schema, table)
	rows, err := db.QueryContext(ctx, sqlstr, schema, table)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*Sequence
	for rows.Next() {
		var s Sequence
		// scan
		if err := rows.Scan(&s.ColumnName); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// Sqlite3TableSequences runs a custom query, returning results as [Sequence].
func Sqlite3TableSequences(ctx context.Context, db DB, schema, table string) ([]*Sequence, error) {
	// query
	sqlstr := `/* ` + schema + ` */ ` +
		`WITH RECURSIVE ` +
		`a AS ( ` +
		`SELECT name, lower(replace(replace(sql, char(13), ' '), char(10), ' ')) AS sql ` +
		`FROM sqlite_master ` +
		`WHERE lower(sql) LIKE '%integer% autoincrement%' ` +
		`), ` +
		`b AS ( ` +
		`SELECT name, trim(substr(sql, instr(sql, '(') + 1)) AS sql ` +
		`FROM a ` +
		`), ` +
		`c AS ( ` +
		`SELECT b.name, sql, '' AS col ` +
		`FROM b ` +
		`UNION ALL ` +
		`SELECT ` +
		`c.name, ` +
		`trim(substr(c.sql, ifnull(nullif(instr(c.sql, ','), 0), instr(c.sql, ')')) + 1)) AS sql, ` +
		`trim(substr(c.sql, 1, ifnull(nullif(instr(c.sql, ','), 0), instr(c.sql, ')')) - 1)) AS col ` +
		`FROM c JOIN b ON c.name = b.name ` +
		`WHERE c.sql != '' ` +
		`), ` +
		`d AS ( ` +
		`SELECT name, substr(col, 1, instr(col, ' ') - 1) AS col ` +
		`FROM c ` +
		`WHERE col LIKE '%autoincrement%' ` +
		`) ` +
		`SELECT col AS column_name ` +
		`FROM d ` +
		`WHERE name = $1`
	// run
	logf(sqlstr, table)
	rows, err := db.QueryContext(ctx, sqlstr, table)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*Sequence
	for rows.Next() {
		var s Sequence
		// scan
		if err := rows.Scan(&s.ColumnName); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// SqlserverTableSequences runs a custom query, returning results as [Sequence].
func SqlserverTableSequences(ctx context.Context, db DB, schema, table string) ([]*Sequence, error) {
	// query
	const sqlstr = `SELECT ` +
		`COL_NAME(o.object_id, c.column_id) AS column_name ` +
		`FROM sys.objects o ` +
		`INNER JOIN sys.columns c ON o.object_id = c.object_id ` +
		`WHERE c.is_identity = 1 ` +
		`AND o.type = 'U' ` +
		`AND SCHEMA_NAME(o.schema_id) = @p1 ` +
		`AND o.name = @p2`
	// run
	logf(sqlstr, schema, table)
	rows, err := db.QueryContext(ctx, sqlstr, schema, table)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*Sequence
	for rows.Next() {
		var s Sequence
		// scan
		if err := rows.Scan(&s.ColumnName); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// OracleTableSequences runs a custom query, returning results as [Sequence].
func OracleTableSequences(ctx context.Context, db DB, schema, table string) ([]*Sequence, error) {
	// query
	const sqlstr = `SELECT ` +
		`LOWER(c.column_name) AS column_name ` +
		`FROM all_tab_columns c ` +
		`WHERE c.identity_column='YES' ` +
		`AND c.owner = UPPER(:1) ` +
		`AND c.table_name  = UPPER(:2)`
	// run
	logf(sqlstr, schema, table)
	rows, err := db.QueryContext(ctx, sqlstr, schema, table)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// load results
	var res []*Sequence
	for rows.Next() {
		var s Sequence
		// scan
		if err := rows.Scan(&s.ColumnName); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}
