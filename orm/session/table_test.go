package session

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"geek/orm/dialect"
)

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

func TestSession_CreateTable(t *testing.T) {
	driver := "sqlite3"

	t.Logf("support dialect list: %v", dialect.Dialects())
	sqliteDial, _ := dialect.GetDialect(driver)
	t.Logf("sqliteDial: %v", sqliteDial)

	db, err := sql.Open(driver, "geek.db")
	if err != nil {
		t.Errorf("open db err: %v", err)
	}
	t.Logf("db info: %v", db)

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		t.Logf("current dialect list: %v", dialect.Dialects())
		t.Errorf("dialect %q not found", driver)
	}
	t.Logf("current dialect list: %v", dialect.Dialects())
	t.Logf("dialect info: %v", dial)

	s := New(db, dial).Model(&Group{})
	s.DropTable()
	s.DropTable()
	s.CreateTable()
	if !s.HasTable() {
		t.Errorf("create table %q failed", s.RefTable().Name)
	}
}
