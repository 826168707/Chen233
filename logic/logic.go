package logic

import (
	"LedgerProject/dao"
	models "LedgerProject/models"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

//不存在返回true  存在返回false
func EmailCheck(email *string) bool {
	//检查数据库中是否存在相同email
	err, _ := models.FindUserByEmail(email)
	if err != nil { //没有找到相同数据
		return true
	} else {
		return false
	}
}

//存在用户返回true  不存在返回false
func UserCheck(email, password *string) bool {
	err, _ := models.FindUserByEmailAndPassword(email, password)
	if err != nil { //没找到相同数据
		return false
	} else {
		return true
	}

}

func LogicGetHome(email string) (err error, username, days string, money, usefulMoney int) {
	var user models.User
	err, user = models.FindUserByEmail(&email)

	username = user.Username
	money = user.Money

	//判断deadline是否为nil----是:days(您还没有设置日期哦!)---否:days为deadline减now的天数
	if user.Deadline == "nil" {
		days = "您还没有设置日期哦!"
		usefulMoney = 0
	} else {
		//检测日期是否过期
		if !IsExpired(user.Deadline) { //过期的情况
			days = "设置的日期已经过期啦,请重新设置本月计划吧!"
			usefulMoney = money
		} else { //未过期
			//计算天数
			days_Int := CalculateDays(user.Deadline)
			days = strconv.Itoa(days_Int)
			//计算可用余额
			usefulMoney = CalculateUsefulMoney(days_Int, money, user.DailyExpenses)

		}
	}
	return
}

//检测是否过期  true表未过期
func IsExpired(deadlline string) bool {
	t2 := StringToTime(deadlline)
	t1 := time.Now()
	return t1.Before(t2)
}

//计算可用余额  可用余额=余额-天数*每日固定支出
func CalculateUsefulMoney(days, money, dailyExpenses int) (usefulMoney int) {
	usefulMoney = money - days*dailyExpenses
	return
}

func CalculateDays(deadline string) (days int) {
	t1 := time.Now()
	t2 := StringToTime(deadline)

	days = int(t2.Sub(t1).Hours() / 24)

	return
}

func VisualCalculateDays(now, setDate string) (days int) {
	t1 := StringToTime(setDate)
	t2 := StringToTime(now)

	days = int(t2.Sub(t1).Hours() / 24)

	return
}

//从session获取email并进行类型修正
func GetEmailFromSession(c *gin.Context) (email string) {
	email = sessions.Default(c).Get("loginuser").(string)
	fmt.Printf("session == %v\n", email)
	return
}

//日期变为时间戳
func DataToTimeStr(deadline *string) (timeSta int64, err error) {
	timeLayout := "2006-01-02"
	loc, err := time.LoadLocation("Asia/Shanghai")
	theTime, err := time.ParseInLocation(timeLayout, *deadline, loc)
	timeSta = theTime.Unix()
	return
}

//将string变为time类型
func StringToTime(str string) (theTime time.Time) {
	timeLayout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Shanghai")
	theTime, _ = time.ParseInLocation(timeLayout, str, loc)
	return
}

//更新对应用户的金额,设置日期,截止日期,日常固定支出
func UpdateCount(email, deadline string, money, dailyExpenses int) (err error) {
	now := time.Now().Format("2006-01-02")
	//更新
	err = models.UpdateMoneyAndDeadline(email, money, dailyExpenses, deadline, now)
	return
}

//添加历史记录
func AddHistory(email string, kind int, money int, comment string) (err error) {
	//获取现在日期
	now := time.Now().Format("2006-01-02")
	err = models.AddOneHistory(email, kind, money, comment, now)
	return
}

//支出提示 根据用户的可用余额与本次花费的关系进行人性化提醒
func CostTip(email string, cost int) (err error, remainMoney int) {

	//返回usefulMoney - cost 值, 前端加上提醒语句
	err, _, _, _, usefulMoney := LogicGetHome(email)
	remainMoney = usefulMoney - cost
	return
}

//根据email获取所有支出记录
func GetCostHistory(email string) (error, []models.History) {
	err, histories := models.FindCostHistoriesByEmail(email)
	return err, histories
}

//根据email获取所有收入记录
func GetIncomeHistory(email string) (error, []models.History) {
	err, histories := models.FindIncomeHistoriesByEmail(email)
	return err, histories
}

//可视数据  算出每种kind的花费之和
func VisualData(email string) (error, map[int]int) {

	err, histories := models.FindCostHistoriesByEmail(email)
	err, user := models.FindUserByEmail(&email)

	sum := map[int]int{}
	for _, v := range histories {
		sum[v.Kind] += v.Money
	}
	//计算到今天的日常花费  days * dailyExpenses

	sum[11] = VisualCalculateDays(time.Now().Format("2006-01-02"), user.SetDate) * user.DailyExpenses

	return err, sum
}

//推荐模块  根据用户的可用余额进行推荐
func GetRecommend(email string) (error, []models.Commodity) {

	//获取用户的可用余额
	err, user := models.FindUserByEmail(&email)
	usefulMoney := CalculateUsefulMoney(CalculateDays(user.Deadline), user.Money, user.DailyExpenses)

	//usefulMoney > 1000 推荐吃喝玩乐

	var commodities []models.Commodity

	if usefulMoney > 1000 {
		err, commodities = models.FindKind1()
	} else if usefulMoney > 500 {
		err, commodities = models.FindKind2() //推荐高性价比商品
	} else if usefulMoney > 100 {
		err, commodities = models.FindKind3() //推荐书
	} else {
		err, commodities = models.FindKind4() //推荐兼职
	}

	return err, commodities
}

//发送验证码
func SendEmail(email string) bool {
	//判断是否填写了email
	if email == "" {
		return false
	}
	//获取随机数
	rand.Seed(time.Now().Unix())
	num := rand.Intn(10000)
	//随机数存入redis,定时5分钟
	if err := dao.AddCaptcha(num); err != nil {
		fmt.Printf("AddCaptcha failed, err:%v\n", err)
		return false
	}

	auth := smtp.PlainAuth("", "826168707@qq.com", "uwdwkbupldkcbeda", "smtp.qq.com")
	to := []string{email}
	nickname := "流云规划"
	user := "826168707@qq.com"
	subject := "流云规划--验证码"
	contentType := "Content-Type: text/plain; charset=UTF-8"
	body := "验证码是\t" + strconv.Itoa(num)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil {
		fmt.Printf("send mail error: %v", err)
		return false
	}
	return true
}

//验证码验证
func CaptchaCheck(captcha string) bool {
	num, err := dao.GetCaptcha()
	if err != nil {
		fmt.Printf("GetCaptcha failed ,err%v\n", err)
		return false
	}

	if captcha == num { //相同
		return true
	} else {
		return false
	}

}
