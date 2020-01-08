package convert

import "fmt"

func Interface2String(inter interface{}) string {
	var str = fmt.Sprintf("%v", inter)
	return str
}
