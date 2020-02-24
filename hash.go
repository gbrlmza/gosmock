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
		fmt.Fprint(digester, k)
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	return string(digester.Sum(nil))
}
