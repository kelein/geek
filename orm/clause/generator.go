package clause

import (
	"fmt"
	"strings"
)

var generators = make(map[Kind]generator, 6)

type generator func(values ...interface{}) (string, []interface{})

func init() {
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderby
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// INSERT INTO $table ($fields)
func _insert(values ...interface{}) (string, []interface{}) {
	table := values[0]
	fields := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("INSERT INTO %s (%v)", table, fields), []interface{}{}
}

// VALUES ($v1), ($v2), ...
func _values(values ...interface{}) (string, []interface{}) {
	bindVals := ""
	sql := strings.Builder{}
	vars := []interface{}{}

	sql.WriteString("VALUES")
	for i, val := range values {
		v := val.([]interface{})
		if bindVals == "" {
			bindVals = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindVals))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}

	return sql.String(), vars
}

// SELECT $fields FROM $table
func _select(values ...interface{}) (string, []interface{}) {
	table := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, table), []interface{}{}
}

// LIMIT $num
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// WHERE $desc
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

// ORDER BY $field $desc
func _orderby(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// UPDATE $table SET $field=$value
func _update(values ...interface{}) (string, []interface{}) {
	table := values[0]
	valMap := values[1].(map[string]interface{})
	keys := make([]string, len(valMap))
	vars := make([]interface{}, len(valMap))
	for k, v := range valMap {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	sql := fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(keys, ", "))
	return sql, vars
}

// DELETE FROM $table
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// SELECT count(*) FROM $table
func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}

func genBindVars(num int) string {
	s := strings.Split(strings.Repeat("?", num), "")
	return strings.Join(s, ", ")
}

func genBindVars2(num int) string {
	vars := make([]string, num)
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}
