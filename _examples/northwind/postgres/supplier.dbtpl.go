package postgres

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Supplier represents a row from 'public.suppliers'.
type Supplier struct {
	SupplierID   int            `json:"supplier_id"`   // supplier_id
	CompanyName  string         `json:"company_name"`  // company_name
	ContactName  sql.NullString `json:"contact_name"`  // contact_name
	ContactTitle sql.NullString `json:"contact_title"` // contact_title
	Address      sql.NullString `json:"address"`       // address
	City         sql.NullString `json:"city"`          // city
	Region       sql.NullString `json:"region"`        // region
	PostalCode   sql.NullString `json:"postal_code"`   // postal_code
	Country      sql.NullString `json:"country"`       // country
	Phone        sql.NullString `json:"phone"`         // phone
	Fax          sql.NullString `json:"fax"`           // fax
	Homepage     sql.NullString `json:"homepage"`      // homepage
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Supplier] exists in the database.
func (s *Supplier) Exists() bool {
	return s._exists
}

// Deleted returns true when the [Supplier] has been marked for deletion
// from the database.
func (s *Supplier) Deleted() bool {
	return s._deleted
}

// Insert inserts the [Supplier] to the database.
func (s *Supplier) Insert(ctx context.Context, db DB) error {
	switch {
	case s._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case s._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO public.suppliers (` +
		`supplier_id, company_name, contact_name, contact_title, address, city, region, postal_code, country, phone, fax, homepage` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12` +
		`)`
	// run
	logf(sqlstr, s.SupplierID, s.CompanyName, s.ContactName, s.ContactTitle, s.Address, s.City, s.Region, s.PostalCode, s.Country, s.Phone, s.Fax, s.Homepage)
	if _, err := db.ExecContext(ctx, sqlstr, s.SupplierID, s.CompanyName, s.ContactName, s.ContactTitle, s.Address, s.City, s.Region, s.PostalCode, s.Country, s.Phone, s.Fax, s.Homepage); err != nil {
		return logerror(err)
	}
	// set exists
	s._exists = true
	return nil
}

// Update updates a [Supplier] in the database.
func (s *Supplier) Update(ctx context.Context, db DB) error {
	switch {
	case !s._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case s._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with composite primary key
	const sqlstr = `UPDATE public.suppliers SET ` +
		`company_name = $1, contact_name = $2, contact_title = $3, address = $4, city = $5, region = $6, postal_code = $7, country = $8, phone = $9, fax = $10, homepage = $11 ` +
		`WHERE supplier_id = $12`
	// run
	logf(sqlstr, s.CompanyName, s.ContactName, s.ContactTitle, s.Address, s.City, s.Region, s.PostalCode, s.Country, s.Phone, s.Fax, s.Homepage, s.SupplierID)
	if _, err := db.ExecContext(ctx, sqlstr, s.CompanyName, s.ContactName, s.ContactTitle, s.Address, s.City, s.Region, s.PostalCode, s.Country, s.Phone, s.Fax, s.Homepage, s.SupplierID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Supplier] to the database.
func (s *Supplier) Save(ctx context.Context, db DB) error {
	if s.Exists() {
		return s.Update(ctx, db)
	}
	return s.Insert(ctx, db)
}

// Upsert performs an upsert for [Supplier].
func (s *Supplier) Upsert(ctx context.Context, db DB) error {
	switch {
	case s._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO public.suppliers (` +
		`supplier_id, company_name, contact_name, contact_title, address, city, region, postal_code, country, phone, fax, homepage` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12` +
		`)` +
		` ON CONFLICT (supplier_id) DO ` +
		`UPDATE SET ` +
		`company_name = EXCLUDED.company_name, contact_name = EXCLUDED.contact_name, contact_title = EXCLUDED.contact_title, address = EXCLUDED.address, city = EXCLUDED.city, region = EXCLUDED.region, postal_code = EXCLUDED.postal_code, country = EXCLUDED.country, phone = EXCLUDED.phone, fax = EXCLUDED.fax, homepage = EXCLUDED.homepage `
	// run
	logf(sqlstr, s.SupplierID, s.CompanyName, s.ContactName, s.ContactTitle, s.Address, s.City, s.Region, s.PostalCode, s.Country, s.Phone, s.Fax, s.Homepage)
	if _, err := db.ExecContext(ctx, sqlstr, s.SupplierID, s.CompanyName, s.ContactName, s.ContactTitle, s.Address, s.City, s.Region, s.PostalCode, s.Country, s.Phone, s.Fax, s.Homepage); err != nil {
		return logerror(err)
	}
	// set exists
	s._exists = true
	return nil
}

// Delete deletes the [Supplier] from the database.
func (s *Supplier) Delete(ctx context.Context, db DB) error {
	switch {
	case !s._exists: // doesn't exist
		return nil
	case s._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM public.suppliers ` +
		`WHERE supplier_id = $1`
	// run
	logf(sqlstr, s.SupplierID)
	if _, err := db.ExecContext(ctx, sqlstr, s.SupplierID); err != nil {
		return logerror(err)
	}
	// set deleted
	s._deleted = true
	return nil
}

// SupplierBySupplierID retrieves a row from 'public.suppliers' as a [Supplier].
//
// Generated from index 'suppliers_pkey'.
func SupplierBySupplierID(ctx context.Context, db DB, supplierID int) (*Supplier, error) {
	// query
	const sqlstr = `SELECT ` +
		`supplier_id, company_name, contact_name, contact_title, address, city, region, postal_code, country, phone, fax, homepage ` +
		`FROM public.suppliers ` +
		`WHERE supplier_id = $1`
	// run
	logf(sqlstr, supplierID)
	s := Supplier{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, supplierID).Scan(&s.SupplierID, &s.CompanyName, &s.ContactName, &s.ContactTitle, &s.Address, &s.City, &s.Region, &s.PostalCode, &s.Country, &s.Phone, &s.Fax, &s.Homepage); err != nil {
		return nil, logerror(err)
	}
	return &s, nil
}
