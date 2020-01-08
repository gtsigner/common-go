//author: https://github.com/zhaojunlike
//date: 2019/12/11
package common

import (
	"container/list"
	"encoding/json"
	funk "github.com/thoas/go-funk"
)

func UniqString(arr []string) []string {
	return funk.UniqString(arr)
}

//转array
func ListToArray(li *list.List) []interface{} {
	var ls = make([]interface{}, li.Len())
	var ix = 0
	for e := li.Front(); e != nil; e = e.Next() {
		ls[ix] = e.Value
		ix++
	}
	return ls
}

func JSONStringify(v interface{}) string {
	var bytes, _ = json.Marshal(v)
	return string(bytes[:])
}

// 打印列表
func LoopList(l list.List, call func(index int, e *list.Element)) {
	var line = 0
	for e := l.Front(); e != nil; e = e.Next() {
		call(line, e)
		line++
	}
}
