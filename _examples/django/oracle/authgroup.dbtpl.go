// Package oracle contains generated code for schema 'django'.
package oracle

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// AuthGroup represents a row from 'django.auth_group'.
type AuthGroup struct {
	ID   int64          `json:"id"`   // id
	Name sql.NullString `json:"name"` // name
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [AuthGroup] exists in the database.
func (ag *AuthGroup) Exists() bool {
	return ag._exists
}

// Deleted returns true when the [AuthGroup] has been marked for deletion
// from the database.
func (ag *AuthGroup) Deleted() bool {
	return ag._deleted
}

// Insert inserts the [AuthGroup] to the database.
func (ag *AuthGroup) Insert(ctx context.Context, db DB) error {
	switch {
	case ag._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case ag._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO django.auth_group (` +
		`name` +
		`) VALUES (` +
		`:1` +
		`) RETURNING id INTO :2`
	// run
	logf(sqlstr, ag.Name)
	var id int64
	if _, err := db.ExecContext(ctx, sqlstr, ag.Name, sql.Out{Dest: &id}); err != nil {
		return logerror(err)
	} // set primary key
	ag.ID = int64(id)
	// set exists
	ag._exists = true
	return nil
}

// Update updates a [AuthGroup] in the database.
func (ag *AuthGroup) Update(ctx context.Context, db DB) error {
	switch {
	case !ag._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case ag._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE django.auth_group SET ` +
		`name = :1 ` +
		`WHERE id = :2`
	// run
	logf(sqlstr, ag.Name, ag.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ag.Name, ag.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [AuthGroup] to the database.
func (ag *AuthGroup) Save(ctx context.Context, db DB) error {
	if ag.Exists() {
		return ag.Update(ctx, db)
	}
	return ag.Insert(ctx, db)
}

// Upsert performs an upsert for [AuthGroup].
func (ag *AuthGroup) Upsert(ctx context.Context, db DB) error {
	switch {
	case ag._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `MERGE django.auth_groupt ` +
		`USING (` +
		`SELECT :1 id, :2 name ` +
		`FROM DUAL ) s ` +
		`ON s.id = t.id ` +
		`WHEN MATCHED THEN ` +
		`UPDATE SET ` +
		`t.name = s.name ` +
		`WHEN NOT MATCHED THEN ` +
		`INSERT (` +
		`name` +
		`) VALUES (` +
		`s.name` +
		`);`
	// run
	logf(sqlstr, ag.ID, ag.Name)
	if _, err := db.ExecContext(ctx, sqlstr, ag.ID, ag.Name); err != nil {
		return logerror(err)
	}
	// set exists
	ag._exists = true
	return nil
}

// Delete deletes the [AuthGroup] from the database.
func (ag *AuthGroup) Delete(ctx context.Context, db DB) error {
	switch {
	case !ag._exists: // doesn't exist
		return nil
	case ag._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM django.auth_group ` +
		`WHERE id = :1`
	// run
	logf(sqlstr, ag.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ag.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	ag._deleted = true
	return nil
}

// AuthGroupByID retrieves a row from 'django.auth_group' as a [AuthGroup].
//
// Generated from index 'auth_group_id_idx'.
func AuthGroupByID(ctx context.Context, db DB, id int64) (*AuthGroup, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name ` +
		`FROM django.auth_group ` +
		`WHERE id = :1`
	// run
	logf(sqlstr, id)
	ag := AuthGroup{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&ag.ID, &ag.Name); err != nil {
		return nil, logerror(err)
	}
	return &ag, nil
}

// AuthGroupByName retrieves a row from 'django.auth_group' as a [AuthGroup].
//
// Generated from index 'auth_group_name_idx'.
func AuthGroupByName(ctx context.Context, db DB, name sql.NullString) (*AuthGroup, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name ` +
		`FROM django.auth_group ` +
		`WHERE name = :1`
	// run
	logf(sqlstr, name)
	ag := AuthGroup{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, name).Scan(&ag.ID, &ag.Name); err != nil {
		return nil, logerror(err)
	}
	return &ag, nil
}
