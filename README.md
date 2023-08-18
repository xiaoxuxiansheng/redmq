# redmq
<p align="center">
<img src="https://github.com/xiaoxuxiansheng/redmq/blob/main/img/redmq_frame.png" height="400px/"><br/><br/>
<b>redmq: 纯 redis 实现的消息队列</b>
<br/><br/>
</p>

## 📚 前言
使用此 sdk 进行实践前，建议先行了解与 redis streams 有关的特性，做到知行合一<br/><br/>
<a href="https://redis.io/docs/data-types/streams/">redis streams</a> <br/><br/>

## 💡 `redmq` 技术原理
<a href="https://xxxx">基于 redis 实现消息队列</a> <br/><br/>

## 🖥 接入 sop
用户需要先行完成 topic 和 consumer group 的创建<br/><br/>
- 创建 topic：my_test_topic<br/><br/>
```redis
127.0.0.1:6379> xadd my_test_topic * first_key first_val
"1692066364494-0"
```
- 创建 consumer group<br/><br/>
```redis
127.0.0.1:6379> XGROUP CREATE my_test_topic my_test_group 0-0
OK
```
- 构造 redis 客户端实例<br/><br/>
```go
import "github.com/xiaoxuxiansheng/redmq/redis"
func main(){
    redisClient := redis.NewClient("tcp","my_address","my_password")
    // ...
}
```

- 启动生产者 producer<br/><br/>
```go
import (
	"context"

	"github.com/xiaoxuxiansheng/redmq"
)
func main(){
    // ...
	producer := redmq.NewProducer(redisClient, redmq.WithMsgQueueLen(10))
	ctx := context.Background()
	msgID, err := producer.SendMsg(ctx, topic, "test_kk", "test_vv")
}
```

- 启动消费者 consumer<br/><br/>
```go
import (
	"github.com/xiaoxuxiansheng/redmq"
)
func main(){
    // ...
    // 构造并启动消费者
	consumer, _ := redmq.NewConsumer(redisClient, topic, consumerGroup, consumerID, callbackFunc,
		// 每条消息最多重试 2 次
		redmq.WithMaxRetryLimit(2),
		// 每轮接收消息的超时时间为 2 s
		redmq.WithReceiveTimeout(2*time.Second),
		// 注入自定义实现的死信队列
		redmq.WithDeadLetterMailbox(demoDeadLetterMailbox))
	defer consumer.Stop()
}
```

## 🐧 使用示例
完整的使用示例代码也可以参见 package example：
- mock 生产者投递消息流程<br/><br/>
```go
import (
	"context"
	"testing"

	"github.com/xiaoxuxiansheng/redmq"
	"github.com/xiaoxuxiansheng/redmq/redis"
)

const (
	network  = "tcp"
	address  = "请输入 redis 地址"
	password = "请输入 redis 密码"
	topic    = "请输入 topic 名称"
)

func Test_Producer(t *testing.T) {
	client := redis.NewClient(network, address, password)
	// 最多保留十条消息
	producer := redmq.NewProducer(client, redmq.WithMsgQueueLen(10))
	ctx := context.Background()
	msgID, err := producer.SendMsg(ctx, topic, "test_k", "test_v")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(msgID)
}
```
- mock 消费者消费消息流程<br></br> 
```go
import (
	"context"
	"testing"
	"time"

	"github.com/xiaoxuxiansheng/redmq"
	"github.com/xiaoxuxiansheng/redmq/redis"
)

const (
	network       = "tcp"
	address       = "请输入 redis 地址"
	password      = "请输入 redis 密码"
	topic         = "请输入 topic 名称"
	consumerGroup = "请输入消费者组名称"
	consumerID    = "请输入消费者名称"
)

// 自定义实现的死信队列
type DemoDeadLetterMailbox struct {
	do func(msg *redis.MsgEntity)
}

func NewDemoDeadLetterMailbox(do func(msg *redis.MsgEntity)) *DemoDeadLetterMailbox {
	return &DemoDeadLetterMailbox{
		do: do,
	}
}

// 死信队列接收消息的处理方法
func (d *DemoDeadLetterMailbox) Deliver(ctx context.Context, msg *redis.MsgEntity) error {
	d.do(msg)
	return nil
}

func Test_Consumer(t *testing.T) {
	client := redis.NewClient(network, address, password)

	// 接收到消息后的处理函数
	callbackFunc := func(ctx context.Context, msg *redis.MsgEntity) error {
		t.Logf("receive msg, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
		return nil
	}

	// 自定义实现的死信队列
	demoDeadLetterMailbox := NewDemoDeadLetterMailbox(func(msg *redis.MsgEntity) {
		t.Logf("receive dead letter, msg id: %s, msg key: %s, msg val: %s", msg.MsgID, msg.Key, msg.Val)
	})

	// 构造并启动消费者
	consumer, err := redmq.NewConsumer(client, topic, consumerGroup, consumerID, callbackFunc,
		// 每条消息最多重试 2 次
		redmq.WithMaxRetryLimit(2),
		// 每轮接收消息的超时时间为 2 s
		redmq.WithReceiveTimeout(2*time.Second),
		// 注入自定义实现的死信队列
		redmq.WithDeadLetterMailbox(demoDeadLetterMailbox))
	if err != nil {
		t.Error(err)
		return
	}
	defer consumer.Stop()

	// 十秒后退出单测程序
	<-time.After(10 * time.Second)
}
```
