package oracle

// Code generated by dbtpl. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// A2In2Out calls the stored procedure 'a_bit_of_everything.a_2_in_2_out(number, number) (number, number)' on db.
func A2In2Out(ctx context.Context, db DB, paramOne, paramTwo int64) (int64, int64, error) {
	// call a_bit_of_everything.a_2_in_2_out
	const sqlstr = `BEGIN a_2_in_2_out(:param_one, :param_two, :return_two, :return_one); END;`
	// run
	var returnTwo int64
	var returnOne int64
	logf(sqlstr, paramOne, paramTwo)
	if _, err := db.ExecContext(ctx, sqlstr, sql.Named("param_one", paramOne), sql.Named("param_two", paramTwo), sql.Out{Dest: &returnTwo}, sql.Out{Dest: &returnOne}); err != nil {
		return 0, 0, logerror(err)
	}
	return returnTwo, returnOne, nil
}
