package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/proxy"
)

func initJs() {

	ss, err := fetchHtml(cfg.Md5js, false)
	if err != nil {
		panic(fmt.Sprintf("get md5js error : %v", err))
	}

	defer ss.Body.Close()

	jsstr, err := ioutil.ReadAll(ss.Body)
	if err != nil {
		panic(fmt.Sprintf("get md5js content error : %v", err))
	}

	jsContent = string(jsstr)

}

/// 转换字符串编码
func convertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//根据url 创建http 请求的 request
//网站有反爬虫策略 wireshark 不解释
func buildRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req.Header.Set("Accept-Language", `zh-CN,zh;q=0.9`)
	req.Header.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.106 Safari/537.36`)
	req.Header.Set("X-Forwarded-For", random_ip())
	req.Header.Set("referer", url)
	req.Header.Set("Content-Type", `multipart/form-data; session_language=cn_CN`)

	return req
}

func fetchHtml(url string, useProxy bool) (*http.Response, error) {

	var (
		resp *http.Response
		err  error
	)

	req := buildRequest(url)

	if useProxy {
		dialSocksProxy, err := proxy.SOCKS5("tcp", "localhost:1087", nil, proxy.Direct)

		tr := &http.Transport{Dial: dialSocksProxy.Dial}

		client := &http.Client{
			Transport: tr,
			Timeout:   time.Second * 5,
		}

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

	} else {

		client := &http.Client{
			Timeout: time.Second * 5,
		}

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("get http code error %d", resp.StatusCode)
	}

	return resp, nil

}

/// 获取远端服务器的HTML页面
func getHtml(url string, useProxy bool) (*goquery.Document, error) {

	res, err := fetchHtml(url, useProxy)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromResponse(res)

}

func fetchFromJS(jsStr string) string {

	jsStr = strings.ReplaceAll(jsStr, "document.write", "")
	jsStr = strings.ReplaceAll(jsStr, "<!--", "")
	jsStr = strings.ReplaceAll(jsStr, "//-->", "")

	vm := otto.New()

	jsContent := fmt.Sprintf(`%s;%s`, jsContent, jsStr)
	// fmt.Println(jsContent)

	v, err := vm.Run(jsContent)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	dd := `<source [^>]*src=['"](?P<src>([^'"]+))[^>]*>`
	reg := regexp.MustCompile(dd)

	if reg.NumSubexp() == 2 {
		list := reg.FindStringSubmatch(v.String())
		if len(list) >= 2 {
			return list[1]
		}

		fmt.Printf("FindStringSubmatch error , list %v \n", list)

		return ""
	}

	return ""

}
