package mysql

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
)

// A1In0Out calls the stored procedure 'a_bit_of_everything.a_1_in_0_out(int)' on db.
func A1In0Out(ctx context.Context, db DB, aParam int) error {
	// call a_bit_of_everything.a_1_in_0_out
	const sqlstr = `CALL a_bit_of_everything.a_1_in_0_out(?)`
	// run
	logf(sqlstr)
	if _, err := db.ExecContext(ctx, sqlstr, aParam); err != nil {
		return logerror(err)
	}
	return nil
}
