package job

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	ioconfig "ims-server/pkg/config"
	"net/smtp"
)

func SendEmail(recipient string, vcode string, url string) error {
	emailConf := ioconfig.GetEmailConf()
	html := fmt.Sprintf(`<div>
        <div>
            尊敬的用户，您好！
        </div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>你本次的验证码为 %s,为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。</p>
			<p><a href="%s">点此返回上一界面</a></p>
        </div>
        <div>
            <p>此邮箱为系统邮箱，请勿回复。</p>
        </div>
    </div>`, vcode, url)

	mail := email.NewEmail()
	mail.From = fmt.Sprintf("io-club ims <%s>", emailConf.MailUserName)
	mail.To = []string{recipient}
	mail.Subject = "io-club ims 验证码"
	mail.HTML = []byte(html)
	err := mail.SendWithTLS(emailConf.Addr, smtp.PlainAuth("", emailConf.MailUserName, emailConf.MailPassword, emailConf.Host),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		return err
	}
	return nil
}
