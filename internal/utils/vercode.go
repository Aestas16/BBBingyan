package utils

import (
    "fmt"
    "math/rand"
    "time"
    "gopkg.in/gomail.v2"
    "github.com/redis/go-redis/v9"

    "user-management-system/internal/config"
)

var RedisClient *redis.Client

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     config.Config.Redis.Host,
        Password: config.Config.Redis.Password,
        DB:       config.Config.Redis.DB,
    })
    _, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        panic(err)
    }
}

func CreateVerCode(email string, code string) error {
    ctx := context.Background()
    err := RedisClient.Set(ctx, email, code, 0).Err()
    if err != nil {
        return err
    }
    err = RedisClient.Expire(ctx, email, time.Duration(config.Config.Captcha.Expire)*time.Second).Err()
    return err
}

func ValidateVerCode(email string, code string) (bool, error) {
    ctx := context.Background()
    vercode, err := RedisClient.Get(ctx, email).Result()
    if err != nil {
        return false, err
    }
    if vercode == code {
        return true, nil
    }
    return false, nil
}

func GenerateVerCode(length int) string {
    rand.Seed(time.Now().UnixNano())
    var chars = []rune("0123456789")
    b := make([]rune, length)
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