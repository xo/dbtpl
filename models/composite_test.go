package models

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"
)

type queryExpectation struct {
	contains []string
	argsLen  int
}

type mockConnector struct {
	t           *testing.T
	expectation queryExpectation
	columns     []string
	rows        [][]driver.Value
}

type mockConn struct {
	connector *mockConnector
}

type mockRows struct {
	columns []string
	rows    [][]driver.Value
	idx     int
}

func (c *mockConnector) Connect(_ context.Context) (driver.Conn, error) {
	return &mockConn{connector: c}, nil
}

func (c *mockConnector) Driver() driver.Driver {
	return &mockDriver{connector: c}
}

type mockDriver struct {
	connector *mockConnector
}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{connector: d.connector}, nil
}

func (c *mockConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("not implemented")
}

func (c *mockConn) Close() error { return nil }

func (c *mockConn) Begin() (driver.Tx, error) {
	return nil, errors.New("not implemented")
}

func (c *mockConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	for _, fragment := range c.connector.expectation.contains {
		if !strings.Contains(query, fragment) {
			c.connector.t.Fatalf("query missing expected fragment %q: %s", fragment, query)
		}
	}
	if c.connector.expectation.argsLen != 0 && len(args) != c.connector.expectation.argsLen {
		c.connector.t.Fatalf("expected %d args, got %d", c.connector.expectation.argsLen, len(args))
	}
	return &mockRows{columns: c.connector.columns, rows: c.connector.rows}, nil
}

func (r *mockRows) Columns() []string { return r.columns }

func (r *mockRows) Close() error { return nil }

func (r *mockRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.idx])
	r.idx++
	return nil
}

func TestPostgresCompositesFiltersTables(t *testing.T) {
	connector := &mockConnector{
		t: t,
		expectation: queryExpectation{
			contains: []string{"t.typtype = 'c'", "n.nspname = $1", "c.relkind = 'c'"},
			argsLen:  1,
		},
		columns: []string{"type_name", "comment"},
		rows:    [][]driver.Value{{"address", ""}},
	}
	db := sql.OpenDB(connector)
	composites, err := PostgresComposites(context.Background(), db, "public")
	if err != nil {
		t.Fatalf("PostgresComposites unexpected error: %v", err)
	}
	if len(composites) != 1 {
		t.Fatalf("expected 1 composite, got %d", len(composites))
	}
	if composites[0].TypeName != "address" {
		t.Fatalf("unexpected composite name: %s", composites[0].TypeName)
	}
}

func TestPostgresCompositeAttributesFiltersTables(t *testing.T) {
	connector := &mockConnector{
		t: t,
		expectation: queryExpectation{
			contains: []string{"n.nspname = $1", "t.typname = $2", "c.relkind = 'c'"},
			argsLen:  2,
		},
		columns: []string{"field_ordinal", "attribute", "data_type", "not_null", "default_value", "comment"},
		rows:    [][]driver.Value{{int64(1), "street", "text", true, "", ""}},
	}
	db := sql.OpenDB(connector)
	attrs, err := PostgresCompositeAttributes(context.Background(), db, "public", "address")
	if err != nil {
		t.Fatalf("PostgresCompositeAttributes unexpected error: %v", err)
	}
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(attrs))
	}
	if attrs[0].Attribute != "street" {
		t.Fatalf("unexpected attribute name: %s", attrs[0].Attribute)
	}
}
