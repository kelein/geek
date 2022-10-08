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

	glog.Info("database connect success!")
	return &Engine{db: db, dialect: dialect}, nil
}

// Close closes the database connection
func (e *Engine) Close() { e.db.Close() }

// NewSession create a session of this ORM Engine
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
