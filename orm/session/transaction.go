package session

import "geek/glog"

// Begin start a transaction
func (s *Session) Begin() error {
	glog.Info("begin transaction")
	tx, err := s.db.Begin()
	if err != nil {
		glog.Error(err)
		return err
	}
	s.tx = tx
	return nil
}

// Commit commits the transaction
func (s *Session) Commit() error {
	glog.Info("commit transaction")
	if err := s.tx.Commit(); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

// Rollback rollbacked the transaction
func (s *Session) Rollback() error {
	glog.Info("rollback transaction")
	if err := s.tx.Rollback(); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
