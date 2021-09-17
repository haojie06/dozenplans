package mailer

import (
	"context"
	"dozenplans/models/dao"
	"dozenplans/models/tables"
	"dozenplans/utils"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

const (
	SMTP_MAIL_HOST     = "smtp.mailgun.org"
	SMTP_MAIL_PORT     = "587"
	SMTP_MAIL_USER     = "plan@rmrf.online"
	SMTP_MAIL_NICKNAME = "DozenPlans"
)

var (
	SMTP_MAIL_PWD       string
	MAILGUN_KEY         string
	MAILGUN_PRIVATE_KEY string
	MAIL_SENDER         = "plan@rmrf.online"
	MAILGUN_DOMAIN      = "rmrf.online"
)

type Mail struct {
	Sender    string
	Subject   string
	Body      string
	Recipient string
	NickName  string
}

// 考虑到发信人一般不变
func NewMail() Mail {
	return Mail{
		Sender:   MAIL_SENDER,
		NickName: SMTP_MAIL_NICKNAME,
	}
}

// 测试使用smtp发信 建议另外开goroutine
func SendEmailViaSmtp(mail Mail) (err error) {
	SMTP_MAIL_PWD = os.Getenv("SMTP_MAIL_PWD")
	auth := smtp.PlainAuth("", SMTP_MAIL_USER, SMTP_MAIL_PWD, SMTP_MAIL_HOST)
	contentType := "Content-Type: text/html; charset=UTF-8"
	// 发送的信息
	s := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s", mail.Recipient, mail.NickName, mail.Sender, mail.Subject, contentType, mail.Body)
	msg := []byte(s)
	addr := fmt.Sprintf("%s:%s", SMTP_MAIL_HOST, SMTP_MAIL_PORT)
	err = smtp.SendMail(addr, auth, SMTP_MAIL_USER, []string{mail.Recipient}, msg)
	if err != nil {
		utils.LogErr(err, "SMTP SEND MAIL")
	} else {
		log.Println("邮件已发送")
	}
	return
}

// 直接使用api发信 发图片需要base64处理
func SendEmailViaApi(mail Mail) {
	MAILGUN_PRIVATE_KEY = os.Getenv("MAILGUN_PRIVATE_KEY")
	mg := mailgun.NewMailgun(MAILGUN_DOMAIN, MAILGUN_PRIVATE_KEY)
	message := mg.NewMessage(mail.Sender, mail.Subject, "", mail.Recipient)
	message.SetHtml(mail.Body)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)
	utils.LogErr(err, "Send Mail via API "+MAILGUN_PRIVATE_KEY)
	log.Printf("Send Mail ID: %s Resp: %s\n", id, resp)
}

// 注册成功通知
func SendRegisterNotification(user tables.User) {
	mail := NewMail()
	mail.Recipient = user.Email
	mail.Subject = "欢迎使用DozenPlans"
	mail.Body = fmt.Sprintf(`
	<h1>您已注册DozenPlans</h1>
	<p>%s 欢迎使用</p>
	`, user.UserName)
	go SendEmailViaSmtp(mail)
}

// 任务超时失败通知
func SendTaskFailedNotification(task tables.Task) {
	mail := NewMail()
	user, _ := dao.GetUserById(task.UserId)
	mail.Recipient = user.Email
	mail.Subject = "任务超时提醒"
	mail.Body = fmt.Sprintf(taskFailedEmailTemplate, task.TaskName, task.Content)
	go SendEmailViaSmtp(mail)
}

// 发送邮件并更新状态
func SendTaskNotificationEmail(task tables.Task) {
	// 修改task的状态，避免重复发送
	task.MailSend = true
	dao.UpdateTask(&task)
	// 根据task获取用户
	user, _ := dao.GetUserById(task.UserId)
	mail := NewMail()
	mail.Recipient = user.Email
	// 两种邮件使用不同的模板
	if task.NotifyMode == "timing" {
		mail.Subject = "DozenPlans 定时提醒"
		mail.Body = fmt.Sprintf(taskNotificationEmailTemplate, "定时提醒", task.TaskName)
	} else if task.NotifyMode == "interval" {
		mail.Subject = "DozenPlans 间隔提醒"
		mail.Body = fmt.Sprintf(taskNotificationEmailTemplate, "间隔提醒", task.TaskName)
	}
	// 发送邮件
	err := SendEmailViaSmtp(mail)
	if err != nil {
		log.Println("邮件发送失败", err.Error())
	} else {
		// 发送成功 更新数据库
		if task.NotifyMode == "timing" {
			task.Notified = true
			dao.UpdateTask(&task)
		} else {
			// 选取的间隔暂时是秒
			task.NotifyTime = time.Now().Add(time.Duration(task.NotifyInterval * int64(time.Second)))
			dao.UpdateTask(&task)
		}
		log.Println("已发送通知邮件 TASK:", task.TaskName)
	}
	// 解除限制
	task.MailSend = false
	dao.UpdateTask(&task)
}
