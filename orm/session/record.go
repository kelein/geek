package session

import (
	"errors"
	"reflect"

	"geek/orm/clause"
)

// Insert save record into database.
// Example:
// s := orm.NewEngine("sqlite3", "geek.db").NewSession()
// u1 := &User{Name: "Plato", Age: 23}
// u2 := &User{Name: "Kathy", Age: 28}
// s.Insert(u1, u2)
func (s *Session) Insert(records ...interface{}) (int64, error) {
	recordVals := make([]interface{}, len(records))
	for _, v := range records {
		table := s.Model(v).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordVals = append(recordVals, table.RecordValues(v))
	}

	s.clause.Set(clause.VALUES, recordVals...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find search records from database.
// Example:
// s := orm.NewEngine("sqlite3", "geek.db").NewSession()
// users := []User{}
// s.Find(&users)
func (s *Session) Find(records interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(records))
	// * 获取切片中的元素类型
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)

	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		values := make([]interface{}, len(table.FieldNames))
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}

	return rows.Close()
}

// Update record by map or key-value pair params.
// Params: map[string]interface{} / (k1, v1, k2, v2)
func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		count := len(kv)
		for i := 0; i < count; i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	val := int64(0)
	if err := row.Scan(&val); err != nil {
		return -1, err
	}
	return val, nil
}

// Limit set limit clause to sql statement
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where set where clause to sql statement
func (s *Session) Where(desc string, args ...interface{}) *Session {
	vars := []interface{}{}
	vars = append(append(vars, desc), args...)
	s.clause.Set(clause.WHERE, vars)
	return s
}

// OrderBy set order by statement to sql
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// Get return first row which matches clause.
// Example:
// u := &User{}
// s.OrderBy("name", "DESC").Get(u)
func (s *Session) Get(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type().Elem()))
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}

	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}

	dest.Set(destSlice.Index(0))
	return nil
}
