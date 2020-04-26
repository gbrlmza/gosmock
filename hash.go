package gosmock

import (
	"crypto"
	_ "crypto/md5"
	"fmt"
	"reflect"
)

func hash(objs ...interface{}) string {
	digester := crypto.MD5.New()
	for k, ob := range objs {
		typeOf := reflect.TypeOf(ob)
		obStr := fmt.Sprint(ob)
		kindStr := ""

		if isNullable(reflect.TypeOf(ob)) && (ob == nil || reflect.ValueOf(ob).IsNil()) {
			kindStr = "nil"
			obStr = "nil"
		} else {
			kind := typeOf.Kind()
			kindStr = kind.String()
			if kind == reflect.Ptr {
				obStr = reflect.ValueOf(ob).Elem().String()
			}
		}

		fmt.Fprint(digester, k)
		fmt.Fprint(digester, kindStr)
		fmt.Fprint(digester, obStr)
	}

	bytes := string(digester.Sum(nil))
	str := fmt.Sprintf("%x", bytes)
	return str
}
