package notificationservice

import (
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

const (
	FROM   = "noreply@cashfish.gmail.com"
	SignUp = "newusercreation"
)

func Notifier(msg []byte, topic string) {
	switch topic {
	case SignUp:
		if err := signUpNotifier(msg); err != nil {
			log.Println(err)
		}
	}
}

func signUpNotifier(msg []byte) error {
	var model ProducerModel
	if err := json.Unmarshal(msg, &model); err != nil {
		return err
	}

	m := gomail.NewMessage()
	headers := map[string][]string{
		"From":    {FROM},
		"To":      {model.Email},
		"Subject": {"New Account Creation"},
	}
	m.SetHeaders(headers)
	body := fmt.Sprintf(`Hi %s,

	Congratulations! 
	It's time to join the new era of digital payment and be part of the new wave of cashless transactions.

	Hope you will find our service useful for your needs.`, model.Username)
	m.SetBody("text/plain", body)

	newDial := gomail.NewDialer("mailhog", 1025, "", "")
	return newDial.DialAndSend(m)
}
