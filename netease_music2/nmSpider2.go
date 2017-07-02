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

	//json解析
	"github.com/bitly/go-simplejson"

	// 字符串处理包
	"regexp"
	"strconv"
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
	Name:        "网易云音乐2",
	Description: "网易云音乐2 [music.163.com]",
	// Pausetime: 300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			//Aid调用Rule中的AidFunc
			//ctx.Aid(map[string]interface{}{"Rule": "抓取歌手"}, "抓取歌手")

			//抓取专辑
			//ctx.Aid(map[string]interface{}{"Rule": "判断专辑页数"}, "判断专辑页数")

			//抓取专辑中的歌曲
			//ctx.Aid(map[string]interface{}{"Rule": "抓取专辑中的歌曲"}, "抓取专辑中的歌曲")

			//抓取单独的歌曲 抓取歌曲
			//ctx.Aid(map[string]interface{}{"Rule": "抓取歌曲"}, "抓取歌曲")

			//抓取歌曲评论
			ctx.Aid(map[string]interface{}{"Rule": "抓取歌曲评论"}, "抓取歌曲评论")

			//测试
			//ctx.Aid(map[string]interface{}{"Rule": "测试"}, "测试")
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
					"专辑id",
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

					albums := make(map[int]string) //key是albumid,value是专辑名称
					query.Find(".tit.f-thide.s-fc0").Each(func(i int, s *goquery.Selection) {
						albumName := s.Text()
						albumUrl, _ := s.Attr("href")

						re, _ := regexp.Compile("[1-9]\\d*")
						albumId := 0
						if re.MatchString(albumUrl) {
							albumId, _ = strconv.Atoi(re.FindString(albumUrl))
						}
						albums[albumId] = albumName

						ctx.Output(map[int]interface{}{
							0: albumName,
							1: albumId,
							2: albumUrl,
						})
					})

					ctx.Aid(map[string]interface{}{"Albums": albums}, "抓取专辑中的歌曲")
				},
			},

			"抓取专辑中的歌曲": {
				ItemFields: []string{
					"歌名",
					"歌曲id",
					"歌曲alias",
				},
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					albums := aid["Albums"].(map[int]string)

					for k, _ := range albums {
						url := fmt.Sprintf("http://music.163.com/album?id=%d", k)
						ctx.AddQueue(
							&request.Request{
								Url:  url,
								Rule: "抓取专辑中的歌曲",
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

					div := query.Find("#song-list-pre-cache").Find("textarea")
					song_list := div.Text()
					json1, err := simplejson.NewJson([]byte(song_list))
					if err != nil {
						log.Fatal(err)
					}
					json2, err := json1.Array()
					if err != nil {
						log.Fatal(err)
					}

					for i := 0; i < len(json2); i++ {
						songName, err := json1.GetIndex(i).Get("name").String()
						if err != nil {
							log.Fatal(err)
						}
						songId, err := json1.GetIndex(i).Get("id").Int()
						if err != nil {
							log.Fatal(err)
						}
						songAliases, err := json1.GetIndex(i).Get("alias").StringArray()
						songAlias := ""
						for n, v := range songAliases {
							songAlias += v
							if n != len(songAliases)-1 {
								songAlias += "|"
							}
						}

						ctx.Output(map[int]interface{}{
							0: songName,
							1: songId,
							2: songAlias,
						})
					}
				},
			},

			"抓取歌曲": {
				//ItemFields: []string{
				//	"歌手",
				//	"歌手主页",
				//},
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					ctx.AddQueue(
						&request.Request{
							Url:  "http://music.163.com/song?id=27566765",
							Rule: aid["Rule"].(string),
						},
					)
					return nil
				},
				ParseFunc: func(ctx *Context) {
					//query := ctx.GetDom()

					queryString := string(ctx.GetText())
					err := ioutil.WriteFile("test.txt", []byte(queryString), 0644)
					if err != nil {
						log.Fatal(err)
					}

					//query.Find(".nm.nm-icn.f-thide.s-fc0").Each(func(i int, s *goquery.Selection) {
					//	artistName := s.Text()
					//	artistUrl, _ := s.Attr("href")
					//
					//	ctx.Output(map[int]interface{}{
					//		0: artistName,
					//		1: artistUrl,
					//	})
					//})
				},
			},

			"抓取歌曲评论": {
				//ItemFields: []string{
				//	"歌手",
				//	"歌手主页",
				//},
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//查找歌手的url http://music.163.com/discover/artist/cat?id=1001&initial=90
					//url := fmt.Sprintf("http://music.163.com/discover/artist/cat?id=%d&initial=%d", v, i)
					ctx.AddQueue(
						&request.Request{
							Url:    "http://music.163.com/weapi/v1/resource/comments/R_SO_4_186016/?csrf_token=",
							Rule:   aid["Rule"].(string),
							Method: "post",
							//PostData: "params=EDegnd4Lwi8uPnPUYBgZAuhYKCi25JtXoOxAHPzjQqQuxjqMCVFZltaFAUUOtyhr2z3%2FAL7HU67DF3kIWMz6iqiafonsqsJ0pzE14o1dsDxOA8%2BaYEZ55ICivfnuJeZ0Oxuu7tcCr8tdXDVRfxPQsF1bbzgA8vnV5WoOgk2ITTIRPmc80fEHRdtY%2FY3kVFBS&encSecKey=62f38822dc22ec228c7d59ed389768479b6e1ccbb32cba72dcd2752d787bc2e86bc3bfba0764176b72455759d2a7a6d798ce69ac22b7103c4f669d1af2f8f7d6794f6d20cdf938869ccc00e8b9e2b699affc86a7b65fad0ae1575d88bcf02119fb743193526b6528447b2b527707d19c541fdb2899e9072cc7dc6d7638b0c1e2",
							PostData: "params=tDcRMwjaL%2FQrCTIClHhKiWz7X6TrwJu%2F0EpnccpuCIZD7Bd9msm4DT9Nmm26EOmz08MGR%2F3W1n7Jg%2BQAQKEGUF%2B1IpW0Ze4bj%2FgbAQpydkbGL40Lkgr2iPIYk5rGQplkChUNYuBnsoLXE1gfF24OHvQ%2FoC2W0Mt7ZvqupSmTGl%2FbsYZUbXRW6AhwFn%2FCYbq2&encSecKey=088847ca15a7e0ef6c1ecf57a67d00491c90b8187af97d84aeac20a83d3df386d0c9e5d950943786d7b2ac4856a23d29f6318d41f16bbd7a55624c892eabd73d94546390cb924ad21a8ad9f607d757d4fc9cea6998b67bb69d167ec3ddff628234d83562a36757786533663d2b9bff58457bf3adda25735756d87db2008ac165",
							Header:   map[string][]string{"referer": []string{"http://music.163.com/song?id=29728114"}},
						},
					)
					return nil
				},
				ParseFunc: func(ctx *Context) {
					//query := ctx.GetDom()

					queryString := string(ctx.GetText())
					err := ioutil.WriteFile("test.txt", []byte(queryString), 0644)
					if err != nil {
						log.Fatal(err)
					}

					//query.Find(".nm.nm-icn.f-thide.s-fc0").Each(func(i int, s *goquery.Selection) {
					//	artistName := s.Text()
					//	artistUrl, _ := s.Attr("href")
					//
					//	ctx.Output(map[int]interface{}{
					//		0: artistName,
					//		1: artistUrl,
					//	})
					//})
				},
			},

			"测试": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					ctx.AddQueue(
						&request.Request{
							Url:  "http://music.163.com/album?id=3067018",
							Rule: aid["Rule"].(string),
						},
					)
					return nil
				},
				ParseFunc: func(ctx *Context) {
					//query := ctx.GetDom()

					queryString := string(ctx.GetText())
					err := ioutil.WriteFile("test.txt", []byte(queryString), 0644)
					if err != nil {
						log.Fatal(err)
					}
				},
			},
		},
	},
}
