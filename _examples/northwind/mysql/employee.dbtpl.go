package mysql

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Employee represents a row from 'northwind.employees'.
type Employee struct {
	EmployeeID      int16          `json:"employee_id"`       // employee_id
	LastName        string         `json:"last_name"`         // last_name
	FirstName       string         `json:"first_name"`        // first_name
	Title           sql.NullString `json:"title"`             // title
	TitleOfCourtesy sql.NullString `json:"title_of_courtesy"` // title_of_courtesy
	BirthDate       sql.NullTime   `json:"birth_date"`        // birth_date
	HireDate        sql.NullTime   `json:"hire_date"`         // hire_date
	Address         sql.NullString `json:"address"`           // address
	City            sql.NullString `json:"city"`              // city
	Region          sql.NullString `json:"region"`            // region
	PostalCode      sql.NullString `json:"postal_code"`       // postal_code
	Country         sql.NullString `json:"country"`           // country
	HomePhone       sql.NullString `json:"home_phone"`        // home_phone
	Extension       sql.NullString `json:"extension"`         // extension
	Photo           []byte         `json:"photo"`             // photo
	Notes           sql.NullString `json:"notes"`             // notes
	ReportsTo       sql.NullInt64  `json:"reports_to"`        // reports_to
	PhotoPath       sql.NullString `json:"photo_path"`        // photo_path
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Employee] exists in the database.
func (e *Employee) Exists() bool {
	return e._exists
}

// Deleted returns true when the [Employee] has been marked for deletion
// from the database.
func (e *Employee) Deleted() bool {
	return e._deleted
}

// Insert inserts the [Employee] to the database.
func (e *Employee) Insert(ctx context.Context, db DB) error {
	switch {
	case e._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case e._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO northwind.employees (` +
		`employee_id, last_name, first_name, title, title_of_courtesy, birth_date, hire_date, address, city, region, postal_code, country, home_phone, extension, photo, notes, reports_to, photo_path` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, e.EmployeeID, e.LastName, e.FirstName, e.Title, e.TitleOfCourtesy, e.BirthDate, e.HireDate, e.Address, e.City, e.Region, e.PostalCode, e.Country, e.HomePhone, e.Extension, e.Photo, e.Notes, e.ReportsTo, e.PhotoPath)
	if _, err := db.ExecContext(ctx, sqlstr, e.EmployeeID, e.LastName, e.FirstName, e.Title, e.TitleOfCourtesy, e.BirthDate, e.HireDate, e.Address, e.City, e.Region, e.PostalCode, e.Country, e.HomePhone, e.Extension, e.Photo, e.Notes, e.ReportsTo, e.PhotoPath); err != nil {
		return logerror(err)
	}
	// set exists
	e._exists = true
	return nil
}

// Update updates a [Employee] in the database.
func (e *Employee) Update(ctx context.Context, db DB) error {
	switch {
	case !e._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case e._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE northwind.employees SET ` +
		`last_name = ?, first_name = ?, title = ?, title_of_courtesy = ?, birth_date = ?, hire_date = ?, address = ?, city = ?, region = ?, postal_code = ?, country = ?, home_phone = ?, extension = ?, photo = ?, notes = ?, reports_to = ?, photo_path = ? ` +
		`WHERE employee_id = ?`
	// run
	logf(sqlstr, e.LastName, e.FirstName, e.Title, e.TitleOfCourtesy, e.BirthDate, e.HireDate, e.Address, e.City, e.Region, e.PostalCode, e.Country, e.HomePhone, e.Extension, e.Photo, e.Notes, e.ReportsTo, e.PhotoPath, e.EmployeeID)
	if _, err := db.ExecContext(ctx, sqlstr, e.LastName, e.FirstName, e.Title, e.TitleOfCourtesy, e.BirthDate, e.HireDate, e.Address, e.City, e.Region, e.PostalCode, e.Country, e.HomePhone, e.Extension, e.Photo, e.Notes, e.ReportsTo, e.PhotoPath, e.EmployeeID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Employee] to the database.
func (e *Employee) Save(ctx context.Context, db DB) error {
	if e.Exists() {
		return e.Update(ctx, db)
	}
	return e.Insert(ctx, db)
}

// Upsert performs an upsert for [Employee].
func (e *Employee) Upsert(ctx context.Context, db DB) error {
	switch {
	case e._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO northwind.employees (` +
		`employee_id, last_name, first_name, title, title_of_courtesy, birth_date, hire_date, address, city, region, postal_code, country, home_phone, extension, photo, notes, reports_to, photo_path` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`employee_id = VALUES(employee_id), last_name = VALUES(last_name), first_name = VALUES(first_name), title = VALUES(title), title_of_courtesy = VALUES(title_of_courtesy), birth_date = VALUES(birth_date), hire_date = VALUES(hire_date), address = VALUES(address), city = VALUES(city), region = VALUES(region), postal_code = VALUES(postal_code), country = VALUES(country), home_phone = VALUES(home_phone), extension = VALUES(extension), photo = VALUES(photo), notes = VALUES(notes), reports_to = VALUES(reports_to), photo_path = VALUES(photo_path)`
	// run
	logf(sqlstr, e.EmployeeID, e.LastName, e.FirstName, e.Title, e.TitleOfCourtesy, e.BirthDate, e.HireDate, e.Address, e.City, e.Region, e.PostalCode, e.Country, e.HomePhone, e.Extension, e.Photo, e.Notes, e.ReportsTo, e.PhotoPath)
	if _, err := db.ExecContext(ctx, sqlstr, e.EmployeeID, e.LastName, e.FirstName, e.Title, e.TitleOfCourtesy, e.BirthDate, e.HireDate, e.Address, e.City, e.Region, e.PostalCode, e.Country, e.HomePhone, e.Extension, e.Photo, e.Notes, e.ReportsTo, e.PhotoPath); err != nil {
		return logerror(err)
	}
	// set exists
	e._exists = true
	return nil
}

// Delete deletes the [Employee] from the database.
func (e *Employee) Delete(ctx context.Context, db DB) error {
	switch {
	case !e._exists: // doesn't exist
		return nil
	case e._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM northwind.employees ` +
		`WHERE employee_id = ?`
	// run
	logf(sqlstr, e.EmployeeID)
	if _, err := db.ExecContext(ctx, sqlstr, e.EmployeeID); err != nil {
		return logerror(err)
	}
	// set deleted
	e._deleted = true
	return nil
}

// EmployeeByEmployeeID retrieves a row from 'northwind.employees' as a [Employee].
//
// Generated from index 'employees_employee_id_pkey'.
func EmployeeByEmployeeID(ctx context.Context, db DB, employeeID int16) (*Employee, error) {
	// query
	const sqlstr = `SELECT ` +
		`employee_id, last_name, first_name, title, title_of_courtesy, birth_date, hire_date, address, city, region, postal_code, country, home_phone, extension, photo, notes, reports_to, photo_path ` +
		`FROM northwind.employees ` +
		`WHERE employee_id = ?`
	// run
	logf(sqlstr, employeeID)
	e := Employee{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, employeeID).Scan(&e.EmployeeID, &e.LastName, &e.FirstName, &e.Title, &e.TitleOfCourtesy, &e.BirthDate, &e.HireDate, &e.Address, &e.City, &e.Region, &e.PostalCode, &e.Country, &e.HomePhone, &e.Extension, &e.Photo, &e.Notes, &e.ReportsTo, &e.PhotoPath); err != nil {
		return nil, logerror(err)
	}
	return &e, nil
}

// EmployeesByReportsTo retrieves a row from 'northwind.employees' as a [Employee].
//
// Generated from index 'reports_to'.
func EmployeesByReportsTo(ctx context.Context, db DB, reportsTo sql.NullInt64) ([]*Employee, error) {
	// query
	const sqlstr = `SELECT ` +
		`employee_id, last_name, first_name, title, title_of_courtesy, birth_date, hire_date, address, city, region, postal_code, country, home_phone, extension, photo, notes, reports_to, photo_path ` +
		`FROM northwind.employees ` +
		`WHERE reports_to = ?`
	// run
	logf(sqlstr, reportsTo)
	rows, err := db.QueryContext(ctx, sqlstr, reportsTo)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Employee
	for rows.Next() {
		e := Employee{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&e.EmployeeID, &e.LastName, &e.FirstName, &e.Title, &e.TitleOfCourtesy, &e.BirthDate, &e.HireDate, &e.Address, &e.City, &e.Region, &e.PostalCode, &e.Country, &e.HomePhone, &e.Extension, &e.Photo, &e.Notes, &e.ReportsTo, &e.PhotoPath); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// Employee returns the Employee associated with the [Employee]'s (ReportsTo).
//
// Generated from foreign key 'employees_ibfk_1'.
func (e *Employee) Employee(ctx context.Context, db DB) (*Employee, error) {
	return EmployeeByEmployeeID(ctx, db, int16(e.ReportsTo.Int64))
}
