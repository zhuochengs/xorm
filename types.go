package xorm

import (
	"github.com/zhuochengs/xorm-core"
	"reflect"
)

var (
	ptrPkType = reflect.TypeOf(&core.PK{})
	pkType    = reflect.TypeOf(core.PK{})
)
