package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"parking/utils"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	// 登录页面
	cli := utils.NewHttpCli()
	_, _, err := cli.Get(utils.LoginPageUrl)
	if err != nil {
		fmt.Println("login page request error", err)
	}

	// 登录
	passBase64 := base64.StdEncoding.EncodeToString([]byte(utils.PassWd))
	paramReader := strings.NewReader("uname=" + utils.UsrName + "&password=" + passBase64 + "&fid=-1&t=true&refer=http://i.chaoxing.com")
	body, cookie, err := cli.Post(utils.LoginUrl, paramReader)
	r, err := regexp.Compile(",\"status\":(.*)}")
	if err != nil {
		fmt.Println("name regexp failed", err)
	}
	status := r.FindStringSubmatch(body)
	if status[1] == "false" {
		err = fmt.Errorf("usrName or passWd error")
	}
	if err != nil {
		fmt.Println("login error", err)
	}
	// 保存个人cookie数据
	usrCookie := make(map[string]string, 0)
	for _, v := range cookie {
		usrCookie[v.Name] =  v.Value
	}

	// 获取名字
	body, _, err = cli.Get(utils.AccountManageUrl)
	if err != nil {
		fmt.Println("account manage page request failed", err)
	}
	r, err = regexp.Compile("id=\"messageName\">(.*)</p></div>")
	if err != nil {
		fmt.Println("name regexp failed", err)
	}
	name := r.FindStringSubmatch(body)
	fmt.Println("welcome ", name[1])

	//获取课程列表
	param := strings.NewReader("courseType=1&courseFolderId=0&courseFolderSize=0")
	body, _, err = cli.Post(utils.CourseListUrl, param)
	//提取课程id 班级id
	r, err = regexp.Compile("<li class=\"course clearfix\" courseId=\"(.*)\" clazzId=\"(.*)\" personId=\"")
	if err != nil {
		fmt.Println("class regexp failed", err)
	}
	class := r.FindAllStringSubmatch(body,-1)
	classes := make(map[string]string, 0)
	for _, v := range class {
		classes[v[1]] = v[2]

	}

	// 检查是否有签到活动
	fmt.Println("检测是否存在需要签到的课程")
	var (
		wg sync.WaitGroup
		active []map[string]string
		activeLock sync.Mutex
	)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	for course, classId := range classes {
		wg.Add(1)
		go func(cou, cla string) {
			defer wg.Done()
			err = fmt.Errorf("begin")
			for i := 0; i < 3 && err != nil; i++ {

				body, _, err = cli.Get(utils.ActiveCourseListUrl + "?fid=0&courseId=" + cou +"&classId=" + cla + "&_=" + now)
				info := classInfo{}
				err = json.Unmarshal([]byte(body),&info)
				if err != nil || len(info.Data.ActiveList) < 1 {
					continue
				}
				otherId, _ := strconv.Atoi(info.Data.ActiveList[0].OtherID)
				// 判断是否有效签到活动
				if  otherId >= 0 && otherId <= 5 && info.Data.ActiveList[0].Status == 1 {
					// 活动开始超过一小时则忽略
					if (time.Now().Unix() - info.Data.ActiveList[0].StartTime) / 1000 < 7200 {
						fmt.Println("检测到活动：", info.Data.ActiveList[0].NameOne)
						activeLock.Lock()
						active = append(active, map[string]string {
							"aid": strconv.Itoa(info.Data.ActiveList[0].ID),
							"name": info.Data.ActiveList[0].NameOne,
							"otherId": info.Data.ActiveList[0].OtherID,
						})
						activeLock.Unlock()
					}
				}
			}
		}(course, classId)
	}
	wg.Wait()

	// 将各个活动签到签一遍
	fmt.Println("检测完毕 开始对每个课程两小时内最新的签到活动签到")
	for _, v := range active {
		signType, err := strconv.Atoi(v["otherId"])
		fmt.Println("signType",signType)
		if err != nil {
			fmt.Println(" sign type not a num")
		}
		if signType == 0 {
			objectId := uploadPic(usrCookie["_uid"], cli)
			url := utils.SignUrl + "?activeId=" + v["aid"] + "&uid=" + usrCookie["_uid"] +
				"&clientip=&useragent=&latitude=-1&longitude=-1&appType=15&fid=" + usrCookie["fid"] +
				"&objectId=" + objectId
			bo, _, _ := cli.Get(url)
			if bo != "success" && bo != "您已签到过了"{
				fmt.Println("failed", bo)
			}
		} else if signType == 2 {
			fmt.Println("课堂" + v["name"] + "为二维码签到,请输入enc:")
			var enc string
			fmt.Scanf("%s\n", &enc)
			cli.Get(utils.SignUrl +
				"?enc=" + enc + "&activeId=" + v["aid"] + "&uid=" + usrCookie["_uid"] +
				"&clientip=&useragent=&latitude=-1&longitude=-1&fid=" + usrCookie["fid"] + "&appType=15",
			)
		} else if signType == 3 || signType == 5 {
			b, _, e := cli.Get(utils.SignUrl + "?activeId=" + v["aid"] + "&uid=" + usrCookie["_uid"] +
				"&clientip=&latitude=-1&longitude=-1&appType=15&fid=" + usrCookie["fid"],
			)
			fmt.Println("body:", b,"error:", e,utils.SignUrl + "?activeId=" + v["aid"] + "&uid=" + usrCookie["_uid"] +
				"&clientip=&latitude=-1&longitude=-1&appType=15&fid=" + usrCookie["fid"])
		} else if signType == 4 {
			fmt.Print("课堂" + v["name"] + "为位置签到,请选择一个位置，输入该位置前序号:\n")
			var choice int
			location := [][]string{{"湘潭大学逸夫楼", "27.892341", "112.868043"}, {"计算机中心", "27.886079", "112.868109"}}
			for index, va := range location {
				fmt.Println(index, va)
			}
			fmt.Scanf("%d\n", &choice)
			cli.Get(utils.SignUrl + "?address=" + url.QueryEscape(location[choice][0]) + "&activeId=" + v["aid"] +
				"&uid=" + usrCookie["_uid"] + "&clientip=&latitude=" + location[choice][1] + "&longitude=" +
				location[choice][2] + "&fid=" + usrCookie["fid"] + "&appType=15&ifTiJiao=1",
			)
		}
	}
}

func uploadPic(uid string, cli *utils.HttpCli) string {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", "PIC002413.png")
	if err != nil {
		fmt.Println(err)
	}

	// 从文件读取数据，写入表单
	srcFile, err := os.Open("./conf/test.png")
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		fmt.Println(err)
	}
	_ = writer.WriteField("puid", uid)

	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	token := tokenRes{}
	b, _, _ := cli.Get("https://pan-yz.chaoxing.com/api/token/uservalid")
	json.Unmarshal([]byte(b), &token)
	b, _, _ = cli.PostPic("https://pan-yz.chaoxing.com/upload?_token=" + token.Token + "&puid=" + uid, body)
	res := uploadPicRes{}
	json.Unmarshal([]byte(b), &res)
	fmt.Println(res)
	return res.ObjectID
}

// </li>
//                                   <li class="course clearfix" courseId="216216767" clazzId="35985025" personId="106312467" id="course_216216767_35985025">
//                       <div class="course-cover">
//                                                           <a href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=216216767&clazzid=35985025&vc=1&cpi=106312467&ismooc2=1" target="_blank">
//                                          <img src="https://p.ananas.chaoxing.com/star3/240_130c/5be9e34ed211a22e88040a7f5234fe19.png">
//                                   </a>
//                                                  <ul class="hanlde-list">
//                               <li class="move-to movetobtn">移动到                               </li>
//                                                                    <li class="quit-course quitCourseBtn">退课</li>
//                                                        </ul>
//                                                                                                       </div>
//                       <div class="course-info">
//                           <h3 class="inlineBlock">
//                                                       <a  class="color1" href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=216216767&clazzid=35985025&vc=1&cpi=106312467&ismooc2=1" target="_b
//lank">
//                                         <span class="course-name overHidden2" title="计网劳动课">计网劳动课</span>
//                                   </a>
//                                                   </h3>
//                           <p class="margint10 line2" title=""></p>
//                           <p class="line2" title="文磊">文磊</p>
//                                                     <p class="overHidden1">班级：默认班级</p>
//                       </div>
//                   </li>
//                                   <li class="course clearfix" courseId="214482214" clazzId="36010063" personId="92073247" id="course_214482214_36010063">
//                       <div class="course-cover">
//                                                           <a href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=214482214&clazzid=36010063&vc=1&cpi=92073247&ismooc2=1" target="_blank">
//                                          <img src="https://p.ananas.chaoxing.com/star3/240_130c/3b7e2b39a73c4e93969c903206cb48a8.jpg">
//                                   </a>
//                                                  <ul class="hanlde-list">
//                               <li class="move-to movetobtn">移动到                               </li>
//                                                                    <li class="quit-course quitCourseBtn">退课</li>
//                                                        </ul>
//                                                                                                       </div>
//                       <div class="course-info">
//                           <h3 class="inlineBlock">
//                                                       <a  class="color1" href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=214482214&clazzid=36010063&vc=1&cpi=92073247&ismooc2=1" target="_bl
//ank">
//                                         <span class="course-name overHidden2" title="计算机组成与体系结构">计算机组成与体系结构</span>
//                                   </a>
//                                                   </h3>
//                           <p class="margint10 line2" title="本课程是高等院校工科本科软件工程专业的一门学科基础课（必修）。通过对计算机硬件系统中的各部件工作过程、组成结构、部件之间信息传送的讲述，使学生掌握
//计算机硬件的基本概念、计算机系统的组成、各部件的工作原理和工作过程以及它们之间的信息传送方法，初步具有计算机硬件系统的分析、设计的能力。">本课程是高等院校工科本科软件工程专业的一门学科基础课（必修）。通过对
//计算机硬件系统中的各部件工作过程、组成结构、部件之间信息传送的讲述，使学生掌握计算机硬件的基本概念、计算机系统的组成、各部件的工作原理和工作过程以及它们之间的信息传送方法，初步具有计算机硬件系统的分析、设计
//的能力。</p>
//                           <p class="line2" title="成洁">成洁</p>
//                                                     <p class="overHidden1">班级：2019网络工程2班</p>
//                       </div>
//                   </li>
//                                   <li class="course clearfix" courseId="214812873" clazzId="34570944" personId="92073247" id="course_214812873_34570944">
//                       <div class="course-cover">
//                                                           <a href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=214812873&clazzid=34570944&vc=1&cpi=92073247&ismooc2=1" target="_blank">
//                                          <img src="https://p.ananas.chaoxing.com/star3/240_130c/f1b0790809b71680500e5022471cd111.jpg">
//                                   </a>
//                                                  <ul class="hanlde-list">
//                               <li class="move-to movetobtn">移动到                               </li>
//                                                                    <li class="quit-course quitCourseBtn">退课</li>
//                                                        </ul>
//                                                                                                       </div>
//                       <div class="course-info">
//                           <h3 class="inlineBlock">
//                                                       <a  class="color1" href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=214812873&clazzid=34570944&vc=1&cpi=92073247&ismooc2=1" target="_bl
//ank">
//                                         <span class="course-name overHidden2" title="计算机组成原理与数字电路">计算机组成原理与数字电路</span>
//                                   </a>
//                                                   </h3>
//                           <p class="margint10 line2" title=""></p>
//                           <p class="line2" title="周黎黎">周黎黎</p>
//                                                     <p class="overHidden1">班级：2019网络工程2班</p>
//                       </div>
//                   </li>
//                                   <li class="course clearfix" courseId="215001161" clazzId="33106345" personId="92073247" id="course_215001161_33106345">
//                       <div class="course-cover">
//                                                           <a href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=215001161&clazzid=33106345&vc=1&cpi=92073247&ismooc2=1" target="_blank">
//                                          <img src="https://p.ananas.chaoxing.com/star3/240_130c/8a0d7e35ecbb20ec4b19f185eed0bc2b.jpg">
//                                   </a>
//                                                  <ul class="hanlde-list">
//                               <li class="move-to movetobtn">移动到                               </li>
//                                                                    <li class="quit-course quitCourseBtn">退课</li>
//                                                        </ul>
//                                                                                                       </div>
//                       <div class="course-info">
//                           <h3 class="inlineBlock">
//                                                       <a  class="color1" href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=215001161&clazzid=33106345&vc=1&cpi=92073247&ismooc2=1" target="_bl
//ank">
//                                         <span class="course-name overHidden2" title="面向对象程序设计(Java)">面向对象程序设计(Java)</span>
//                                   </a>
//                                                   </h3>
//                           <p class="margint10 line2" title=""></p>
//                           <p class="line2" title="李枚毅">李枚毅</p>
//                                                     <p class="overHidden1">班级：默认班级</p>
//                       </div>
//                   </li>
//                                   <li class="course clearfix" courseId="214829096" clazzId="32607627" personId="106312467" id="course_214829096_32607627">
//                       <div class="course-cover">
//                                                           <a href="https://mooc1-1.chaoxing.com/visit/stucoursemiddle?courseid=214829096&clazzid=32607627&vc=1&cpi=106312467&ismooc2=1" target="_blank">
//                                          <img src="https://p.ananas.chaoxing.com/star3/240_130c/3a7c4e2929905646967168696f78ff01.png">
//                                   </a>
//                                                  <ul class="hanlde-list">
//                               <li class="move-to movetobtn">移动到                               </li>
//                                                                    <li class="quit-course quitCourseBtn">退课</li>
//                                                        </ul>
//