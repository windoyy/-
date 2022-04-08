package utils

import (
	"fmt"
	"os"
	//使用了beego框架的配置文件读取模块
	"github.com/beego/beego/v2/core/config"
)

var (
	LoginPageUrl        string //登录页面url
	LoginUrl            string
	AccountManageUrl    string
	CourseListUrl       string
	ActiveCourseListUrl string
	SignUrl           string

	UsrName  string //项目名称
	PassWd  string //服务器ip地址
)

func InitConfig() {
	//从配置文件读取配置信息
	if _, err := os.Stat("./conf/app.conf"); os.IsNotExist(err) {
		fmt.Println("config file not exit, creat a blank config file.")
		if _, err = os.Stat("./conf"); os.IsNotExist(err) {
			if os.Mkdir("./conf", 0755) != nil {
				fmt.Println("mkdir conf failed.")
				return
			}
		}
		if file, err := os.Create("./conf/app.conf"); err != nil{
			fmt.Println("creat app.conf failed", err)
			return
		} else {
			file.Close()
		}
	}
	appConf, err := config.NewConfig("ini", "./conf/app.conf")
	if err != nil {
		fmt.Println(err)
		return
	}

	LoginPageUrl        = "http://passport2.chaoxing.com/mlogin?fid=&newversion=true&refer=http%3A%2F%2Fi.chaoxing.com"
	LoginUrl            = "http://passport2.chaoxing.com/fanyalogin"
	AccountManageUrl    = "http://passport2.chaoxing.com/mooc/accountManage"
	CourseListUrl       = "http://mooc1-1.chaoxing.com/visit/courselistdata"
	ActiveCourseListUrl = "https://mobilelearn.chaoxing.com/v2/apis/active/student/activelist"
	SignUrl           = "https://mobilelearn.chaoxing.com/pptSign/stuSignajax"

	UsrName, err = appConf.String("usrName")
	if err != nil || len(UsrName) == 0 {
		fmt.Println("请输入手机号")
		fmt.Scanf("%s\n", &UsrName)
	} else {
		fmt.Println("从配置文件读取到手机号")
	}

	PassWd, err  = appConf.String("passWd")
	if err != nil || len(UsrName) == 0 {
		fmt.Println("请输入手密码")
		fmt.Scanf("%s\n", &UsrName)
	} else {
		fmt.Println("从配置文件读取到密码")
	}

}

func init() {
	InitConfig()
}