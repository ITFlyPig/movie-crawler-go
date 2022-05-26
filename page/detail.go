package page

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

// MovieInfo 影片信息
type MovieInfo struct {
	Name        string   // 名字
	Year        string   // 出品时间
	Cover       string   // 封面
	Director    string   // 导演
	Writers     []string // 编剧
	Actors      []string // 主演
	MovieType   []string // 类型
	Country     string   // 国家
	Language    []string // 语言
	ReleaseDate []string // 上映日期
	Mins        int      // 时长
	Alias       []string // 别名
}

type IMovieDetail interface {
	// ParseDetail 解析得到具体的电影的信息
	ParseDetail(response *colly.Response) (error, *MovieInfo)
}

//====================从豆瓣获取电影的详细信息============================

// DoubanDetailParser 豆瓣解析具体的电影信息
type DoubanDetailParser struct {
	// itemHandler 不同item的处理器
	itemHandler map[string]func(value string, movie *MovieInfo)
}

func New() *DoubanDetailParser {
	return &DoubanDetailParser{
		itemHandler: map[string]func(value string, movie *MovieInfo){
			"国家": func(value string, movie *MovieInfo) {
				movie.Country = value
			},
			"类型": func(value string, movie *MovieInfo) {
				movie.MovieType = split(value)
			},
			"语言": func(value string, movie *MovieInfo) {
				movie.Language = split(value)
			},
			"上映": func(value string, movie *MovieInfo) {
				movie.ReleaseDate = split(value)
			},
			"片长": func(value string, movie *MovieInfo) {
				minStr := strings.ReplaceAll(value, "分钟", "")
				minStr = strings.ReplaceAll(minStr, " ", "")
				mins, err := strconv.Atoi(minStr)
				if err != nil {
					fmt.Printf("电影 %v 将时长 %v 转为分钟失败：%v", movie.Name, value, err)
					return
				}
				movie.Mins = mins
			},
			"又名": func(value string, movie *MovieInfo) {
				movie.Alias = split(value)
			},
			"导演": func(value string, movie *MovieInfo) {
				movie.Director = value
			},
			"编剧": func(value string, movie *MovieInfo) {
				movie.Writers = split(value)
			},
			"主演": func(value string, movie *MovieInfo) {
				movie.Actors = split(value)
			},
		},
	}

}

// split 去除空格，并使用/分割为数组
func split(str string) []string {
	newStr := strings.ReplaceAll(str, " ", "")
	return strings.Split(newStr, "/")
}

func (d *DoubanDetailParser) ParseDetail(response *colly.Response) (err error, movie *MovieInfo) {
	if response == nil {
		err = errors.New("响应Response不能为空")
		return
	}
	if response.StatusCode != 200 {
		err = errors.New("响应码不为200")
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
	if err != nil {
		return
	}

	movie = new(MovieInfo)
	// 获取封面
	cover, exists := doc.Find("#mainpic>a>img").Attr("src")
	if exists {
		movie.Cover = cover
	}
	// 获取名字
	name := doc.Find("#content > h1 > span:nth-child(1)").Text()
	movie.Name = name
	// 获取其他信息
	infoText := doc.Find("#info").Text()
	arr := strings.Split(infoText, "\n")
	for _, item := range arr {
		for key, f := range d.itemHandler {
			if strings.Contains(item, key) {
				start := strings.Index(item, ":")
				total := len(item)
				if start >= 0 && start+1 < total {
					value := item[start+1 : total]
					f(value, movie)
				}

			}
		}
	}
	return
}
