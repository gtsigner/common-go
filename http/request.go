//author: https://github.com/zhaojunlike
//date: 2019/12/11
//enable global dns cached
package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/viki-org/dnscache"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func NewDefaultHttpRes() *Res {
	res := Res{Code: "0", Message: "", Data: nil, Ok: false, StatusCode: 599,}
	return &res
}

type Client struct {
	Http    *http.Client
	Options *Options
}

func (client *Client) Destroy() {
	client.Http.CloseIdleConnections()
}

//发送请求
func (client *Client) Request(config *Config) (*Res, error) {
	return SendHttpRequest(client.Http, config)
}

//set cookirjar
func (client *Client) SetCookieJar(jar *cookiejar.Jar) {
	client.Options.CookieJar = jar
	client.Http.Jar = client.Options.CookieJar
}

// 设置单个cookiejar的cookie数据
func (client *Client) SetCookie(uri string, cookie *http.Cookie) error {
	var r, err = url.Parse(uri)
	if err != nil {
		return err
	}
	var arr = make([]*http.Cookie, 1)
	arr[0] = cookie
	if client.Http.Jar != nil {
		client.Http.Jar.SetCookies(r, arr)
		return nil
	}
	return ErrNotSetCookieJar
}

// 移除COokie
func (client *Client) RemoveCookie(uri string, name string) error {
	c := &http.Cookie{Name: name, Value: "", Expires: time.Unix(0, 0)}
	return client.SetCookie(uri, c)
}

// 移除C0okie
func (client *Client) GetCookie(uri string, name string) *http.Cookie {
	var urs, _ = url.Parse(uri) //
	var cookies = client.Http.Jar.Cookies(urs)
	for _, v := range cookies {
		if v.Name == name {
			return v
		}
	}
	return nil
}

var (
	resolver = dnscache.New(time.Minute * 5)
)

//构造HttpClient
//TODO 提供tsConf进行配置
func NewHttpClient(opt *Options) (*Client, error) {
	tsConf := &tls.Config{
		NextProtos:             []string{"http/1.1", "h2"},
		SessionTicketsDisabled: false,
		//InsecureSkipVerify: true,
		//PreferServerCipherSuites: true,
		//MinVersion: tls.VersionTLS10,
		//MaxVersion: tls.VersionTLS12,
		//CipherSuites: []uint16{
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_RSA_WITH_RC4_128_SHA,
		//	tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		//	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		//	tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		//	tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		//},
	}

	//pool := x509.NewCertPool()
	//certData, err := ioutil.ReadFile("../certs/ca.cer")
	//fmt.Println(err)
	//pool.AppendCertsFromPEM(certData)
	//tsConf.RootCAs = pool

	tr := &http.Transport{
		TLSClientConfig:     tsConf,
		MaxIdleConnsPerHost: 1024,
		//dns cached
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			separator := strings.LastIndex(addr, ":")
			ip, _ := resolver.FetchOneString(addr[:separator])
			return net.Dial("tcp", ip+addr[separator:])
		},
	}
	if opt.Proxy != nil {
		uri := url.URL{}
		var proxyUrl *url.URL
		var urlStr = "http://"
		if opt.Proxy.Security {
			urlStr = "https://"
		}
		if opt.Proxy.Username != "" && opt.Proxy.Password != "" {
			opt.Proxy.NeedAuth = true
			urlStr += fmt.Sprintf("%s:%s@", opt.Proxy.Username, opt.Proxy.Password)
		}
		urlStr += fmt.Sprintf("%s:%d", opt.Proxy.Host, opt.Proxy.Port)
		proxyUrl, _ = uri.Parse(urlStr)
		if proxyUrl != nil {
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	if opt.UseHttp2 {
		if err := http2.ConfigureTransport(tr); err != nil {
			return nil, err
		}
	}

	if opt.Timeout <= 0 {
		opt.Timeout = 30 * time.Second
	}
	client := Client{Options: opt}
	client.Http = &http.Client{Transport: tr, Timeout: opt.Timeout}
	if opt.CookieJar != nil {
		client.Http.Jar = opt.CookieJar
	}
	//禁止自动转发
	if opt.AllowRedirect == false {
		client.Http.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return &client, nil
}

//http 请求
func SendHttpRequest(client *http.Client, config *Config) (*Res, error) {
	ts := time.Now()

	//#region start
	res := NewDefaultHttpRes()
	res.Req = Req{Config: config}
	res.Url = config.Url
	var req *http.Request
	var err error
	//set default method
	if config.Method == "" {
		config.Method = "GET"
	}
	//useragent
	if config.Headers == nil {
		config.Headers = map[string]string{}
	}
	//TODO 全部转小写的headers
	if _, ext := config.Headers["User-Agent"]; ext == false {
		config.Headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36"
	}

	if config.Data == nil {
		//GET 之类的
		req, err = http.NewRequest(config.Method, config.Url, nil)
	} else {
		//包含数据的请求
		var ct = "application/json"
		//默认
		var v, ext = config.Headers["Content-Type"]
		if ext {
			ct = v
		}
		if config.PostForm {
			//TODO 判断类型，自动解析成form

			var str = config.Data.(url.Values).Encode()
			req, err = http.NewRequest(config.Method, config.Url, strings.NewReader(str))
			if err == nil {
				ct = "application/x-www-form-urlencoded"
			}
		} else {
			var str, err = json.Marshal(config.Data)
			if err == nil {
				req, err = http.NewRequest(config.Method, config.Url, bytes.NewReader(str))
			}
		}
		//设置请求
		if err == nil {
			config.Headers["Content-Type"] = ct
		}
	}
	//整体判断错误
	if err != nil {
		res.Message = err.Error()
		return res, err
	}

	//便利header
	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	rsp, err := client.Do(req)
	if err != nil {
		res.Message = fmt.Sprintf("error getting: %v\n", err)
		return res, nil
	}
	defer rsp.Body.Close()
	//status
	res.StatusCode = rsp.StatusCode
	res.Ok = rsp.StatusCode < 300
	res.Proto = rsp.Proto
	//解析headers
	res.Headers = rsp.Header

	//解析Cookies DEFAULT
	if config.AutoParseCookie == true {
		res.Cookies = rsp.Cookies()
	}

	//body
	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		res.Message = fmt.Sprintf("error reading response: %v\n", err)
		return res, nil
	}
	res.RespStr = string(buf)
	ct := rsp.Header["Content-Type"]

	//TODO 优化一下逻辑，根据content-type 进行自动解析，或者通过jsonparse 强制解析
	//全部尝试用json解析
	var obj interface{}
	if nil == json.Unmarshal(buf, &obj) {
		res.Data = obj
	}

	if config.ParseJson == false && ct != nil && len(ct) > 0 {
		cty := ct[0]
		var obj interface{}
		//判断是否需要解析成JSON
		if strings.Index(cty, "json") != -1 && nil == json.Unmarshal(buf, &obj) {
			res.Data = obj
		}
	}
	//#endregion

	te := time.Now()
	tw := te.Sub(ts)
	res.Time = int(tw.Milliseconds())

	return res, nil
}
