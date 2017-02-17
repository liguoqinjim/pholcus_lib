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

var QuestionNum = 0

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

var WeiboAskRunMode = "test"

//var WeiboAskRunMode = "run"

var WeiboAskSpider = &Spider{
	Name:         "微博问答问题",
	Description:  "微博问答问题爬虫",
	Pausetime:    2000,
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
							//if k != 48 {
							//	continue
							//}
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
							if WeiboAskRunMode == "test" {
								break
							}
						} else {
							//todo 测试
							if WeiboAskRunMode == "test" {
								break
							}

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
						if WeiboAskRunMode == "test" {
							pageCount = 1
						}

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
							if WeiboAskRunMode == "test" {
								pageCount = 1
							}

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
						//fmt.Printf("查询答主问题价格 %+v\n", v)
						ctx.Aid(map[string]interface{}{"data": v}, "查询答主问题价格")
						//todo 测试
						if WeiboAskRunMode == "test" {
							break
						}
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

					total_count, _ := strconv.Atoi(answererData.Data.Total_count)
					if total_count > 0 {
						//fmt.Printf("%+v\n", answererData.Data.Pager_info)
						ctx.Aid(map[string]interface{}{"data": askData, "answererData": answererData}, "查询问题")
					}
				},
			},

			"查询问题": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					askData := aid["data"].(AskerData)
					answererData := aid["answererData"].(*AnswererData)

					total_count, _ := strconv.Atoi(answererData.Data.Total_count)
					if total_count > 0 {
						for i := 1; i <= answererData.Data.Pager_info.Total_page; i++ {
							//http://e.weibo.com/v1/public/aj/qa/getselleranswer?uid=1979899604&page=2

							uid := strings.Split(askData.Content_url, "=")[1]
							url := fmt.Sprintf("http://e.weibo.com/v1/public/aj/qa/getselleranswer?uid=%s&page=%d", uid, i)
							referer := fmt.Sprintf("http://e.weibo.com/v1/public/center/qauthor?uid=%s", uid)
							ctx.AddQueue(
								&request.Request{
									Url: url,
									Header: http.Header{
										"User-Agent": []string{"Mozilla/5.0 (Linux; Android 4.4.2; Lenovo A3300-T Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36 Weibo (LENOVO-Lenovo A3300-T__weibo__7.0.0__android__android4.4.2)"},
										"Cookie":     []string{ask_cookies2},
										"referer":    []string{referer},
									},
									Temp: aid,
									Rule: "查询问题",
								},
							)
							//todo 测试
							if WeiboAskRunMode == "test" {
								break
							}
						}
					}

					return nil
				},
				ParseFunc: func(ctx *Context) {
					askData := ctx.GetTemp("data", "test").(AskerData)
					answererData := ctx.GetTemp("answererData", "test").(*AnswererData)

					json1, err := simplejson.NewJson([]byte(ctx.GetDom().Text()))
					if err != nil {
						fmt.Println(err)
						return
					}

					json2, err := json1.Get("data").Get("list").Array()
					if err != nil {
						fmt.Println(err)
						return
					}

					//fmt.Println(ctx.Request.GetUrl())
					//fmt.Println("查询问题---", ctx.GetDom().Text())
					//fmt.Println("查询问题-->", len(json2))
					for i := 0; i < len(json2); i++ {
						content_url, _ := json1.Get("data").Get("list").GetIndex(i).Get("content_url").String()

						ctx.Aid(map[string]interface{}{"data": askData, "answererData": answererData, "question_url": content_url}, "查询问题详细")
						//todo 测试
						if WeiboAskRunMode == "test" {
							break
						}
					}
				},
			},

			"查询问题详细": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					question_url := aid["question_url"].(string)

					question_id := strings.Split(question_url, "=")[1]
					//注意：%号转义
					url_temp := "http://api.weibo.cn/2/question/show?networktype=wifi&uicode=10000432&moduleID=700&wb_version=3319&c=android&i=61e6992&s=773ce9dd&ua=LENOVO-Lenovo%%20A3300-T__weibo__7.0.0__android__android4.4.2&wm=2468_1001&aid=01ApIEZ_RFW8QFgeOItuEYX1q0tJxDA9C2a8HBmnmEK9iF5K8.&oid=%s&v_f=2&from=1070095010&gsid=_2A251ofP_DeRxGeBP6FcQ9inFwzuIHXVUDxQurDV6PUJbkdAKLWLNkWqC3NU0yJ_c7bOoukrk9g-NKkmkKg..&lang=zh_CN&skin=default&vuid=6135167987&oldwm=2468_1001&sflag=1&luicode=80000001"
					url := fmt.Sprintf(url_temp, question_id)
					ctx.AddQueue(
						&request.Request{
							Url: url,
							Header: http.Header{
								"X-Log-Uid":  []string{"6135167987"},
								"User-Agent": []string{"Lenovo A3300-T_4.4.2_weibo_7.0.0_android"},
							},
							Temp: aid,
							Rule: "查询问题详细",
						},
					)

					return nil
				},
				ParseFunc: func(ctx *Context) {

					json1, err := simplejson.NewJson([]byte(ctx.GetDom().Text()))
					if err != nil {
						fmt.Println(err)
						return
					}
					aid := ctx.GetTemps()

					askData := ctx.GetTemp("data", "test").(AskerData)
					answererData := ctx.GetTemp("answererData", "test").(*AnswererData)
					question_url := ctx.GetTemp("question_url", "test").(string)
					aid["answer_detail"] = json1

					//ask_content, _ := json1.Get("ask_content").String() //问题内容
					//oid, _ := json1.Get("object_id").String()           //问题id
					//qa_price, _ := json1.Get("qa_price").String()       //问题价格
					//
					//ask_at, _ := json1.Get("ask_at").String()           //问题时间
					//answer_at, _ := json1.Get("answer_at").String()     //回答时间
					//
					//json2 := json1.Get("answerer")
					//answerer_name, _ := json2.Get("name").String()
					//answerer_id, _ := json2.Get("id").Int()
					//followers_count, _ := json2.Get("followers_count").Int()
					//description, _ := json2.Get("description").String()
					//
					//json3 := json1.Get("asker")
					//asker_id, _ := json3.Get("id").Int()
					//asker_name, _ := json3.Get("name").String()
					ctx.Aid(map[string]interface{}{"data": askData, "answererData": answererData, "question_url": question_url, "answer_detail": json1}, "查询问题围观数")
				},
			},

			"查询问题围观数": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					question_url := aid["question_url"].(string)

					question_id := strings.Split(question_url, "=")[1]
					fid := strings.Split(question_id, ":")[1]
					//注意：%号转义
					url_temp := "http://api.weibo.cn/2/question/extend?networktype=wifi&uicode=10000432&moduleID=700&wb_version=3319&c=android&i=61e6992&s=773ce9dd&ua=LENOVO-Lenovo%%20A3300-T__weibo__7.0.0__android__android4.4.2&wm=2468_1001&aid=01ApIEZ_RFW8QFgeOItuEYX1q0tJxDA9C2a8HBmnmEK9iF5K8.&fid=%s&oid=%s&v_f=2&from=1070095010&gsid=_2A251ofP_DeRxGeBP6FcQ9inFwzuIHXVUDxQurDV6PUJbkdAKLWLNkWqC3NU0yJ_c7bOoukrk9g-NKkmkKg..&lang=zh_CN&read=false&skin=default&vuid=6135167987&oldwm=2468_1001&sflag=1&luicode=80000001"
					url := fmt.Sprintf(url_temp, fid, question_id)
					ctx.AddQueue(
						&request.Request{
							Url: url,
							Header: http.Header{
								"X-Log-Uid":  []string{"6135167987"},
								"User-Agent": []string{"Lenovo A3300-T_4.4.2_weibo_7.0.0_android"},
							},
							Temp: aid,
							Rule: "查询问题围观数",
						},
					)

					return nil
				},
				ItemFields: []string{
					"问题内容",
					"问题id",
					"问题价格",
					"问题围观数",
					"问题时间",
					"回答时间",
					"答题者微博名",
					"答题者id",
					"答题者粉丝数",
					"答题者描述",
					"提问者微博名",
					"提问者id",
				},
				ParseFunc: func(ctx *Context) {
					jsonA, err := simplejson.NewJson([]byte(ctx.GetDom().Text()))
					if err != nil {
						fmt.Println(err)
						return
					}
					interact_count, _ := jsonA.Get("interact_user_info").Get("interact_count").String()

					json1 := ctx.GetTemps()["answer_detail"].(*simplejson.Json)
					ask_content, _ := json1.Get("ask_content").String() //问题内容
					oid, _ := json1.Get("object_id").String()           //问题id
					qa_price, _ := json1.Get("qa_price").String()       //问题价格
					ask_at, _ := json1.Get("ask_at").String()           //问题时间
					answer_at, _ := json1.Get("answer_at").String()     //回答时间

					json2 := json1.Get("answerer")
					answerer_name, _ := json2.Get("name").String()
					answerer_id, _ := json2.Get("id").Int()
					followers_count, _ := json2.Get("followers_count").Int()
					description, _ := json2.Get("description").String()

					json3 := json1.Get("asker")
					asker_id, _ := json3.Get("id").Int()
					asker_name, _ := json3.Get("name").String()

					ctx.Output(map[int]interface{}{
						0: ask_content,
						1: oid,
						2: qa_price,
						3: interact_count, 4: ask_at, 5: answer_at,
						6: answerer_name, 7: answerer_id, 8: followers_count, 9: description,
						10: asker_name, 11: asker_id,
					})
				},
			},
		},
	},
}
