package schema

import (
	"go/ast"
	"reflect"

	"geek/orm/dialect"
)

var schemaTag = "GEORM"

// Field stands for database table column
type Field struct{ Name, Type, Tag string }

// Schema stands for database table
type Schema struct {
	Name       string
	Model      interface{}
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// GetField return a Field by name
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// Parse convert any object into ORM Schema
func Parse(dest interface{}, dialect dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Name:     modelType.Name(),
		Model:    dest,
		fieldMap: make(map[string]*Field),
	}

	num := modelType.NumField()
	for i := 0; i < num; i++ {
		m := modelType.Field(i)
		if !m.Anonymous && ast.IsExported(m.Name) {
			field := &Field{
				Name: m.Name,
				Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(m.Type))),
			}

			if v, ok := m.Tag.Lookup(schemaTag); ok {
				field.Tag = v
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, m.Name)
			schema.fieldMap[m.Name] = field
		}
	}

	return schema
}
