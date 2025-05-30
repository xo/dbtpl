package mysql

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Order represents a row from 'northwind.orders'.
type Order struct {
	OrderID        int16           `json:"order_id"`         // order_id
	CustomerID     sql.NullString  `json:"customer_id"`      // customer_id
	EmployeeID     sql.NullInt64   `json:"employee_id"`      // employee_id
	OrderDate      sql.NullTime    `json:"order_date"`       // order_date
	RequiredDate   sql.NullTime    `json:"required_date"`    // required_date
	ShippedDate    sql.NullTime    `json:"shipped_date"`     // shipped_date
	ShipVia        sql.NullInt64   `json:"ship_via"`         // ship_via
	Freight        sql.NullFloat64 `json:"freight"`          // freight
	ShipName       sql.NullString  `json:"ship_name"`        // ship_name
	ShipAddress    sql.NullString  `json:"ship_address"`     // ship_address
	ShipCity       sql.NullString  `json:"ship_city"`        // ship_city
	ShipRegion     sql.NullString  `json:"ship_region"`      // ship_region
	ShipPostalCode sql.NullString  `json:"ship_postal_code"` // ship_postal_code
	ShipCountry    sql.NullString  `json:"ship_country"`     // ship_country
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Order] exists in the database.
func (o *Order) Exists() bool {
	return o._exists
}

// Deleted returns true when the [Order] has been marked for deletion
// from the database.
func (o *Order) Deleted() bool {
	return o._deleted
}

// Insert inserts the [Order] to the database.
func (o *Order) Insert(ctx context.Context, db DB) error {
	switch {
	case o._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case o._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO northwind.orders (` +
		`order_id, customer_id, employee_id, order_date, required_date, shipped_date, ship_via, freight, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, o.OrderID, o.CustomerID, o.EmployeeID, o.OrderDate, o.RequiredDate, o.ShippedDate, o.ShipVia, o.Freight, o.ShipName, o.ShipAddress, o.ShipCity, o.ShipRegion, o.ShipPostalCode, o.ShipCountry)
	if _, err := db.ExecContext(ctx, sqlstr, o.OrderID, o.CustomerID, o.EmployeeID, o.OrderDate, o.RequiredDate, o.ShippedDate, o.ShipVia, o.Freight, o.ShipName, o.ShipAddress, o.ShipCity, o.ShipRegion, o.ShipPostalCode, o.ShipCountry); err != nil {
		return logerror(err)
	}
	// set exists
	o._exists = true
	return nil
}

// Update updates a [Order] in the database.
func (o *Order) Update(ctx context.Context, db DB) error {
	switch {
	case !o._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case o._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE northwind.orders SET ` +
		`customer_id = ?, employee_id = ?, order_date = ?, required_date = ?, shipped_date = ?, ship_via = ?, freight = ?, ship_name = ?, ship_address = ?, ship_city = ?, ship_region = ?, ship_postal_code = ?, ship_country = ? ` +
		`WHERE order_id = ?`
	// run
	logf(sqlstr, o.CustomerID, o.EmployeeID, o.OrderDate, o.RequiredDate, o.ShippedDate, o.ShipVia, o.Freight, o.ShipName, o.ShipAddress, o.ShipCity, o.ShipRegion, o.ShipPostalCode, o.ShipCountry, o.OrderID)
	if _, err := db.ExecContext(ctx, sqlstr, o.CustomerID, o.EmployeeID, o.OrderDate, o.RequiredDate, o.ShippedDate, o.ShipVia, o.Freight, o.ShipName, o.ShipAddress, o.ShipCity, o.ShipRegion, o.ShipPostalCode, o.ShipCountry, o.OrderID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Order] to the database.
func (o *Order) Save(ctx context.Context, db DB) error {
	if o.Exists() {
		return o.Update(ctx, db)
	}
	return o.Insert(ctx, db)
}

// Upsert performs an upsert for [Order].
func (o *Order) Upsert(ctx context.Context, db DB) error {
	switch {
	case o._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO northwind.orders (` +
		`order_id, customer_id, employee_id, order_date, required_date, shipped_date, ship_via, freight, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`order_id = VALUES(order_id), customer_id = VALUES(customer_id), employee_id = VALUES(employee_id), order_date = VALUES(order_date), required_date = VALUES(required_date), shipped_date = VALUES(shipped_date), ship_via = VALUES(ship_via), freight = VALUES(freight), ship_name = VALUES(ship_name), ship_address = VALUES(ship_address), ship_city = VALUES(ship_city), ship_region = VALUES(ship_region), ship_postal_code = VALUES(ship_postal_code), ship_country = VALUES(ship_country)`
	// run
	logf(sqlstr, o.OrderID, o.CustomerID, o.EmployeeID, o.OrderDate, o.RequiredDate, o.ShippedDate, o.ShipVia, o.Freight, o.ShipName, o.ShipAddress, o.ShipCity, o.ShipRegion, o.ShipPostalCode, o.ShipCountry)
	if _, err := db.ExecContext(ctx, sqlstr, o.OrderID, o.CustomerID, o.EmployeeID, o.OrderDate, o.RequiredDate, o.ShippedDate, o.ShipVia, o.Freight, o.ShipName, o.ShipAddress, o.ShipCity, o.ShipRegion, o.ShipPostalCode, o.ShipCountry); err != nil {
		return logerror(err)
	}
	// set exists
	o._exists = true
	return nil
}

// Delete deletes the [Order] from the database.
func (o *Order) Delete(ctx context.Context, db DB) error {
	switch {
	case !o._exists: // doesn't exist
		return nil
	case o._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM northwind.orders ` +
		`WHERE order_id = ?`
	// run
	logf(sqlstr, o.OrderID)
	if _, err := db.ExecContext(ctx, sqlstr, o.OrderID); err != nil {
		return logerror(err)
	}
	// set deleted
	o._deleted = true
	return nil
}

// OrdersByCustomerID retrieves a row from 'northwind.orders' as a [Order].
//
// Generated from index 'customer_id'.
func OrdersByCustomerID(ctx context.Context, db DB, customerID sql.NullString) ([]*Order, error) {
	// query
	const sqlstr = `SELECT ` +
		`order_id, customer_id, employee_id, order_date, required_date, shipped_date, ship_via, freight, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country ` +
		`FROM northwind.orders ` +
		`WHERE customer_id = ?`
	// run
	logf(sqlstr, customerID)
	rows, err := db.QueryContext(ctx, sqlstr, customerID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Order
	for rows.Next() {
		o := Order{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&o.OrderID, &o.CustomerID, &o.EmployeeID, &o.OrderDate, &o.RequiredDate, &o.ShippedDate, &o.ShipVia, &o.Freight, &o.ShipName, &o.ShipAddress, &o.ShipCity, &o.ShipRegion, &o.ShipPostalCode, &o.ShipCountry); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// OrdersByEmployeeID retrieves a row from 'northwind.orders' as a [Order].
//
// Generated from index 'employee_id'.
func OrdersByEmployeeID(ctx context.Context, db DB, employeeID sql.NullInt64) ([]*Order, error) {
	// query
	const sqlstr = `SELECT ` +
		`order_id, customer_id, employee_id, order_date, required_date, shipped_date, ship_via, freight, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country ` +
		`FROM northwind.orders ` +
		`WHERE employee_id = ?`
	// run
	logf(sqlstr, employeeID)
	rows, err := db.QueryContext(ctx, sqlstr, employeeID)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Order
	for rows.Next() {
		o := Order{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&o.OrderID, &o.CustomerID, &o.EmployeeID, &o.OrderDate, &o.RequiredDate, &o.ShippedDate, &o.ShipVia, &o.Freight, &o.ShipName, &o.ShipAddress, &o.ShipCity, &o.ShipRegion, &o.ShipPostalCode, &o.ShipCountry); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// OrderByOrderID retrieves a row from 'northwind.orders' as a [Order].
//
// Generated from index 'orders_order_id_pkey'.
func OrderByOrderID(ctx context.Context, db DB, orderID int16) (*Order, error) {
	// query
	const sqlstr = `SELECT ` +
		`order_id, customer_id, employee_id, order_date, required_date, shipped_date, ship_via, freight, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country ` +
		`FROM northwind.orders ` +
		`WHERE order_id = ?`
	// run
	logf(sqlstr, orderID)
	o := Order{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, orderID).Scan(&o.OrderID, &o.CustomerID, &o.EmployeeID, &o.OrderDate, &o.RequiredDate, &o.ShippedDate, &o.ShipVia, &o.Freight, &o.ShipName, &o.ShipAddress, &o.ShipCity, &o.ShipRegion, &o.ShipPostalCode, &o.ShipCountry); err != nil {
		return nil, logerror(err)
	}
	return &o, nil
}

// OrdersByShipVia retrieves a row from 'northwind.orders' as a [Order].
//
// Generated from index 'ship_via'.
func OrdersByShipVia(ctx context.Context, db DB, shipVia sql.NullInt64) ([]*Order, error) {
	// query
	const sqlstr = `SELECT ` +
		`order_id, customer_id, employee_id, order_date, required_date, shipped_date, ship_via, freight, ship_name, ship_address, ship_city, ship_region, ship_postal_code, ship_country ` +
		`FROM northwind.orders ` +
		`WHERE ship_via = ?`
	// run
	logf(sqlstr, shipVia)
	rows, err := db.QueryContext(ctx, sqlstr, shipVia)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()
	// process
	var res []*Order
	for rows.Next() {
		o := Order{
			_exists: true,
		}
		// scan
		if err := rows.Scan(&o.OrderID, &o.CustomerID, &o.EmployeeID, &o.OrderDate, &o.RequiredDate, &o.ShippedDate, &o.ShipVia, &o.Freight, &o.ShipName, &o.ShipAddress, &o.ShipCity, &o.ShipRegion, &o.ShipPostalCode, &o.ShipCountry); err != nil {
			return nil, logerror(err)
		}
		res = append(res, &o)
	}
	if err := rows.Err(); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// Customer returns the Customer associated with the [Order]'s (CustomerID).
//
// Generated from foreign key 'orders_ibfk_1'.
func (o *Order) Customer(ctx context.Context, db DB) (*Customer, error) {
	return CustomerByCustomerID(ctx, db, o.CustomerID.String)
}

// Employee returns the Employee associated with the [Order]'s (EmployeeID).
//
// Generated from foreign key 'orders_ibfk_2'.
func (o *Order) Employee(ctx context.Context, db DB) (*Employee, error) {
	return EmployeeByEmployeeID(ctx, db, int16(o.EmployeeID.Int64))
}

// Shipper returns the Shipper associated with the [Order]'s (ShipVia).
//
// Generated from foreign key 'orders_ibfk_3'.
func (o *Order) Shipper(ctx context.Context, db DB) (*Shipper, error) {
	return ShipperByShipperID(ctx, db, int16(o.ShipVia.Int64))
}
