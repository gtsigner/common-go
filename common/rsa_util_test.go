package common

import (
	"fmt"
	"testing"
)

func TestRsaEncryptString(t *testing.T) {
	var res, err = RsaEncryptString("846600Es", "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDe3Dg7/mmAK8gwjzj49zlpvMer7vMjagHANU4tv26yjtebNw/tffj9eXw4FQQ95liZEIooDWUnjbs2KRBRM8K8l/8wh6hT/TqDpr82LA+Pa716vk7IpsFqW5qOyZavM9RGHeqM7DC4weeoh+yIvpCXQBx67xO2Ce89jR55Xl5uQQIDAQAB")
	fmt.Println(res, err)
}
func TestNewRsa(t *testing.T) {
	var mapper = make(map[string]interface{}, 0)
	mapper["a"] = "123"
	fmt.Println(mapper)
}
