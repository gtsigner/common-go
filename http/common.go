package http

import (
	"strconv"
	"strings"
)

//解析proxy
func ParseProxy(str string) *Proxy {
	str = strings.Trim(str, " ")
	var proxy = &Proxy{}
	if str == "" {
		return nil
	}
	var arr = strings.Split(str, ":")
	if len(arr) == 2 {
		proxy.Host = arr[0]
		var port, err = strconv.Atoi(arr[1])
		if err != nil {
			return nil
		}
		proxy.Port = port
		return proxy
	}
	//如果有用户名或者密码
	if len(arr) == 4 {
		proxy.Host = arr[0]
		var port, err = strconv.Atoi(arr[1])
		if err != nil {
			return nil
		}
		proxy.Port = port
		proxy.Username = arr[2]
		proxy.Password = arr[3]
		return proxy
	}
	return nil

}
