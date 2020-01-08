//author: https://github.com/zhaojunlike
//date: 2019/12/15
package common

import "time"

func CreateTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func CreateUnix() int64 {
	return time.Now().Unix()
}
