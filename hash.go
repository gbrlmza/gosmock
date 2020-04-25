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
		kind := reflect.TypeOf(ob).Kind()
		obStr := fmt.Sprint(ob)
		if kind == reflect.Ptr {
			obStr = reflect.ValueOf(ob).Elem().String()
		}
		fmt.Fprint(digester, k)
		fmt.Fprint(digester, kind)
		fmt.Fprint(digester, obStr)
	}
	bytes := string(digester.Sum(nil))
	str := fmt.Sprintf("%x", bytes)
	return str
}
