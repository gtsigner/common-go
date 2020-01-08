//author: https://github.com/zhaojunlike
//date: 2019/12/13
package http

import (
	"fmt"
	"testing"
)

func TestMapToQueryString(t *testing.T) {
	var pa = map[string]interface{}{
		"a": 1,
		"b": true,
		"c": "123",
	}
	var str = MapToQueryString(pa)
	fmt.Println(str)
}
