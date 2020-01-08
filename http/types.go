//author: https://github.com/zhaojunlike
//date: 2019/12/11
package http

import (
    "crypto/tls"
    "errors"
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "time"
)

//全局HTTP ERROR
var (
    ErrParseProxy      = errors.New("parse proxy fail")
    ErrNotSetCookieJar = errors.New("not set cookie jar")
)

type Options struct {
    Proxy         *Proxy
    Timeout       time.Duration
    AllowRedirect bool
    CookieJar     *cookiejar.Jar
    tslConfig     *tls.Config
    UseHttp2      bool
    Debug         bool
}

func NewOptions() *Options {
    return &Options{Proxy: nil, Timeout: 30 * time.Second, AllowRedirect: false, UseHttp2: true}
}

type Proxy struct {
    Host     string
    Port     int
    Username string
    Password string
    NeedAuth bool
    Security bool //use https
}

//请求
type Req struct {
    Config *Config
}

//响应
type Res struct {
    Code       string
    Message    string
    Data       interface{}
    Ok         bool
    RespStr    string
    Time       int
    StatusCode int
    Url        string
    Proto      string
    Headers    map[string][]string
    Cookies    []*http.Cookie
    Req        Req
    Success    bool
}

func (res *Res) Println() {
    PrintlnHttpRes(*res)
}
func (res *Res) PrintData() {
    PrintlnResData(*res)
}

func (res *Res) Destroy() {

}
func PrintlnHttpRes(res Res) {
    fmt.Printf("URL：%s\t耗时:%dms\t状态码:%d\n", res.Url, res.Time, res.StatusCode)
}
func PrintlnResData(res Res) {
    fmt.Printf("URL：%s\t耗时:%dms\t状态码:%d\n返回:%s\n", res.Url, res.Time, res.StatusCode, res.RespStr)
}
