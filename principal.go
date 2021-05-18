package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func main() {
	// 连接："amqp://用户名:密码@服务地址:服务端口/"
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnErrPrincipal(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrPrincipal(err, "无法打开信道")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"principal",
		true, // durable 持久性
		false,
		false,
		false,
		nil,
	)
	failOnErrPrincipal(err, "无法声明队列")

	args := os.Args
	body := strings.Join(args[1:], " ")
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnErrPrincipal(err, "消息发送失败")
	log.Println("消息发送成功！")
}

// failOnErr 检查异常并终断程序输出错误
func failOnErrPrincipal(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
