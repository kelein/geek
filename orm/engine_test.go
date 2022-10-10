package orm

import (
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"geek/orm/session"
)

type User struct {
	Name string `GEORM:"PRIMARY KEY"`
	Age  int
}

func testOpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "geek.db")
	if err != nil {
		t.Fatalf("open sqlite3 failed: %v", err)
	}
	return engine
}

func TestTransaction_Rollback(t *testing.T) {
	engine := testOpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{}).DropTable()

	txFn := func(s *session.Session) (interface{}, error) {
		s.Model(&User{}).CreateTable()
		s.Insert(&User{"Virgo", 24})
		return nil, errors.New("mock error")
	}

	_, err := engine.Transaction(txFn)
	if err == nil || s.HasTable() {
		t.Fatal("transaction rollbac failed")
	}
}

func TestTransaction_Commit(t *testing.T) {
	engine := testOpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.Model(&User{}).DropTable()

	txFn := func(s *session.Session) (interface{}, error) {
		s.Model(&User{}).CreateTable()
		result, err := s.Insert(&User{"Virgo", 24})
		return result, err
	}

	_, err := engine.Transaction(txFn)
	u := &User{}
	s.Get(u)
	if err != nil || u.Name != "Virgo" {
		t.Fatal("transaction commits failed")
	}
}

func TestEngine_Transaction(t *testing.T) {
	tests := []struct {
		name string
		args func(t *testing.T)
	}{
		{"Rollback", TestTransaction_Rollback},
		{"Commit", TestTransaction_Commit},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args(t)
		})
	}
}
