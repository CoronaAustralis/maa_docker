package utils

import (
	"maa-server/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func EmailPush(){
	mail := gomail.NewMessage()
    mail.SetHeader("From", config.Conf.EmailPush.EmailAddress) 
    mail.SetHeader("To", config.Conf.EmailPush.EmailAddress) 
    mail.SetHeader("Subject", "Maa_docker Event Failed Notification")
    mail.SetBody("text/plain", "maa_docker 事件执行失败")

    dialer := gomail.NewDialer(config.Conf.EmailPush.Host, config.Conf.EmailPush.Port, config.Conf.EmailPush.EmailAddress, config.Conf.EmailPush.Token)

    if err := dialer.DialAndSend(mail); err != nil {
        log.Errorln("failed send email, err: ",err)
		return
    }

    log.Infoln("email sent successfully")
}