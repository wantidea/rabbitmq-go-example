package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

func main() {
	// 参数录入
	routingKey := ""
	args := os.Args
	if len(args) > 1 {
		routingKey = args[1]
	} else {
		log.Fatalln("请输入绑定的 RoutingKey")
	}

	// 连接："amqp://用户名:密码@服务地址:服务端口/"
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnErrTeacherExchangeTopic(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrTeacherExchangeTopic(err, "无法打开信道")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"instructorTopic", // name 交换机名称
		"topic",           // 交换机类型
		true,
		false,
		false,
		false,
		nil,
	)
	failOnErrTeacherExchangeTopic(err, "无法声明交换机")

	// 声明一个全新的非持久空队列
	q, err := ch.QueueDeclare(
		"",    // queue name 队列名称
		false, // durable 持久性
		false,
		true, // exclusive
		false,
		nil,
	)
	failOnErrTeacherExchangeTopic(err, "无法声明队列")

	// 绑定队列
	err = ch.QueueBind(
		q.Name,            // queue name 队列名称
		routingKey,        // routingKey 路由键
		"instructorTopic", // exchange 交换机名称
		false,
		nil,
	)
	failOnErrTeacherExchangeTopic(err, "队列绑定失败")

	msgs, err := ch.Consume(
		q.Name, // queue name 队列名称
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	failOnErrTeacherExchangeTopic(err, "注册消费者失败")

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

// failOnErrTeacherExchangeTopic 检查异常并终断程序输出错误
func failOnErrTeacherExchangeTopic(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
