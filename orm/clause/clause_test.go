package clause

import (
	"reflect"
	"testing"
)

func TestClause_Set(t *testing.T) {

}

func TestClause_Build(t *testing.T) {
	type fields struct {
		sql  map[Kind]string
		vars map[Kind]interface{}
	}
	type args struct {
		orders []Kind
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Clause{
				sql:  tt.fields.sql,
				vars: tt.fields.vars,
			}
			got, got1 := c.Build(tt.args.orders...)
			if got != tt.want {
				t.Errorf("Clause.Build() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Clause.Build() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
