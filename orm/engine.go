package orm

import (
	"database/sql"

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
