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
	failOnErrPrincipalExchangeHeaders(err, "无法连接 RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrPrincipalExchangeHeaders(err, "无法打开信道")
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
	failOnErrPrincipalExchangeHeaders(err, "无法声明交换机")

	// 此处模拟校长学生信息给教导主任
	args := os.Args
	body := strings.Join(args[1:], " ")
	err = ch.Publish(
		"instructorHeaders", // exchange 交换机名称
		"",                  // routing key 当前设置为空字符串
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
			Headers: amqp.Table{
				"teacher": "t_chinese",
				"student": "s_ming",
			},
		},
	)
	failOnErrPrincipalExchangeHeaders(err, "消息发送失败")
	log.Println("消息发送成功！")
}

// failOnErr 检查异常并终断程序输出错误
func failOnErrPrincipalExchangeHeaders(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
