package query

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

type Params map[string]interface{}

type Query struct {
	query  string
	params Params
}

func New(base string) *Query {
	return &Query{
		query:  base,
		params: nil,
	}
}

func (q *Query) Where(cond string) *Query {
	q.query += " WHERE " + cond
	return q
}

func (q *Query) And(cond string) *Query {
	q.query += " AND " + cond
	return q
}

func (q *Query) Or(cond string) *Query {
	q.query += " OR " + cond
	return q
}

func (q *Query) OrderBy(field string, fields ...string) *Query {
	q.query += " ORDER BY " + strings.Join(append([]string{field}, fields...), ", ")
	return q
}

func (q *Query) Wrap(wrapper string) *Query {
	panic("unimplemented")
}

var ErrParamNotFound = errors.New("param not found")

func (e *executor) Param(name string) (string, error) {
	idx, found := e.m[name]
	if !found {
		if v, found := e.params[name]; !found {
			return "", fmt.Errorf("%w: %s", ErrParamNotFound, name)
		} else {
			e.args = append(e.args, v)
			idx, e.m[name] = len(e.args), len(e.args)
		}
	}
	return fmt.Sprintf("$%d", idx), nil
}

type executor struct {
	params Params
	m      map[string]int
	args   []interface{}
}

func (e *executor) Args() []interface{} {
	return e.args
}

func (q *Query) Build(params Params) (string, []interface{}, error) {
	e := executor{
		m:      make(map[string]int),
		params: params,
	}
	t, err := template.New("query-template").Parse(q.query)
	if err != nil {
		return "", nil, err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, &e); err != nil {
		return "", nil, err
	}
	return buf.String(), e.Args(), nil
}
