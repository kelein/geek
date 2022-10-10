package session

import (
	"reflect"

	"geek/glog"
)

// Hooks Kind Constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) invoke(method string, value interface{}) {
	fn := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fn = reflect.ValueOf(value).MethodByName(method)
	}

	param := []reflect.Value{reflect.ValueOf(s)}
	if fn.IsValid() {
		if v := fn.Call(param); len(v) > 0 {
			err, ok := v[0].Interface().(error)
			if ok {
				glog.Error(err)
			}
		}
	}
	return
}

// IAfterQuery for abstract interface
type IAfterQuery interface {
	AfterQuery(s *Session) error
}

// IBeforeInsert for abstract interface
type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

// CallMethod invokes the registered hooks
func (s *Session) CallMethod(method string, value interface{}) {
	param := reflect.ValueOf(value)

	switch method {
	case AfterQuery:
		if fn, ok := param.Interface().(IAfterQuery); ok {
			fn.AfterQuery(s)
		}
	case BeforeInsert:
		if fn, ok := param.Interface().(IBeforeInsert); ok {
			fn.BeforeInsert(s)
		}
	default:
		panic("unsupported hook: " + method)
	}
	return
}
