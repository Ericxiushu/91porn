package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/zyxar/argo/rpc"
)

func sendToAria2(content Content, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	videoFileName := fmt.Sprintf("%s.mp4", content.title)

	videoFile := path.Join(cfg.SaveDir, videoFileName)
	fmt.Println("download video to ", videoFile)
	if isExist(videoFile) {
		fmt.Printf("video %s is exist \n", videoFile)
		return
	}

	client, err := rpc.New(context.Background(), cfg.Aria2.Uri, cfg.Aria2.Token, time.Second*10, nil)
	if err != nil {
		fmt.Println("rpc error ", err)
		return
	}

	_, err = client.AddURI([]string{content.videoURL}, map[string]string{"out": videoFileName, "dir": cfg.SaveDir})
	if err != nil {
		fmt.Printf("add file %s to aria2 fail : %v , url : %s \n", content.title, err, content.videoURL)
		return
	}

	fmt.Printf("add file %s to aria2 success \n", content.title)

}

/// 下载文件
func downloadFile(url string, fileName string, c chan int) {
	//fileName := getNameFromUrl(url)

	defer func() {
		c <- 0
	}()

	req := buildRequest(url)
	if req == nil {
		fmt.Println("buildRequest error")
		return
	}
	http.DefaultClient.Timeout = 10 * time.Second
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed download , ", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed download " + url)
		return
	}

	defer func() {
		resp.Body.Close()
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	localFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)

	if _, err := io.Copy(localFile, resp.Body); err != nil {
		panic("failed save " + fileName)
	}

	fmt.Println("success download " + fileName)
}

/// 下载页面采集内容
func downloadContent(content Content, c chan int) {
	fmt.Println("begin download " + content.title)

	basePath := fmt.Sprintf("%s/%s/", cfg.SaveDir, content.title)
	basePath = strings.Replace(basePath, "\n", "", -1)
	basePath = strings.Replace(basePath, " ", "", -1)

	//var c1 chan int
	c1 := make(chan int)
	chanCount := 0

	// 创建目录
	if !isExist(basePath) {
		os.MkdirAll(basePath, 0777)
	}

	// 下载缩略图
	thumbFile := basePath + "thumb.jpg"
	if !isExist(thumbFile) {
		chanCount += 1
		go downloadFile(content.thumbURL, thumbFile, c1)
	}

	// 下载视频
	videoFile := basePath + "1.mp4"
	if !isExist(videoFile) {
		chanCount += 1
		go downloadFile(content.videoURL, videoFile, c1)
	}

	for i := 0; i < chanCount; i++ {
		<-c1
	}
	c <- 0
}

/// 获取远端服务器的内容页面
func getContent(url string, content *Content) {

	// 获取远程服务器的页面
	doc, err := getHtml(url, USE_PROXY)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 视频缩略图url
	v := doc.Find("video")
	thumbURL, _ := v.Attr("poster")
	if len(thumbURL) > 0 {
		content.thumbURL = thumbURL
		//fmt.Println(thumbURL)
	}

	// 视频url
	src := v.Find("source")
	videoURL, _ := src.Attr("src")
	if len(videoURL) > 0 {
		content.videoURL = videoURL
		//fmt.Println(videoURL)
	}

	// 从script拿
	if len(content.videoURL) < 1 {
		content.videoURL = fetchFromJS(v.Text())
	}

	// 从分享链接拿
	if len(content.videoURL) < 1 {
		shareURL := doc.Find("#linkForm2 #fm-video_link").Text()
		getContent(shareURL, content)
	}

	if len(content.videoURL) < 1 {
		fmt.Printf("get %s video url failure \n", url)
	}

	// 标题
	t := doc.Find("div#viewvideo-title")
	title := t.Text()
	title = strings.TrimSpace(title)
	if len(title) > 0 {
		content.title = title
		//fmt.Println(title)
	}
}
