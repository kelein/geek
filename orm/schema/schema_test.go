package schema

import (
	"encoding/json"
	"testing"
	"time"

	"geek/orm/dialect"
)

var sqlite, _ = dialect.GetDialect("sqlite3")

type User struct {
	Name string `GEORM:"PRIMARY KEY"`
	Age  int    `GEORM:"age"`
}

type Group struct {
	ID      int `GEORM:"PRIMARY KEY"`
	Name    string
	Users   []*User
	Updated time.Time
}

func TestParse(t *testing.T) {
	type args struct {
		dest    interface{}
		dialect dialect.Dialect
	}
	tests := []struct {
		name string
		args args
	}{
		{"A", args{&User{}, sqlite}},
		{"B", args{&Group{}, sqlite}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.args.dest, tt.args.dialect)
			t.Logf("dialect.Parse() got = %v", got)
			s, _ := json.Marshal(got)
			t.Logf("marshal schema: %v", string(s))
		})
	}
}
