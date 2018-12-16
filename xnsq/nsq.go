package xnsq

import (
	"github.com/hyacinthus/x/xlog"
	nsq "github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

var log = xlog.Get()
var nsqdAddr, nsqLookupdAddr string

// Init 初始化nsq相关地址
// nsqd 需要tcp地址 默认端口4150
// nsqlookupd 需要http地址 默认端口4161
func Init(nsqd, nsqlookupd string) {
	nsqdAddr = nsqd
	nsqLookupdAddr = nsqlookupd
}

// Producer 是 nsq 的生产者
func Producer() *nsq.Producer {
	producer, err := nsq.NewProducer(nsqdAddr, nsq.NewConfig())
	if err != nil {
		logrus.WithError(err).Panic("init nsq producer failed")
	}
	producer.SetLogger(NewLogrusLoggerAtLevel(logrus.WarnLevel))
	log.Info("NSQ Producer 初始化完成。")
	return producer
}

// Reg 注册一个消费者处理函数 func(msg *nsq.Message) error
func Reg(topic, channel string, handler nsq.HandlerFunc) {
	q, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		log.WithError(err).Panic("init nsq comsumer failed")
	}
	q.SetLogger(NewLogrusLoggerAtLevel(logrus.WarnLevel))
	q.AddHandler(handler)
	err = q.ConnectToNSQLookupd(nsqLookupdAddr)
	if err != nil {
		log.WithError(err).Panic("nsq comsumer connect to lookupd failed")
	}
	log.Infof("订阅 nsq topic %s at %s", topic, channel)
}
