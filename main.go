package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"movie-crawler/page"
	"net"
	"net/http"
	"time"
)

func main() {
	c := colly.NewCollector(
		// 设置Ua
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
	)

	// 设置代理和超时
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		detailParser := page.New()
		_, movie := detailParser.ParseDetail(r)
		fmt.Printf("解析得到的信息：%v", movie)

	})

	c.Visit("https://movie.douban.com/subject/26933210/")
}
