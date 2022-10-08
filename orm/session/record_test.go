package session

import (
	"database/sql"
	"geek/orm/clause"
	"geek/orm/dialect"
	"geek/orm/schema"
	"strings"
	"testing"
)

func TestSession_Insert(t *testing.T) {
	type fields struct {
		db       *sql.DB
		dialect  dialect.Dialect
		refTable *schema.Schema
		sql      strings.Builder
		vals     []interface{}
		clause   clause.Clause
	}
	type args struct {
		records []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				db:       tt.fields.db,
				dialect:  tt.fields.dialect,
				refTable: tt.fields.refTable,
				sql:      tt.fields.sql,
				vals:     tt.fields.vals,
				clause:   tt.fields.clause,
			}
			got, err := s.Insert(tt.args.records...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Session.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Session.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Find(t *testing.T) {
	type fields struct {
		db       *sql.DB
		dialect  dialect.Dialect
		refTable *schema.Schema
		sql      strings.Builder
		vals     []interface{}
		clause   clause.Clause
	}
	type args struct {
		records interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				db:       tt.fields.db,
				dialect:  tt.fields.dialect,
				refTable: tt.fields.refTable,
				sql:      tt.fields.sql,
				vals:     tt.fields.vals,
				clause:   tt.fields.clause,
			}
			if err := s.Find(tt.args.records); (err != nil) != tt.wantErr {
				t.Errorf("Session.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
