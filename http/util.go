//author: https://github.com/zhaojunlike
//date: 2019/12/13
package http

import (
	"fmt"
	"sort"
	"strings"
)

func SortMap(mapper map[string]interface{}) map[string]interface{} {
	var keys = make([]string, 0)
	for k, _ := range mapper {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sorted = make(map[string]interface{}, 0)
	for _, v := range keys {
		sorted[v] = mapper[v]
	}
	return sorted
}
func MapToQueryString(mapper map[string]interface{}) string {
	var str = ""
	//根据key排序
	for k, v := range mapper {
		//判断
		str += fmt.Sprintf("&%s=%v", k, v)
	}
	return strings.Replace(str, "&", "", 1)
}

// 排序map
func MapToQueryStringSort(mapper map[string]interface{}) string {
	var str = ""
	var keys = make([]string, 0)
	for k, _ := range mapper {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//根据key排序
	for _, k := range keys {
		//判断
		str += fmt.Sprintf("&%s=%v", k, mapper[k])
	}
	return strings.Replace(str, "&", "", 1)
}

func MapToSplit(mapper map[string]string) string {
	var str = ""
	var char = ";"
	for k, v := range mapper {
		//判断
		str += fmt.Sprintf("%s=%v%s", k, v, char)
	}
	return str
}
