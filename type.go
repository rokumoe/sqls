package sqls

import (
	"database/sql"
	"time"
)

type Bool bool

func (p *Bool) Scan(src interface{}) error {
	var t sql.NullBool
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = Bool(t.Bool)
	return nil
}

type I64 int64

func (p *I64) Scan(src interface{}) error {
	var t sql.NullInt64
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = I64(t.Int64)
	return nil
}

type I32 int32

func (p *I32) Scan(src interface{}) error {
	var t sql.NullInt32
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = I32(t.Int32)
	return nil
}

type U64 int64

func (p *U64) Scan(src interface{}) error {
	var t sql.NullInt64
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = U64(t.Int64)
	return nil
}

type U32 int32

func (p *U32) Scan(src interface{}) error {
	var t sql.NullInt32
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = U32(t.Int32)
	return nil
}

type F64 float64

func (p *F64) Scan(src interface{}) error {
	var t sql.NullFloat64
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = F64(t.Float64)
	return nil
}

type String string

func (p *String) Scan(src interface{}) error {
	var t sql.NullString
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = String(t.String)
	return nil
}

type Time time.Time

func (p *Time) Scan(src interface{}) error {
	var t sql.NullTime
	err := t.Scan(src)
	if err != nil {
		return err
	}
	*p = Time(t.Time)
	return nil
}
