package pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// . "github.com/henrylee2cn/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	//"regexp"
	"strconv"
	"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
	//"fmt"
	"encoding/json"
	"fmt"
	"net/http"
	//json
	"github.com/bitly/go-simplejson"
)

//答主的标签种类
var FieldTypes = map[int]string{
	8:  "健康医疗",
	48: "科学科普",
	56: "财经",
	37: "段子手",
	55: "读书作家",
	36: "星座命理",
	25: "房产家装",
	11: "军事",
	31: "教育",
	18: "娱乐明星",
	-1: "other", //key是负数的时候，生成链接的参数用value
}

type AskData struct {
	Code string
	Msg  string
	Data AskCoreData
}
type AskCoreData struct {
	List       []AskerData
	Total_page int
	Cur_page   int
}
type AskerData struct {
	Avat_url     string
	Identity     int
	Profile_url  string
	Nickname     string
	Intro        string
	Ask_url      string
	Content_url  string
	Question_num int
	Look_num     int
}

func GetAskerJson(content string) (*AskData, error) {
	var ask AskData
	err := json.Unmarshal([]byte(content), &ask)

	return &ask, err
}

type AnswererData struct {
	Code string
	Msg  string
	Data AnswererCoreData
}
type AnswererCoreData struct {
	//List       []AskerData
	Author_info AuthorInfo
	Ask_enable  int
	List        []QuestionData
	Total_count string
	Pager_info  PagerInfo
}
type QuestionData struct {
	Time           string
	Vtype          int
	Avatar         string
	Profile_url    string
	Content_url    string
	Intro          string
	Nickname       string
	Onlooker_count string
	Look_status    string
	Ask_price      string
	Look_price     int
}
type AuthorInfo struct {
	Ask_url     string
	Nickname    string
	Avatar      string
	Profile_url string
	Label       string
	Vtype       int
	Price       string
	Intro       string
}
type PagerInfo struct {
	Total_page int
	Curpage    int
	Pagesize   int
}

func GetAnswererJson(content string) (*AnswererData, error) {
	var ask AnswererData
	err := json.Unmarshal([]byte(content), &ask)

	return &ask, err
}

//博主信息
type WeiboUserInfo struct {
	FollowNum      int
	FriendNum      int
	Description    string
	VerifiedReason string
}

func init() {
	WeiboAskSpider.Register()
}

//var ask_cookies1 = "_s_tentry=-; Apache=6107357982546.091.1486706299965; SINAGLOBAL=6107357982546.091.1486706299965; ULV=1486706299975:1:1:1:6107357982546.091.1486706299965:; SUB=_2A251mSV3DeRhGeNI4lsU9CrEyDiIHXVWtVU_rDV8PUJbitAKLWatkWtZewsAa7ojD2tv3UZtpQVbYffaDw..; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WWpuEKiUVVcFsGgwUyqlHLY5NHD95QfSo.4SKBX1heXWs4DqcjMi--NiK.Xi-2Ri--ciKnRi-zNSKq41K-XShn0S5tt; SCF=AnIMOiTPi591m41NyWBnF6ieIv41bLg3CyEJ8ZDuGzWFIoHkUeD32kZDpZvjSR_YuQ..; SUHB=0KABWDSJpS2qKc"
var ask_cookies2 = "SUB=_2A251mR_mDeRhGeBP6FcQ9inFwzuIHXVWtUuurDV8PUJbitAKLWrAkWtgXLKVDjW1uPm5b-j_sk-g9JE3aQ..; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9Wh.NW5WSuYsKu1Xkoev7WS45NHD95QceKefeKqN1KnNWs4DqcjMi--NiK.Xi-2Ri--ciKnRi-zNSo20SK2cS0.RS7tt; SCF=AnIMOiTPi591m41NyWBnF6hVAZTRxU4VG6Xc2AhuNBIH_CPB6ScH4PgQ1PxsGffSGQ..; SUHB=0mxU7Nhg2YNZng; _s_tentry=-; Apache=1935663854237.6458.1486712796058; SINAGLOBAL=1935663854237.6458.1486712796058; ULV=1486712796112:1:1:1:1935663854237.6458.1486712796058:"

var WeiboAskSpider = &Spider{
	Name:        "微博问答",
	Description: "微博问答爬虫",
	//Pausetime:    12000,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			//Aid调用Rule中的AidFunc
			ctx.Aid(map[string]interface{}{"Rule": "判断页数"}, "判断页数")
		},

		Trunk: map[string]*Rule{
			//判断标签内有多少页
			"判断页数": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for k, v := range FieldTypes {
						if k > 0 {
							//todo 测试
							//continue
							var tempData = map[string]interface{}{"fieldType": k}
							ctx.AddQueue(
								&request.Request{
									//http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=48&page=1&pagesize=10
									//Url: "http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=48&page=1&pagesize=10",
									Url:  "http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=" + strconv.Itoa(k) + "&page=1&pagesize=10",
									Rule: aid["Rule"].(string),
									Temp: tempData,
									Header: http.Header{
										"Cookie": []string{ask_cookies2},
									},
								},
							)

							//todo 测试
							//break
						} else {
							//todo 测试
							//break
							var tempData = map[string]interface{}{"fieldType": k}
							ctx.AddQueue(
								&request.Request{
									//http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=48&page=1&pagesize=10
									//Url: "http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=48&page=1&pagesize=10",
									Url:  "http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=" + v + "&page=1&pagesize=10",
									Rule: aid["Rule"].(string),
									Temp: tempData,
									Header: http.Header{
										"Cookie": []string{ask_cookies2},
									},
								},
							)
						}
					}

					return nil
				},
				ParseFunc: func(ctx *Context) {
					fieldType := ctx.GetTemp("fieldType", "test").(int)
					askData, err := GetAskerJson(ctx.GetDom().Text())
					if err != nil {
						fmt.Println("err=", err)
					}
					totalPage := askData.Data.Total_page

					ctx.Aid(map[string]interface{}{"totalPage": totalPage, "fieldType": fieldType}, "按类按页查询答主")
				},
			},

			"按类按页查询答主": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=8&page=1&pagesize=10
					pageCount := aid["totalPage"].(int)
					fieldType := aid["fieldType"].(int)

					if fieldType > 0 {
						//todo 测试用 pageCount强制等于1
						//pageCount = 1
						for i := 1; i <= pageCount; i++ {
							//注：这里用两个%d，不包含fieldType为负数的情况
							url := fmt.Sprintf("http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=%d&page=%d&pagesize=10", fieldType, i)
							ctx.AddQueue(
								&request.Request{
									Url: url,
									Header: http.Header{
										"Cookie": []string{ask_cookies2},
									},
									Rule: "按类按页查询答主",
								},
							)
						}
					} else {
						for i := 1; i <= pageCount; i++ {
							//todo 测试用 pageCount强制等于1
							//pageCount = 1
							url := fmt.Sprintf("http://e.weibo.com/v1/public/h5/aj/qa/getfamousanswer?fieldtype=%s&page=%d&pagesize=10", FieldTypes[fieldType], i)
							ctx.AddQueue(
								&request.Request{
									Url: url,
									Header: http.Header{
										"Cookie": []string{ask_cookies2},
									},
									Rule: "按类按页查询答主",
								},
							)
						}
					}

					return nil
				},
				ParseFunc: func(ctx *Context) {
					askData, err := GetAskerJson(ctx.GetDom().Text())
					if err != nil {
						fmt.Println("err=", err)
					}

					for _, v := range askData.Data.List {
						ctx.Aid(map[string]interface{}{"data": v}, "查询答主问题价格")
						//todo 测试
						//break
					}
					//ctx.Aid(map[string]interface{}{"totalPage": totalPage, "fieldType": fieldType}, "按类按页查询答主")
				},
			},

			"查询答主问题价格": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					askData := aid["data"].(AskerData)

					//http://e.weibo.com/v1/public/h5/aj/qa/getauthor?uid=2146965345
					uid := strings.Split(askData.Content_url, "=")[1]
					url := fmt.Sprintf("http://e.weibo.com/v1/public/h5/aj/qa/getauthor?uid=%s", uid)
					ctx.AddQueue(
						&request.Request{
							Url: url,
							Header: http.Header{
								"Cookie":  []string{ask_cookies2},
								"referer": []string{askData.Content_url},
							},
							Temp: aid,
							Rule: "查询答主问题价格",
						},
					)

					return nil
				},
				ItemFields: []string{
					"微博名",
					"标签",
					"被围观次数",
					"回答问题次数",
					"提问价格",
				},
				ParseFunc: func(ctx *Context) {
					tmpData := ctx.GetTemp("data", "test")
					var askData AskerData
					if tmpData != "test" {
						askData = tmpData.(AskerData)
					} else {
						return
					}

					answererData, err := GetAnswererJson(ctx.GetDom().Text())
					if err != nil {
						fmt.Println("err=", err)
						return
					}

					ctx.Output(map[int]interface{}{
						0: askData.Nickname,
						1: answererData.Data.Author_info.Label,
						2: askData.Look_num,
						3: answererData.Data.Total_count,
						4: answererData.Data.Author_info.Price,
					})

					//旧
					//ctx.Aid(map[string]interface{}{"data": askData, "answererData": answererData}, "查询答主问题价格(原先存数据的)")

					//查询博主粉丝数量 暂时不用
					//ctx.Aid(map[string]interface{}{"data": askData, "answererData": answererData}, "查询博主粉丝量")
				},
			},

			"查询答主问题价格(原先存数据的)": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					askData := aid["data"].(AskerData)

					//http://e.weibo.com/v1/public/h5/aj/qa/getauthor?uid=2146965345
					uid := strings.Split(askData.Content_url, "=")[1]
					url := fmt.Sprintf("http://e.weibo.com/v1/public/h5/aj/qa/getauthor?uid=%s", uid)
					ctx.AddQueue(
						&request.Request{
							Url: url,
							Header: http.Header{
								"Cookie":  []string{ask_cookies2},
								"referer": []string{askData.Content_url},
							},
							Temp: aid,
							Rule: "查询答主问题价格(原先存数据的)",
						},
					)

					return nil
				},
				ItemFields: []string{
					"微博名",
					"标签",
					"被围观次数",
					"回答问题次数",
					"提问价格",
				},
				ParseFunc: func(ctx *Context) {
					tmpData := ctx.GetTemp("data", "test")
					var askData AskerData
					if tmpData != "test" {
						askData = tmpData.(AskerData)
					} else {
						return
					}

					answererData, err := GetAnswererJson(ctx.GetDom().Text())
					if err != nil {
						fmt.Println("err=", err)
						return
					}

					ctx.Output(map[int]interface{}{
						0: askData.Nickname,
						1: answererData.Data.Author_info.Label,
						2: askData.Look_num,
						3: answererData.Data.Total_count,
						4: answererData.Data.Author_info.Price,
					})
				},
			},

			"查询博主粉丝量": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					askData := aid["data"].(AskerData)

					uid := strings.Split(askData.Content_url, "=")[1]
					//注 这个链接里面有个from，出错的时候可以改这个地方  %%转义
					url := fmt.Sprintf("http://api.weibo.cn/2/profile?networktype=wifi&uicode=10000198&moduleID=708&user_domain=%s&wb_version=3319&c=android&ua=LENOVO-Lenovo%%20A3300-T__weibo__7.0.0__android__android4.4.2&wm=2468_1001&uid=%s&v_f=2&v_p=43&from=1070095010&lang=zh_CN&skin=default&oldwm=2468_1001&sflag=1&cover_width=720&luicode=80000001", uid, uid)
					ctx.AddQueue(
						&request.Request{
							Url:  url,
							Temp: aid,
							Rule: "查询博主粉丝量",
						},
					)

					return nil
				},
				ParseFunc: func(ctx *Context) {
					tmpData := ctx.GetTemp("data", "test")
					var askData AskerData
					if tmpData != "test" {
						askData = tmpData.(AskerData)
					} else {
						return
					}

					tmpData2 := ctx.GetTemp("answererData", "test2")
					var answererData *AnswererData
					if tmpData2 != "test2" {
						answererData = tmpData2.(*AnswererData)
					} else {
						return
					}

					json1, err := simplejson.NewJson([]byte(ctx.GetDom().Text()))
					if err != nil {
						fmt.Println(err)
						return
					}
					follow_num, _ := json1.Get("userInfo").Get("followers_count").Int()
					friend_num, _ := json1.Get("userInfo").Get("friends_count").Int()
					desp, _ := json1.Get("userInfo").Get("description").String()
					reason, _ := json1.Get("userInfo").Get("verified_reason").String()

					weiboUserinfo := &WeiboUserInfo{FollowNum: follow_num, FriendNum: friend_num, Description: desp, VerifiedReason: reason}
					//fmt.Printf("%+v\n", weiboUserinfo)
					//fmt.Printf("%+v\n", askData)
					//fmt.Printf("%+v\n", answererData)

					ctx.Aid(map[string]interface{}{"data": askData, "answererData": answererData, "weiboUserInfo": weiboUserinfo}, "查询博主最近的微博的评论数和点赞数")
				},
			},

			"查询博主最近的微博的评论数和点赞数": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					askData := aid["data"].(AskerData)

					uid := strings.Split(askData.Content_url, "=")[1]
					//注 这个链接里面有个from，出错的时候可以改这个地方
					url := fmt.Sprintf("http://api.weibo.cn/2/cardlist?networktype=wifi&uicode=10000198&moduleID=708&wb_version=3319&c=android&i=61e6992&s=773ce9dd&ua=LENOVO-Lenovo%20A3300-T__weibo__7.0.0__android__android4.4.2&wm=2468_1001&aid=01ApIEZ_RFW8QFgeOItuEYX1q0tJxDA9C2a8HBmnmEK9iF5K8.&fid=1076031676582524_-_WEIBO_SECOND_PROFILE_WEIBO_ORI&uid=%s&v_f=2&v_p=43&from=1070095010&gsid=_2A251mR_mDeRxGeBP6FcQ9inFwzuIHXVUDxQurDV6PUJbkdAKLWLNkWqGgIgE5Wp6TnD5OhJDNlfHzIu_Gw..&imsi=460073497246957&lang=zh_CN&page=1&skin=default&count=20&oldwm=2468_1001&sflag=1&containerid=1076031676582524_-_WEIBO_SECOND_PROFILE_WEIBO_ORI&luicode=80000001&need_head_cards=0", uid)
					ctx.AddQueue(
						&request.Request{
							Url:  url,
							Temp: aid,
							Rule: "查询博主最近的微博的评论数和点赞数",
						},
					)

					return nil
				},
				ItemFields: []string{
					"微博名",
					"标签",
					"被围观次数",
					"回答问题次数",
					"提问价格",
					"粉丝数",
					"最近微博平均回复数",
					"最近微博平均点赞数",
					"最近微博最大回复数",
					"最近微博最大点赞数",
					"最近微博数量",
				},
				ParseFunc: func(ctx *Context) {
					tmpData := ctx.GetTemp("data", "test")
					var askData AskerData
					if tmpData != "test" {
						askData = tmpData.(AskerData)
					} else {
						return
					}

					tmpData2 := ctx.GetTemp("answererData", "test2")
					var answererData *AnswererData
					if tmpData2 != "test2" {
						answererData = tmpData2.(*AnswererData)
					} else {
						return
					}

					tmpData3 := ctx.GetTemp("weiboUserInfo", "test3")
					var weiboUserInfo *WeiboUserInfo
					if tmpData3 != "test3" {
						weiboUserInfo = tmpData3.(*WeiboUserInfo)
					} else {
						return
					}

					json1, err := simplejson.NewJson([]byte(ctx.GetDom().Text()))
					if err != nil {
						fmt.Println(err)
						return
					}
					sum := 0
					maxAtti := 0
					attis := 0
					maxComment := 0
					comments := 0
					json2, err := json1.Get("cards").Array()
					for i := 0; i < len(json2); i++ {
						atti, err2 := json1.Get("cards").GetIndex(i).Get("mblog").Get("attitudes_count").Int()
						if err2 != nil {
							fmt.Println("err2=", err2)
							continue
						}
						comment, err3 := json1.Get("cards").GetIndex(i).Get("mblog").Get("comments_count").Int()
						if err3 != nil {
							fmt.Println("err3=", err3)
							continue
						}
						sum++
						if comment > maxComment {
							maxComment = comment
						}
						if atti > maxAtti {
							maxAtti = atti
						}
						attis += atti
						comments += comment
					}

					ctx.Output(map[int]interface{}{
						0:  askData.Nickname,
						1:  answererData.Data.Author_info.Label,
						2:  askData.Look_num,
						3:  answererData.Data.Total_count,
						4:  answererData.Data.Author_info.Price,
						5:  weiboUserInfo.FollowNum,
						6:  float64(comments / sum),
						7:  float64(attis / sum),
						8:  maxComment,
						9:  maxAtti,
						10: sum,
					})
				},
			},
		},
	},
}
