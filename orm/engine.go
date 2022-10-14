package orm

import (
	"database/sql"
	"fmt"
	"strings"

	"geek/glog"
	"geek/orm/dialect"
	"geek/orm/session"
)

// Engine of ORM
type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// NewEngine create ORM Engine instance
func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		glog.Errorf("open db failed: %v", err)
		return nil, err
	}

	// ensure database connection alive
	if err := db.Ping(); err != nil {
		glog.Errorf("ping db failed: %v", err)
		return nil, err
	}

	// ensure database dialect exists
	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		glog.Errorf("dailect %q not found", driver)
	}
	glog.Errorf(glog.Red.Bold("dailect info: %v"), dialect)

	glog.Info("database connect success!")
	return &Engine{db: db, dialect: dialect}, nil
}

// Close closes the database connection
func (e *Engine) Close() { e.db.Close() }

// NewSession create a session of this ORM Engine
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

// Txfunc stands for Transaction callback function
type Txfunc func(*session.Session) (interface{}, error)

// Transaction create a transaction
func (e *Engine) Transaction(f Txfunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			s.Rollback()
			panic(r)
		} else if err != nil {
			s.Rollback()
		} else {
			// if commit failed update err and rollback
			err = s.Commit()
		}
	}()

	return f(s)
}

// Migrate migrates the given table
func (e *Engine) Migrate(value interface{}) error {
	txFn := func(s *session.Session) (interface{}, error) {
		if !s.Model(value).HasTable() {
			glog.Warn("table %q does not exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}

		table := s.RefTable()
		sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)
		rows, err := s.Raw(sql).QueryRows()
		if err != nil {
			glog.Error("query rows failed: %v", err)
			return nil, err
		}

		columns, _ := rows.Columns()
		addCols := diff(table.FieldNames, columns)
		delCols := diff(columns, table.FieldNames)
		glog.Infof("added columns: %v, deleted columns: %v", addCols, delCols)

		// * Add Table Columns
		for _, col := range addCols {
			fd := table.GetField(col)
			sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, fd.Name, fd.Type)
			if _, err := s.Raw(sql).Exec(); err != nil {
				return nil, err
			}
		}

		if len(delCols) == 0 {
			return nil, nil
		}

		// * Delete Table Coloumns
		tmp := "tmp_" + table.Name
		fileds := strings.Join(table.FieldNames, ",")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fileds, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return nil, err
	}
	_, err := e.Transaction(txFn)
	return err
}

// diff returns differences between a and b
func diff(a, b []string) []string {
	B := make(map[string]any, len(b))
	for _, v := range b {
		B[v] = true
	}

	items := []string{}
	for _, v := range a {
		if _, ok := B[v]; !ok {
			items = append(items, v)
		}
	}
	return items
}
