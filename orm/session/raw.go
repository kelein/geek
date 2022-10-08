package session

import (
	"database/sql"
	"strings"

	"geek/glog"
	"geek/orm/clause"
	"geek/orm/dialect"
	"geek/orm/schema"
)

// Session for Database
type Session struct {
	db       *sql.DB
	dialect  dialect.Dialect
	refTable *schema.Schema
	sql      strings.Builder
	vals     []interface{}
	clause   clause.Clause
}

// New create db session instance
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

// Clear reset sql builder
func (s *Session) Clear() {
	s.sql.Reset()
	s.vals = nil
	s.clause = clause.Clause{}
}

// DB get db instance for session
func (s *Session) DB() *sql.DB { return s.db }

// Raw is a original raw SQL builder
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.vals = append(s.vals, values...)
	return s
}

// Exec execute raw sql with values
func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	glog.Info(s.sql.String(), s.vals)
	result, err := s.DB().Exec(s.sql.String(), s.vals...)
	if err != nil {
		glog.Errorf("Exec err: %v", err)
		return nil, err
	}
	return result, nil
}

// QueryRow get row from database
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	glog.Info(s.sql.String(), s.vals)
	return s.DB().QueryRow(s.sql.String(), s.vals)
}

// QueryRows get multipy rows from database
func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	glog.Info(s.sql.String(), s.vals)
	return s.DB().Query(s.sql.String(), s.vals)
}
