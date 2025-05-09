package mysql

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
	"time"
)

// AuthUser represents a row from 'django.auth_user'.
type AuthUser struct {
	ID          int          `json:"id"`           // id
	Password    string       `json:"password"`     // password
	LastLogin   sql.NullTime `json:"last_login"`   // last_login
	IsSuperuser bool         `json:"is_superuser"` // is_superuser
	Username    string       `json:"username"`     // username
	FirstName   string       `json:"first_name"`   // first_name
	LastName    string       `json:"last_name"`    // last_name
	Email       string       `json:"email"`        // email
	IsStaff     bool         `json:"is_staff"`     // is_staff
	IsActive    bool         `json:"is_active"`    // is_active
	DateJoined  time.Time    `json:"date_joined"`  // date_joined
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [AuthUser] exists in the database.
func (au *AuthUser) Exists() bool {
	return au._exists
}

// Deleted returns true when the [AuthUser] has been marked for deletion
// from the database.
func (au *AuthUser) Deleted() bool {
	return au._deleted
}

// Insert inserts the [AuthUser] to the database.
func (au *AuthUser) Insert(ctx context.Context, db DB) error {
	switch {
	case au._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case au._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO django.auth_user (` +
		`password, last_login, is_superuser, username, first_name, last_name, email, is_staff, is_active, date_joined` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, au.Password, au.LastLogin, au.IsSuperuser, au.Username, au.FirstName, au.LastName, au.Email, au.IsStaff, au.IsActive, au.DateJoined)
	res, err := db.ExecContext(ctx, sqlstr, au.Password, au.LastLogin, au.IsSuperuser, au.Username, au.FirstName, au.LastName, au.Email, au.IsStaff, au.IsActive, au.DateJoined)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	au.ID = int(id)
	// set exists
	au._exists = true
	return nil
}

// Update updates a [AuthUser] in the database.
func (au *AuthUser) Update(ctx context.Context, db DB) error {
	switch {
	case !au._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case au._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE django.auth_user SET ` +
		`password = ?, last_login = ?, is_superuser = ?, username = ?, first_name = ?, last_name = ?, email = ?, is_staff = ?, is_active = ?, date_joined = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, au.Password, au.LastLogin, au.IsSuperuser, au.Username, au.FirstName, au.LastName, au.Email, au.IsStaff, au.IsActive, au.DateJoined, au.ID)
	if _, err := db.ExecContext(ctx, sqlstr, au.Password, au.LastLogin, au.IsSuperuser, au.Username, au.FirstName, au.LastName, au.Email, au.IsStaff, au.IsActive, au.DateJoined, au.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [AuthUser] to the database.
func (au *AuthUser) Save(ctx context.Context, db DB) error {
	if au.Exists() {
		return au.Update(ctx, db)
	}
	return au.Insert(ctx, db)
}

// Upsert performs an upsert for [AuthUser].
func (au *AuthUser) Upsert(ctx context.Context, db DB) error {
	switch {
	case au._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO django.auth_user (` +
		`id, password, last_login, is_superuser, username, first_name, last_name, email, is_staff, is_active, date_joined` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`password = VALUES(password), last_login = VALUES(last_login), is_superuser = VALUES(is_superuser), username = VALUES(username), first_name = VALUES(first_name), last_name = VALUES(last_name), email = VALUES(email), is_staff = VALUES(is_staff), is_active = VALUES(is_active), date_joined = VALUES(date_joined)`
	// run
	logf(sqlstr, au.ID, au.Password, au.LastLogin, au.IsSuperuser, au.Username, au.FirstName, au.LastName, au.Email, au.IsStaff, au.IsActive, au.DateJoined)
	if _, err := db.ExecContext(ctx, sqlstr, au.ID, au.Password, au.LastLogin, au.IsSuperuser, au.Username, au.FirstName, au.LastName, au.Email, au.IsStaff, au.IsActive, au.DateJoined); err != nil {
		return logerror(err)
	}
	// set exists
	au._exists = true
	return nil
}

// Delete deletes the [AuthUser] from the database.
func (au *AuthUser) Delete(ctx context.Context, db DB) error {
	switch {
	case !au._exists: // doesn't exist
		return nil
	case au._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM django.auth_user ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, au.ID)
	if _, err := db.ExecContext(ctx, sqlstr, au.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	au._deleted = true
	return nil
}

// AuthUserByID retrieves a row from 'django.auth_user' as a [AuthUser].
//
// Generated from index 'auth_user_id_pkey'.
func AuthUserByID(ctx context.Context, db DB, id int) (*AuthUser, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, password, last_login, is_superuser, username, first_name, last_name, email, is_staff, is_active, date_joined ` +
		`FROM django.auth_user ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	au := AuthUser{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&au.ID, &au.Password, &au.LastLogin, &au.IsSuperuser, &au.Username, &au.FirstName, &au.LastName, &au.Email, &au.IsStaff, &au.IsActive, &au.DateJoined); err != nil {
		return nil, logerror(err)
	}
	return &au, nil
}

// AuthUserByUsername retrieves a row from 'django.auth_user' as a [AuthUser].
//
// Generated from index 'username'.
func AuthUserByUsername(ctx context.Context, db DB, username string) (*AuthUser, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, password, last_login, is_superuser, username, first_name, last_name, email, is_staff, is_active, date_joined ` +
		`FROM django.auth_user ` +
		`WHERE username = ?`
	// run
	logf(sqlstr, username)
	au := AuthUser{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, username).Scan(&au.ID, &au.Password, &au.LastLogin, &au.IsSuperuser, &au.Username, &au.FirstName, &au.LastName, &au.Email, &au.IsStaff, &au.IsActive, &au.DateJoined); err != nil {
		return nil, logerror(err)
	}
	return &au, nil
}
