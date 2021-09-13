package loader

import (
	"context"
	"regexp"
	"strings"

	"github.com/xo/xo/models"
	xo "github.com/xo/xo/types"
)

func init() {
	Register(&Loader{
		Driver:           "mysql",
		Mask:             "?",
		GoType:           MysqlGoType,
		Schema:           models.MysqlSchema,
		Enums:            models.MysqlEnums,
		EnumValues:       MysqlEnumValues,
		Procs:            models.MysqlProcs,
		ProcParams:       models.MysqlProcParams,
		Tables:           models.MysqlTables,
		TableColumns:     models.MysqlTableColumns,
		TableSequences:   models.MysqlTableSequences,
		TableForeignKeys: models.MysqlTableForeignKeys,
		TableIndexes:     models.MysqlTableIndexes,
		IndexColumns:     models.MysqlIndexColumns,
		ViewCreate:       models.MysqlViewCreate,
		ViewDrop:         models.MysqlViewDrop,
	})
}

// MysqlGoType parse a mysql type into a Go type based on the column
// definition.
func MysqlGoType(ctx context.Context, d xo.Datatype) (string, string, error) {
	var goType, zero string
	switch d.Type {
	case "bit":
		switch {
		case d.Prec == 1 && !d.Nullable:
			goType, zero = "bool", "false"
		case d.Prec == 1 && d.Nullable:
			goType, zero = "null.Bool", "sql.NullBool{}"
		case d.Prec <= 8 && !d.Nullable:
			goType, zero = "uint8", "0"
		case d.Prec <= 16 && !d.Nullable:
			goType, zero = "uint16", "0"
		case d.Prec <= 32 && !d.Nullable:
			goType, zero = "uint32", "0"
		case d.Nullable:
			goType, zero = "null.Int", "sql.NullInt64{}"
		default:
			goType, zero = "uint64", "0"
		}
	case "bool", "boolean":
		goType, zero = "bool", "false"
		if d.Nullable {
			goType, zero = "null.Bool", "sql.NullBool{}"
		}
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		goType, zero = "string", `""`
		if d.Nullable {
			goType, zero = "null.String", "sql.NullString{}"
		}
	case "tinyint":
		switch {
		case d.Prec == 1 && !d.Nullable: // force tinyint(1) as bool
			goType, zero = "bool", "false"
		case d.Prec == 1 && d.Nullable:
			goType, zero = "null.Bool", "sql.NullBool{}"
		case d.Nullable:
			goType, zero = "null.Int", "sql.NullInt64{}"
		default:
			goType, zero = "int8", "0"
		}
	case "smallint", "year":
		goType, zero = "int16", "0"
		if d.Nullable {
			goType, zero = "null.Int", "sql.NullInt64{}"
		}
	case "mediumint", "int", "integer":
		goType, zero = Int32(ctx), "0"
		if d.Nullable {
			goType, zero = "null.Int", "sql.NullInt64{}"
		}
	case "bigint":
		goType, zero = "int64", "0"
		if d.Nullable {
			goType, zero = "null.Int", "sql.NullInt64{}"
		}
	case "float":
		goType, zero = "float32", "0.0"
		if d.Nullable {
			goType, zero = "null.Float", "sql.NullFloat64{}"
		}
	case "decimal", "double":
		goType, zero = "float64", "0.0"
		if d.Nullable {
			goType, zero = "null.Float", "sql.NullFloat64{}"
		}
	case "binary", "blob", "longblob", "mediumblob", "tinyblob", "varbinary":
		goType, zero = "[]byte", "nil"
	case "timestamp", "datetime", "date":
		goType, zero = "time.Time", "time.Time{}"
		if d.Nullable {
			goType, zero = "null.Time", "sql.NullTime{}"
		}
	case "time":
		// time is not supported by the MySQL driver. Can parse the string to time.Time in the user code.
		goType, zero = "string", `""`
		if d.Nullable {
			goType, zero = "null.String", "sql.NullString{}"
		}
	default:
		goType, zero = SchemaGoType(ctx, d.Type, d.Nullable)
	}
	// force []byte for SET('a',...)
	if setRE.MatchString(d.Type) {
		goType, zero = "[]byte", "nil"
	}
	// if unsigned ...
	if intRE.MatchString(goType) && d.Unsigned {
		if goType == Int32(ctx) {
			goType, zero = Uint32(ctx), "0"
		} else {
			goType = "u" + goType
		}
	}
	return goType, zero, nil
}

// setRE is the regexp that matches MySQL SET() type definitions.
var setRE = regexp.MustCompile(`(?i)^set\([^)]*\)$`)

// MysqlEnumValues loads the enum values.
func MysqlEnumValues(ctx context.Context, db models.DB, schema string, enum string) ([]*models.EnumValue, error) {
	// load enum values
	res, err := models.MysqlEnumValues(ctx, db, schema, enum)
	if err != nil {
		return nil, err
	}
	// process enum values
	var values []*models.EnumValue
	for i, val := range strings.Split(res.EnumValues[1:len(res.EnumValues)-1], "','") {
		values = append(values, &models.EnumValue{
			EnumValue:  val,
			ConstValue: i + 1,
		})
	}
	return values, nil
}
