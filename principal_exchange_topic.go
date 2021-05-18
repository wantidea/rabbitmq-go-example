package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func main() {
	// 参数录入
	routingKey := ""
	args := os.Args
	if len(args) > 2 {
		routingKey = args[1]
	} else {
		log.Fatalln("请输入绑定的 RoutingKey 与消息")
	}
	body := strings.Join(args[2:], " ")

	// 连接："amqp://用户名:密码@服务地址:服务端口/"
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	failOnErrPrincipalExchangeTopic(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrPrincipalExchangeTopic(err, "无法打开信道")
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
	failOnErrPrincipalExchangeTopic(err, "无法声明交换机")

	// 此处模拟校长学生信息给教导主任
	err = ch.Publish(
		"instructorTopic", // exchange 交换机名称
		routingKey,        // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnErrPrincipalExchangeTopic(err, "消息发送失败")
	log.Println("消息发送成功！")
}

// failOnErr 检查异常并终断程序输出错误
func failOnErrPrincipalExchangeTopic(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
