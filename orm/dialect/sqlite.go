package dialect

import (
	"fmt"
	"reflect"
	"time"
)

const dialName = "sqlite3"

type sqlite struct{}

var _ Dialect = (*sqlite)(nil)

// func init() { RegisterDialect(dialName, &sqlite{}) }

func init() { RegisterDialect(dialName, &sqlite{}) }

func (s *sqlite) DataTypeOf(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := v.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", v.Type().Name(), v.Kind()))
}

func (s *sqlite) TableExistSQL(tableName string) (string, []interface{}) {
	sql := "SELECT name FROM sqlite_master WHERE type='table' AND name=?"
	args := []interface{}{tableName}
	return sql, args
}
