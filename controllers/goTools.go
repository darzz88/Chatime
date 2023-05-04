package controllers

import (
	"context"

	"github.com/go-redis/redis/v8"
	gomail "gopkg.in/mail.v2"
)

var ctx = context.Background()

func SetRedis(rdb *redis.Client, key string, value string, expiration int) {
	err := rdb.Set(ctx, key, value, 0).Err()
	CheckError(err)
}

func GetRedis(rdb *redis.Client, key string) string {
	val, err := rdb.Get(ctx, key).Result()

	CheckError(err)
	return val
}

func SendMail(to string, subject string, text string) {
	gmail := gomail.NewMessage()

	gmail.SetHeader("From", "if-20001@students.ithb.ac.id")
	gmail.SetHeader("To", to)
	gmail.SetHeader("Subject", subject)
	gmail.SetBody("text/plain", text)

	gm := gomail.NewDialer("smtp.gmail.com", 587, "chatimepbp@gmail.com", "phvmkljjxdzlpsxf")

	err := gm.DialAndSend(gmail)
	CheckError(err)
}

func Task() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	eng := GetRedis(rdb, "eng")
	idn := GetRedis(rdb, "idn")
	go SendMail("chatimepbp@gmail.com", "Chatime - Reminder", eng)
	SendMail("chatimepbp@gmail.com", "Chatime - Pengingat", idn)
}
