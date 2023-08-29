package util

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/google/uuid"
)

func SendRealCreateUserMail(to string, token string) error {
	from := os.Getenv("SMTP_USERADDR")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	apiAddr := os.Getenv("PUBLIC_API_ADDR")

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Message-ID: " + "<" + uuid.New().String() + "@" + host + ">\r\n" +
		"Subject: PLAccounting - メールアドレス確認\r\n\r\n" +
		"本登録を完了するには、下記URLにてメールアドレス確認を行なってください。\n" +
		apiAddr + "/createUser?token=" + token

	fmt.Println("Will send message to ", to)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", host+":"+port, tlsConfig)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	if err = c.Auth(smtp.PlainAuth("", user, pass, host)); err != nil {
		fmt.Println("Auth Error: ", err)
		return err
	}

	if err = c.Mail(from); err != nil {
		fmt.Println("Mail Error: ", err)
		return err
	}

	if err = c.Rcpt(to); err != nil {
		fmt.Println("Rcpt Error: ", err)
		return err
	}

	w, err := c.Data()
	if err != nil {
		fmt.Println("Data Error: ", err)
		return err
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		fmt.Println("Write Error: ", err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Println("Close Error: ", err)
		return err
	}

	fmt.Println("Sent message to ", to)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}
