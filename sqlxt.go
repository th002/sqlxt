package sqlxt

import (
	"bytes"

	"github.com/jmoiron/sqlx"
)

//SqlxT sqlxt struct
type SqlxT struct {
	db     *sqlx.DB
	parser *Parse
}

//New new sqlxt
func New(db *sqlx.DB, path string) (*SqlxT, error) {
	parse, err := NewParse(path)
	if err != nil {
		return nil, err
	}
	return &SqlxT{
		db,
		parse,
	}, nil
}

//Exec exec sql
func (st *SqlxT) Exec(key string, v interface{}) error {
	tmp, err := st.parser.Get(key)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	err = tmp.Execute(buf, v)
	if err != nil {
		return err
	}

	st.db.Exec(buf.String())
	return nil
}

//Get get row
func (st *SqlxT) Get(key string, v interface{}) error {
	tmp, err := st.parser.Get(key)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	err = tmp.Execute(buf, v)
	if err != nil {
		return err
	}
	return st.db.Get(v, buf.String())
}

//Select get rows
func (st *SqlxT) Select(key string, v interface{}) error {
	tmp, err := st.parser.Get(key)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	err = tmp.Execute(buf, v)
	if err != nil {
		return err
	}
	return st.db.Select(v, buf.String())
}
