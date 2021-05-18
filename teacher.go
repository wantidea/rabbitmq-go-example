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
	failOnErrTeacher(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrTeacher(err, "无法打开信道")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"principal",
		true, // durable 持久性
		false,
		false,
		false,
		nil,
	)
	failOnErrTeacher(err, "无法声明队列")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	failOnErrTeacher(err, "注册消费者失败")

	err = ch.Qos(
		1, // prefetch count 预取计数
		0,
		false,
	)
	failOnErrTeacher(err, "无法设置QoS")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到来自校长的学生名单：%s", d.Body)
			log.Printf("电话通知：%s ...", d.Body)

			// 模拟电话通知时长
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)

			log.Printf("已通知学生：%s", d.Body)
			d.Ack(false)
		}
	}()
	<-forever
}

// failOnErrTeacher 检查异常并终断程序输出错误
func failOnErrTeacher(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
