//author: https://github.com/zhaojunlike
//date: 2019/12/19
package common

import (
	"net/http"
	"strings"
)

//解析setCookie
func ParseCookiesFromHeader(sk []string) map[string]http.Cookie {
	var mapper = map[string]http.Cookie{}
	for _, v := range sk {
		var ck = ParseSetCookieSingle(v)
		mapper[ck.Name] = ck
	}
	return mapper
}

func ParseSetCookieSingle(str string) http.Cookie {
	var ss = strings.Split(str, ";")[0]
	var kv = strings.Split(ss, "=")
	var ck = http.Cookie{}
	ck.Name = kv[0]
	ck.Value = kv[1]
	return ck
}

func PackCookiesToString(cks map[string]http.Cookie) string {
	var str = ""
	for k, v := range cks {
		str += k + "=" + v.Value + ";"
	}
	return str
}
