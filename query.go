package query

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"go.uber.org/multierr"
)

type Params map[string]interface{}

type Query struct {
	query  string
	params Params
}

func New(base string) *Query {
	var q Query
	q.query = base
	return &q
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

func (e *executor) param(name string) string {
	idx, found := e.m[name]
	if !found {
		if v, found := e.params[name]; !found {
			e.err = multierr.Append(e.err, fmt.Errorf("%w: %s", ErrParamNotFound, name))
			return ""
		} else {
			e.args = append(e.args, v)
			idx, e.m[name] = len(e.args), len(e.args)
		}
	}
	return fmt.Sprintf("$%d", idx)
}

type executor struct {
	params Params
	m      map[string]int
	args   []interface{}

	err error
}

func (e *executor) Err() error {
	return e.err
}

func (e *executor) Args() []interface{} {
	return e.args
}

func (q *Query) Build(params Params) (string, []interface{}, error) {
	e := executor{
		m:      make(map[string]int),
		params: params,
	}
	t, err := template.New("").Funcs(template.FuncMap{"Param": e.param}).Parse(q.query)
	if err != nil {
		return "", nil, err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, params)
	if err != nil {
		return "", nil, err
	}
	return buf.String(), e.Args(), e.Err()
}
