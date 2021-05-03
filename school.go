package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// 连接："amqp://用户名:密码@服务地址:服务端口/"
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnErr(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErr(err, "无法打开信道")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"school",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErr(err, "无法声明队列")

	body := "同学们，因...，明天不用来上学啦"
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
	failOnErr(err, "消息发送失败")
	log.Println("消息发送成功！")
}

// failOnErr 检查异常并终断程序输出错误
func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
