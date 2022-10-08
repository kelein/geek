package dialect

import (
	"reflect"
	"strings"

	"geek/glog"
)

// var dialects = map[string]Dialect{}
var dialects = make(map[string]Dialect)

// Dialect abstract for multipy DB
type Dialect interface {
	DataTypeOf(v reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect register multipy db dialect
func RegisterDialect(name string, dial Dialect) {
	dialects[name] = dial
}

// GetDialect return a Dialect and exist status
func GetDialect(name string) (Dialect, bool) {
	glog.Errorf(strings.Repeat("-", 20))
	glog.Infof("dialects map cap: %v", len(dialects))
	dial, ok := dialects[name]
	glog.Errorf(strings.Repeat("+", 20))
	glog.Infof("dialect %q exist status: %v, value: %v", name, ok, dial)
	return dial, ok
}

// Dialects debug func
func Dialects() map[string]Dialect {
	return dialects
}
