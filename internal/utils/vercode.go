import (
	"math/rand"
	"time"
	"gopkg.in/gomail.v2"

	"user-management-system/internal/config"
)

func GenerateVerCode(len int) string {
	rand.Seed(time.Now().UnixNano())
	var chars = []rune("0123456789")
	b := make([]rune, len)
	for i := range(b) {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func SendVerCode(code string, email string) error {
	cfg := config.Config.Server.Email
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Username)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Code")
	msg := fmt.Sprintf("Your verification code is %s", code)
	m.SetBody("text/html", msg)
	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	return d.DialAndSend(m)
}