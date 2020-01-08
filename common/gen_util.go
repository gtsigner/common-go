//author: https://github.com/zhaojunlike
//date: 2019/12/17
package common

import (
	uuid "github.com/satori/go.uuid"
)

func CreateUUIDV4() string {
	u1 := uuid.NewV4()
	return u1.String()
}
