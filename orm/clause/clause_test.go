package clause

import (
	"reflect"
	"testing"
)

func TestClause_Select(t *testing.T) {
	clause := new(Clause)
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "user", []string{"*"})
	clause.Set(WHERE, "name = ?", "Tom")
	clause.Set(ORDERBY, "name DESC")

	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Logf("clause.Build() sql = %q, vars = %v", sql, vars)
	wantSQL := "SELECT * FROM user WHERE name = ? ORDER BY name DESC LIMIT ?"
	wantVars := []interface{}{"Tom", 3}
	if sql != wantSQL {
		t.Errorf("clause.Build() sql failed")
	}
	if !reflect.DeepEqual(vars, wantVars) {
		t.Errorf("clause.Build() vars failed")
	}
}

func TestClause_Insert(t *testing.T) {
	// clause := new(Clause)
	// clause.Set(INSERT, "user", []string)
}

func TestClause_Build(t *testing.T) {
	t.Run("SELECT", func(t *testing.T) {
		TestClause_Select(t)
	})
}
