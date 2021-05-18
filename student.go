package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main() {
	// 连接："amqp://用户名:密码@服务地址:服务端口/"
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnErrStudent(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrStudent(err, "无法打开信道")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"school",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErrStudent(err, "无法声明队列")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnErrStudent(err, "注册消费者失败")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到来自学校的信息：%s", d.Body)
		}
	}()
	<-forever
}

// failOnErrStudent 检查异常并终断程序输出错误
func failOnErrStudent(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
