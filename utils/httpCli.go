package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"regexp"
)

type HttpCli struct {
	httpClient *http.Client
}
func NewHttpCli() *HttpCli{
	//声明一个jar存放cookie
	jar, _ := cookiejar.New(nil)
	//设置使用短连接
	tran := http.Transport{
		DisableKeepAlives: true,
	}
	//声明一个httpClient 用于自动带着cookieJar内的cookie发送http请求并接收cookie存入cookieJar
	httpClient := &http.Client{
		CheckRedirect: nil,
		Jar: jar,
		Transport: &tran,
	}
	return &HttpCli{
		httpClient: httpClient,
	}
}

// Get 发起HTTPGet请求并返回resp cookie error
func (c *HttpCli) Get(url string) (string, string, error) {
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", "", err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", "", err
	}

	gCurCookieJar := c.httpClient.Jar
	gCurCookies := gCurCookieJar.Cookies(httpReq.URL)
	cookie, err := json.Marshal(gCurCookies)
	if err != nil {
		return "", "", err
	}

	return string(body), string(cookie), nil
}
// Post 发起HTTP Post请求并返回resp cookie error
func (c *HttpCli) Post(url string, param io.Reader) (string, []*http.Cookie, error) {
	httpReq, err := http.NewRequest("POST", url, param)
	if err != nil {
		return "", nil, err
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", nil, err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", nil, err
	}

	gCurCookieJar := c.httpClient.Jar
	gCurCookies := gCurCookieJar.Cookies(httpReq.URL)
	if err != nil {
		return "", nil, err
	}

	return string(body), gCurCookies, nil
}

// PostPic 发起HTTP Post请求并返回resp cookie error
func (c *HttpCli) PostPic(url string, param *bytes.Buffer) (string, []*http.Cookie, error) {
	httpReq, err := http.NewRequest("POST", url, param)
	if err != nil {
		return "", nil, err
	}


	r, err := regexp.Compile("--(.*)\r\n")
	if err != nil {
		fmt.Println("name regexp failed", err)
	}
	boundary := r.FindStringSubmatch(param.String())
	httpReq.Header.Set("Content-Type", "multipart/form-data;boundary=" + boundary[1])
	//httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", nil, err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", nil, err
	}

	gCurCookieJar := c.httpClient.Jar
	gCurCookies := gCurCookieJar.Cookies(httpReq.URL)
	if err != nil {
		return "", nil, err
	}

	return string(body), gCurCookies, nil
}
