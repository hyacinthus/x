package xmq

import (
	"encoding/json"

	"time"

	"github.com/levigross/grequests"
	nsq "github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// NSQ context
type nsqContext struct {
	msg *nsq.Message
}

// Bind extract payload data to v
func (c *nsqContext) Bind(v interface{}) error {
	return json.Unmarshal(c.msg.Body, v)
}

// NSQ 客户端
type nsqClient struct {
	producer *nsq.Producer
	config   Config
}

// 新建客户端
func newNSQClient(config Config) Client {
	var err error

	var c = &nsqClient{
		config: config,
	}

	c.producer, err = nsq.NewProducer(config.PubHost+":"+config.PubTCP, nsq.NewConfig())
	if err != nil {
		logrus.WithError(err).Panic("init nsq producer failed")
	}
	c.producer.SetLogger(NewLogrusLoggerAtLevel(logrus.WarnLevel))
	log.Info("NSQ Producer 初始化完成。")
	return c
}

// 为了对外隐藏nsq包的对象，把handler转换成使用context接口
func decorate(f HandlerFunc) nsq.HandlerFunc {
	return func(msg *nsq.Message) error {
		c := &nsqContext{
			msg: msg,
		}
		return f(c)
	}
}

// Reg 注册一个消费者处理函数 func(msg *nsq.Message) error
func (c *nsqClient) Sub(topic, channel string, handler HandlerFunc) {
	q, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		log.WithError(err).Panic("init nsq comsumer failed")
	}
	q.SetLogger(NewLogrusLoggerAtLevel(logrus.WarnLevel))
	q.AddHandler(decorate(handler))
	err = q.ConnectToNSQLookupd(c.config.SubHost + ":" + c.config.SubHTTP)
	if err != nil {
		log.WithError(err).Panic("nsq comsumer connect to lookupd failed")
	}
	log.Infof("订阅 nsq topic %s by %s", topic, channel)
}

// Pub 发布消息，用json编码
func (c *nsqClient) Pub(topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.producer.Publish(topic, data)
}

// Delay 延迟发布消息，用json编码
func (c *nsqClient) Delay(topic string, payload interface{}, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.producer.DeferredPublish(topic, delay, data)
}

// CreateTopic create a topic on nsqd by http request
func (c *nsqClient) CreateTopic(topic string) error {
	resp, err := grequests.Post("http://"+c.config.PubHost+":"+c.config.PubHTTP+"/topic/create",
		&grequests.RequestOptions{Data: map[string]string{"topic": topic}})
	if err != nil {
		return err
	}
	log.Info(resp)
	return nil
}
