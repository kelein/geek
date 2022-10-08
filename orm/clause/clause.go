package clause

import "strings"

// Segment Kind List
const (
	INSERT Kind = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Kind stands for subset segment in SQL
type Kind uint8

// Clause for SQL query
type Clause struct {
	sql  map[Kind]string
	vars map[Kind][]interface{}
}

// Set setting up a SQL clause
func (c *Clause) Set(name Kind, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Kind]string)
		c.vars = make(map[Kind][]interface{})
	}
	c.sql[name], c.vars[name] = generators[name](vars...)
}

// Build generate a whole SQL statement
func (c *Clause) Build(orders ...Kind) (string, []interface{}) {
	sqls := []string{}
	vars := []interface{}{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.vars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
