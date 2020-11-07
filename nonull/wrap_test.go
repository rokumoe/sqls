package nonull

import (
	"reflect"
	"testing"
)

func TestWrap(t *testing.T) {
	type User struct {
		Id   int
		Name string `db:"name"`
		Age  int    `nonull:"-"`
	}
	v := &User{}
	t.Log(reflect.TypeOf(v))
	x := Wrap(v)
	t.Log(reflect.TypeOf(x))
}
