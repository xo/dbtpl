package oracle

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// AuthPermission represents a row from 'django.auth_permission'.
type AuthPermission struct {
	ID            int64          `json:"id"`              // id
	Name          sql.NullString `json:"name"`            // name
	ContentTypeID int64          `json:"content_type_id"` // content_type_id
	Codename      sql.NullString `json:"codename"`        // codename
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [AuthPermission] exists in the database.
func (ap *AuthPermission) Exists() bool {
	return ap._exists
}

// Deleted returns true when the [AuthPermission] has been marked for deletion
// from the database.
func (ap *AuthPermission) Deleted() bool {
	return ap._deleted
}

// Insert inserts the [AuthPermission] to the database.
func (ap *AuthPermission) Insert(ctx context.Context, db DB) error {
	switch {
	case ap._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case ap._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO django.auth_permission (` +
		`name, content_type_id, codename` +
		`) VALUES (` +
		`:1, :2, :3` +
		`) RETURNING id INTO :4`
	// run
	logf(sqlstr, ap.Name, ap.ContentTypeID, ap.Codename)
	var id int64
	if _, err := db.ExecContext(ctx, sqlstr, ap.Name, ap.ContentTypeID, ap.Codename, sql.Out{Dest: &id}); err != nil {
		return logerror(err)
	} // set primary key
	ap.ID = int64(id)
	// set exists
	ap._exists = true
	return nil
}

// Update updates a [AuthPermission] in the database.
func (ap *AuthPermission) Update(ctx context.Context, db DB) error {
	switch {
	case !ap._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case ap._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE django.auth_permission SET ` +
		`name = :1, content_type_id = :2, codename = :3 ` +
		`WHERE id = :4`
	// run
	logf(sqlstr, ap.Name, ap.ContentTypeID, ap.Codename, ap.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ap.Name, ap.ContentTypeID, ap.Codename, ap.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [AuthPermission] to the database.
func (ap *AuthPermission) Save(ctx context.Context, db DB) error {
	if ap.Exists() {
		return ap.Update(ctx, db)
	}
	return ap.Insert(ctx, db)
}

// Upsert performs an upsert for [AuthPermission].
func (ap *AuthPermission) Upsert(ctx context.Context, db DB) error {
	switch {
	case ap._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `MERGE django.auth_permissiont ` +
		`USING (` +
		`SELECT :1 id, :2 name, :3 content_type_id, :4 codename ` +
		`FROM DUAL ) s ` +
		`ON s.id = t.id ` +
		`WHEN MATCHED THEN ` +
		`UPDATE SET ` +
		`t.name = s.name, t.content_type_id = s.content_type_id, t.codename = s.codename ` +
		`WHEN NOT MATCHED THEN ` +
		`INSERT (` +
		`name, content_type_id, codename` +
		`) VALUES (` +
		`s.name, s.content_type_id, s.codename` +
		`);`
	// run
	logf(sqlstr, ap.ID, ap.Name, ap.ContentTypeID, ap.Codename)
	if _, err := db.ExecContext(ctx, sqlstr, ap.ID, ap.Name, ap.ContentTypeID, ap.Codename); err != nil {
		return logerror(err)
	}
	// set exists
	ap._exists = true
	return nil
}

// Delete deletes the [AuthPermission] from the database.
func (ap *AuthPermission) Delete(ctx context.Context, db DB) error {
	switch {
	case !ap._exists: // doesn't exist
		return nil
	case ap._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM django.auth_permission ` +
		`WHERE id = :1`
	// run
	logf(sqlstr, ap.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ap.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	ap._deleted = true
	return nil
}

// AuthPermissionByContentTypeIDCodename retrieves a row from 'django.auth_permission' as a [AuthPermission].
//
// Generated from index 'auth_perm_content_t_01ab375a_u'.
func AuthPermissionByContentTypeIDCodename(ctx context.Context, db DB, contentTypeID int64, codename sql.NullString) (*AuthPermission, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, content_type_id, codename ` +
		`FROM django.auth_permission ` +
		`WHERE content_type_id = :1 AND codename = :2`
	// run
	logf(sqlstr, contentTypeID, codename)
	ap := AuthPermission{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, contentTypeID, codename).Scan(&ap.ID, &ap.Name, &ap.ContentTypeID, &ap.Codename); err != nil {
		return nil, logerror(err)
	}
	return &ap, nil
}

// AuthPermissionByContentTypeID retrieves a row from 'django.auth_permission' as a [AuthPermission].
//
// Generated from index 'auth_permi_content_ty_2f476e4b'.
func AuthPermissionByContentTypeID(ctx context.Context, db DB, contentTypeID int64) ([]*AuthPermission, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, content_type_id, codename ` +
		`FROM django.auth_permission ` +
		`WHERE content_type_id = :1`
	// run
	logf(sqlstr, contentTypeID)
	rows, err := db.QueryContext(ctx, sqlstr, contentTypeID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*AuthPermission
	for rows.Next() {
		ap := AuthPermission{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&ap.ID, &ap.Name, &ap.ContentTypeID, &ap.Codename); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &ap)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// AuthPermissionByID retrieves a row from 'django.auth_permission' as a [AuthPermission].
//
// Generated from index 'auth_permission_id_idx'.
func AuthPermissionByID(ctx context.Context, db DB, id int64) (*AuthPermission, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, content_type_id, codename ` +
		`FROM django.auth_permission ` +
		`WHERE id = :1`
	// run
	logf(sqlstr, id)
	ap := AuthPermission{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&ap.ID, &ap.Name, &ap.ContentTypeID, &ap.Codename); err != nil {
		return nil, logerror(err)
	}
	return &ap, nil
}

// DjangoContentType returns the DjangoContentType associated with the [AuthPermission]'s (ContentTypeID).
//
// Generated from foreign key 'auth_perm_content_t_2f476e4b_f'.
func (ap *AuthPermission) DjangoContentType(ctx context.Context, db DB) (*DjangoContentType, error) {
	return DjangoContentTypeByID(ctx, db, ap.ContentTypeID)
}
