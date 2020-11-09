package nonull

import (
	"database/sql"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/rokumoe/sqls"
)

var rtcache struct {
	sync.RWMutex
	m           map[reflect.Type]reflect.Type
	scannerType reflect.Type
}

func deref(i interface{}) reflect.Type {
	return reflect.TypeOf(i).Elem()
}

func init() {
	rtcache.m = map[reflect.Type]reflect.Type{
		deref((*int64)(nil)):     deref((*sqls.I64)(nil)),
		deref((*int32)(nil)):     deref((*sqls.I32)(nil)),
		deref((*uint64)(nil)):    deref((*sqls.U64)(nil)),
		deref((*uint32)(nil)):    deref((*sqls.U32)(nil)),
		deref((*bool)(nil)):      deref((*sqls.Bool)(nil)),
		deref((*string)(nil)):    deref((*sqls.String)(nil)),
		deref((*float64)(nil)):   deref((*sqls.F64)(nil)),
		deref((*time.Time)(nil)): deref((*sqls.Time)(nil)),
	}
	if unsafe.Sizeof(0) == 8 {
		rtcache.m[deref((*int)(nil))] = deref((*sqls.I64)(nil))
	} else {
		rtcache.m[deref((*int)(nil))] = deref((*sqls.I32)(nil))
	}
	rtcache.scannerType = deref((*sql.Scanner)(nil))
}

func mapRType(t reflect.Type) reflect.Type {
	if t.Implements(rtcache.scannerType) {
		return t
	}
	switch t.Kind() {
	case reflect.Int64:
		return deref((*sqls.I64)(nil))
	case reflect.Int32:
		return deref((*sqls.I32)(nil))
	case reflect.Int:
		if unsafe.Sizeof(0) == 8 {
			return deref((*sqls.I64)(nil))
		} else {
			return deref((*sqls.I32)(nil))
		}
	case reflect.Uint64:
		return deref((*sqls.U64)(nil))
	case reflect.Uint32:
		return deref((*sqls.U32)(nil))
	case reflect.Uint:
		if unsafe.Sizeof(uint(0)) == 8 {
			return deref((*sqls.U64)(nil))
		} else {
			return deref((*sqls.U32)(nil))
		}
	case reflect.String:
		return deref((*sqls.String)(nil))
	case reflect.Bool:
		return deref((*sqls.Bool)(nil))
	case reflect.Float64:
		return deref((*sqls.F64)(nil))
	case reflect.Ptr:
		elem := mapCachedType(t.Elem())
		if elem != t.Elem() {
			return reflect.PtrTo(elem)
		}
	case reflect.Struct:
		n := t.NumField()
		fields := make([]reflect.StructField, n)
		needMap := false
		for i := 0; i < n; i++ {
			sf := t.Field(i)
			attr := sf.Tag.Get("nonull")
			if attr != "-" {
				psftyp := reflect.PtrTo(sf.Type)
				pftyp := mapCachedType(psftyp)
				if pftyp != psftyp {
					needMap = true
					sf.Type = pftyp.Elem()
				}
			}
			fields[i] = sf
		}
		if needMap {
			return reflect.StructOf(fields)
		}
	case reflect.Slice:
		elem := mapCachedType(t.Elem())
		if elem != t.Elem() {
			return reflect.SliceOf(elem)
		}
	case reflect.Array:
		elem := mapCachedType(t.Elem())
		if elem != t.Elem() {
			return reflect.ArrayOf(t.Len(), elem)
		}
	}
	return t
}

func mapCachedType(t reflect.Type) reflect.Type {
	r := rtcache.m[t]
	if r == nil {
		r = mapRType(t)
		if r.Kind() == reflect.Struct {
			rtcache.m[t] = r
		}
	}
	return r
}

func Wrap(out interface{}) interface{} {
	if _, ok := out.(sql.Scanner); ok {
		return out
	}
	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Ptr {
		return out
	}
	elem := rv.Type().Elem()
	rtcache.RLock()
	relem := rtcache.m[elem]
	rtcache.RUnlock()
	if relem == nil {
		rtcache.Lock()
		relem = mapCachedType(elem)
		rtcache.Unlock()
	}
	return reflect.NewAt(relem, unsafe.Pointer(rv.Pointer())).Interface()
}
