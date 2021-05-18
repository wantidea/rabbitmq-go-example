package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	// 连接："amqp://用户名:密码@服务地址:服务端口/"
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnErrTeacherExchangeHeaders(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrTeacherExchangeHeaders(err, "无法打开信道")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"instructorHeaders", // name 交换机名称
		"headers",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnErrTeacherExchangeHeaders(err, "无法声明交换机")

	// 声明一个全新的非持久空队列
	q, err := ch.QueueDeclare(
		"",    // queue name 队列名称
		false, // durable 持久性
		false,
		true, // exclusive
		false,
		nil,
	)
	failOnErrTeacherExchangeHeaders(err, "无法声明队列")

	// 绑定队列
	err = ch.QueueBind(
		q.Name, // queue name 队列名称
		"",
		"instructorHeaders", // exchange 交换机名称
		false,
		amqp.Table{
			"x-match": "all", // any 为匹配一个，all 需要匹配所有
			"teacher": "t_chinese",
			"student": "s_ming",
		},
	)
	failOnErrTeacherExchangeHeaders(err, "队列绑定失败")

	msgs, err := ch.Consume(
		q.Name, // queue name 队列名称
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	failOnErrTeacherExchangeHeaders(err, "注册消费者失败")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到来自教导主任的学生名单：%s", d.Body)
			log.Printf("电话通知：%s ...", d.Body)

			// 模拟电话通知时长
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)

			log.Printf("已通知学生：%s", d.Body)
		}
	}()
	<-forever
}

// failOnErrTeacherExchangeHeaders 检查异常并终断程序输出错误
func failOnErrTeacherExchangeHeaders(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
