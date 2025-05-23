package sqlserver

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
)

// AFunc2In calls the stored function 'a_bit_of_everything.a_func_2_in(int, int) int' on db.
func AFunc2In(ctx context.Context, db DB, paramOne, paramTwo int) (int, error) {
	// call a_bit_of_everything.a_func_2_in
	const sqlstr = `SELECT a_bit_of_everything.a_func_2_in(@p1, @p2) AS OUT`
	// run
	var r0 int
	logf(sqlstr, paramOne, paramTwo)
	if err := db.QueryRowContext(ctx, sqlstr, paramOne, paramTwo).Scan(&r0); err != nil {
		return 0, logerror(err)
	}
	return r0, nil
}
