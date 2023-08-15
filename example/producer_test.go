package example

import (
	"context"
	"testing"

	"github.com/xiaoxuxiansheng/redmq"
	"github.com/xiaoxuxiansheng/redmq/redis"
)

func Test_Producer(t *testing.T) {
	client := redis.NewClient(network, address, password)
	// 最多保留十条消息
	producer := redmq.NewProducer(client, redmq.WithMsgQueueLen(10))
	ctx := context.Background()
	msgID, err := producer.SendMsg(ctx, topic, "test_kk", "test_vv")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(msgID)
}
