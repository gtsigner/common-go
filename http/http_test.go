//author: https://github.com/zhaojunlike
//date: 2019/12/12
package http

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewDefaultHttpRes(t *testing.T) {
	var opt = &Options{Proxy: &Proxy{Host: "127.0.0.1", Port: 8888}}
	client, _ := NewHttpClient(opt)
	defer client.Destroy()
	ids := []string{"d8469d31-ca22-474b-a329-450d32adc789", "d8469d31-ca22-474b-a329-450d32adc789", "d8469d31-ca22-474b-a329-450d32adc789"}
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		uri := fmt.Sprintf("https://api.nike.com/launch/launch_views/v2/%s", id)
		conf := NewConfig(uri)
		res, _ := client.Request(conf)
		res.Println()
		wg.Done()
	}
	wg.Wait()
	fmt.Println("ressss")
}
