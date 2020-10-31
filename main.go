// CGO_ENABLED=0 go build -o 91porn

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	USE_PROXY bool = false ///< 是否启用代理
	VERSION        = "v0.0.3"
)

var (
	saveDir = ""

	jsContent = ""

	cfg = &Config{}

	defaultConfg = `{"port":8888,"domain":"1003.workarea8.live","md5js":"http://1003.workarea8.live/js/m.js","saveDir":"/data/91porn","aria2":{"uri":"http://127.0.0.1:6800/jsonrpc","token":"aabbccdd"},"pageType":"rf","pages":{"hot":"当前最热","rp":"最近得分","long":"10分钟以上","md":"本月讨论","tf":"本月收藏","mf":"收藏最多","rf":"最近加精","top":"本月最热"}}`
)

type Config struct {
	Domain   string            `json:"domain"`
	SaveDir  string            `json:"saveDir"`
	Port     int               `json:"port"`
	Md5js    string            `json:"md5js"`
	PageType string            `json:"pageType"`
	Pages    map[string]string `json:"pages"`
	Aria2    Aria2Cfg          `json:"aria2"`
}

type Aria2Cfg struct {
	Uri   string `json:"uri"`
	Token string `json:"token"`
}

type Content struct {
	title      string ///< 标题
	desc       string ///< 描述
	contentURL string ///< 内容入口地址
	coverURL   string ///< 封面地址
	thumbURL   string ///< 缩略图地址
	videoURL   string ///< 视频地址
}

func init() {

	rand.Seed(time.Now().Unix())

	initConf()
	initJs()
	initDir()

}

func initDir() {
	ok, err := PathExists(cfg.SaveDir)
	if err != nil {
		panic(err)
	}

	if !ok {
		err = os.MkdirAll(cfg.SaveDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func initConf() {

	res, err := ioutil.ReadFile("/91porn/conf/conf.json")
	if err == nil {
		err = json.Unmarshal(res, cfg)
		if err == nil {
			return
		}
	}

	fmt.Println("get config error , use default config")
	err = json.Unmarshal([]byte(defaultConfg), cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("load conf succ")

}

/// 获取远端服务器的列表页面
func getPage(pageURL string, contents *[]Content) {
	// 获取远程服务器的页面
	fmt.Println("start get list page")
	doc, err := getHtml(pageURL, USE_PROXY)
	if err != nil {
		fmt.Println(err)
		return
	}

	var content Content

	// 获取内容页面的访问入口url
	doc.Find(".row a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		title, err := item.ChildrenFiltered("span").Html()
		if err != nil {
			fmt.Println("get title error : ", err)
			return
		}
		if title != "" {
			content = Content{
				contentURL: link,
				title:      title,
			}
			*contents = append(*contents, content)
		}
	})

	// 遍历内容页面
	for k, v := range *contents {
		fmt.Printf("fetch %d detail page : %s\n", k, v.contentURL)
		getContent(v.contentURL, &(*contents)[k])
	}

}

/// 爬虫
func spider(pageType string, pageNums ...int) string {

	if len(pageType) <= 0 {
		pageType = cfg.PageType
	}

	if len(pageNums) < 1 {
		pageNums = []int{1}
	}

	// 抓取页面
	var contents []Content
	for _, i := range pageNums {
		pageURL := fmt.Sprintf("http://%s/v.php?category=%s&viewtype=basic&page=%d", cfg.Domain, pageType, i)
		//fmt.Println(pageURL)
		getPage(pageURL, &contents)
	}

	var res strings.Builder

	// 下载视频
	wg := &sync.WaitGroup{}
	for _, s := range contents {

		wg.Add(1)
		go sendToAria2(s, wg)

		res.WriteString("============ " + VERSION + " ===============\n")
		res.WriteString(fmt.Sprintf("页面地址:%s\n", s.contentURL))
		res.WriteString(fmt.Sprintf("标题:%s\n", s.title))
		res.WriteString(fmt.Sprintf("缩略图:%s\n", s.thumbURL))
		res.WriteString(fmt.Sprintf("视频:%s\n", s.videoURL))

		res.WriteString("\n")

	}

	wg.Wait()

	fmt.Println(res.String())
	fmt.Println("all done!")

	return res.String()
}

func cron() {

	ti := time.NewTicker(time.Hour)

	fmt.Println("start cron : ", time.Now())
	spider("rp", 1)

	for {
		select {
		case <-ti.C:

			fmt.Println("start cron : ", time.Now())
			spider("rp", 1)

		}
	}

}

func main() {

	go cron()
	web()

}
