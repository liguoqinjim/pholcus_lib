package netease_music

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	//"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	"regexp"
	"strconv"
	"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
	//"fmt"
	"fmt"
	"io/ioutil"
	"log"
)

func init() {
	NMSpider.Register()
}

//var artist_cat_ids = []int{1001}

var artist_cat_ids = []int{1001, 1002, 1003}

var NMSpider = &Spider{
	Name:        "网易云音乐",
	Description: "网易云音乐 [music.163.com]",
	// Pausetime: 300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			//Aid调用Rule中的AidFunc
			//ctx.Aid(map[string]interface{}{"Rule": "抓取歌手"}, "抓取歌手")

			//抓取专辑
			ctx.Aid(map[string]interface{}{"Rule": "判断专辑页数"}, "判断专辑页数")
		},

		Trunk: map[string]*Rule{
			"抓取歌手": {
				ItemFields: []string{
					"歌手",
					"歌手主页",
				},
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//查找歌手的url http://music.163.com/discover/artist/cat?id=1001&initial=90
					for _, v := range artist_cat_ids {
						for i := 65; i <= 90; i++ {
							url := fmt.Sprintf("http://music.163.com/discover/artist/cat?id=%d&initial=%d", v, i)
							ctx.AddQueue(
								&request.Request{
									Url:  url,
									Rule: aid["Rule"].(string),
								},
							)
						}

					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					queryString := string(ctx.GetText())
					err := ioutil.WriteFile("test.txt", []byte(queryString), 0644)
					if err != nil {
						log.Fatal(err)
					}

					query.Find(".nm.nm-icn.f-thide.s-fc0").Each(func(i int, s *goquery.Selection) {
						artistName := s.Text()
						artistUrl, _ := s.Attr("href")

						ctx.Output(map[int]interface{}{
							0: artistName,
							1: artistUrl,
						})
					})

					//ctx.Output(map[int]interface{}{
					//	0: queryString,
					//})
				},
			},

			"判断专辑页数": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//url
					url := "http://music.163.com/artist/album?id=13193&limit=12&offset=0"
					ctx.AddQueue(
						&request.Request{
							Url:  url,
							Rule: aid["Rule"].(string),
						},
					)
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					pageCount := query.Find(".zpgi").Size()
					ctx.Aid(map[string]interface{}{"PageCount": pageCount, "Rule": "抓取专辑"}, "抓取专辑")
				},
			},

			"抓取专辑": {
				ItemFields: []string{
					"专辑",
					"专辑链接",
				},
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//url http://music.163.com/artist/album?id=13193&limit=12&offset=0
					pageCount := aid["PageCount"].(int)

					for i := 0; i < pageCount; i++ {
						url := fmt.Sprintf("http://music.163.com/artist/album?id=13193&limit=12&offset=%d", 12*i)
						ctx.AddQueue(
							&request.Request{
								Url:  url,
								Rule: aid["Rule"].(string),
							},
						)
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					queryString := string(ctx.GetText())
					err := ioutil.WriteFile("test.txt", []byte(queryString), 0644)
					if err != nil {
						log.Fatal(err)
					}

					query.Find(".tit.f-thide.s-fc0").Each(func(i int, s *goquery.Selection) {
						albumName := s.Text()
						albumUrl, _ := s.Attr("href")

						ctx.Output(map[int]interface{}{
							0: albumName,
							1: albumUrl,
						})
					})
				},
			},

			"生成请求": {
				//单数页是url直接返回,双数页是异步加载,两个url在下面有写
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//Url:  "http://search.jd.com/Search?keyword=" + ctx.GetKeyin() + "&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&bs=1&s=1&click=0&page=" + strconv.Itoa(pageNum),
					//Url:  "http://search.jd.com/s_new.php?keyword=" + ctx.GetKeyin() + "&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&bs=1&s=31&scrolling=y&pos=30&page=" + strconv.Itoa(pageNum),
					pageCount := aid["PageCount"].(int)

					for i := 1; i < pageCount; i++ {
						ctx.AddQueue(
							&request.Request{
								Url:  "http://search.jd.com/Search?keyword=" + ctx.GetKeyin() + "&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&bs=1&s=1&click=0&page=" + strconv.Itoa(i*2-1),
								Rule: "搜索结果",
							},
						)
						ctx.AddQueue(
							&request.Request{
								Url:  "http://search.jd.com/s_new.php?keyword=" + ctx.GetKeyin() + "&enc=utf-8&qrst=1&rt=1&stop=1&vt=2&bs=1&s=31&scrolling=y&pos=30&page=" + strconv.Itoa(i*2),
								Rule: "搜索结果",
							},
						)
					}
					return nil
				},
			},

			"搜索结果": {
				//从返回中解析出数据。注：异步返回的结果页面结构是和单数页的一样的，所以就一套解析就可以了。
				ItemFields: []string{
					"标题",
					"价格",
					"评论数",
					"链接",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					query.Find(".gl-item").Each(func(i int, s *goquery.Selection) {
						// 获取标题
						a := s.Find(".p-name.p-name-type-2 > a")
						title := a.Text()

						re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
						// title = re.ReplaceAllStringFunc(title, strings.ToLower)
						title = re.ReplaceAllString(title, " ")
						title = strings.Trim(title, " \t\n")

						// 获取价格
						price, _ := s.Find("strong[data-price]").First().Attr("data-price")

						// 获取评论数
						//#J_goodsList > ul > li:nth-child(1) > div > div.p-commit
						discuss := s.Find(".p-commit > strong > a").Text()

						// 获取URL
						url, _ := a.Attr("href")
						url = "http:" + url

						// 结果存入Response中转
						if title != "" {
							ctx.Output(map[int]interface{}{
								0: title,
								1: price,
								2: discuss,
								3: url,
							})
						}
					})
				},
			},
		},
	},
}
