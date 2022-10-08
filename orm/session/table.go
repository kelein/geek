package session

import (
	"fmt"
	"reflect"
	"strings"

	"geek/glog"
	"geek/orm/schema"
)

// Model update session refTable via schema
func (s *Session) Model(value interface{}) *Session {
	equal := reflect.TypeOf(value) == reflect.TypeOf(s.refTable.Model)
	if s.refTable == nil || !equal {
		s.refTable = schema.Parse(value, s.dialect)
	}

	return s
}

// RefTable get database session refTable
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		glog.Error("session modal is not set")
	}
	return s.refTable
}

// CreateTable build create table statement and execute
func (s *Session) CreateTable() error {
	table := s.RefTable()
	columns := make([]string, len(table.FieldNames))
	for _, fd := range table.Fields {
		col := fmt.Sprintf("%s %s %s", fd.Name, fd.Type, fd.Tag)
		columns = append(columns, col)
	}

	desc := strings.Join(columns, ", ")
	createSQL := fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)
	_, err := s.Raw(createSQL).Exec()
	return err
}

// DropTable build drop table statement and execute
func (s *Session) DropTable() error {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.RefTable().Name)
	_, err := s.Raw(dropSQL).Exec()
	return err
}

// HasTable check whether table exist
func (s *Session) HasTable() bool {
	sql, vals := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, vals...).QueryRow()

	name := ""
	if err := row.Scan(&name); err != nil {
		glog.Errorf("row scan err: %v", err)
		return false
	}
	return s.RefTable().Name == name
}
