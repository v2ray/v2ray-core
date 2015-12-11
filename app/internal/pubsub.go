package internal

import (
	"github.com/v2ray/v2ray-core/app"
)

type PubsubWithContext interface {
	Publish(context app.Context, topic string, message app.PubsubMessage)
	Subscribe(context app.Context, topic string, handler app.TopicHandler)
}

type contextedPubsub struct {
	context app.Context
	pubsub  PubsubWithContext
}

func (this *contextedPubsub) Publish(topic string, message app.PubsubMessage) {
	this.pubsub.Publish(this.context, topic, message)
}

func (this *contextedPubsub) Subscribe(topic string, handler app.TopicHandler) {
	this.pubsub.Subscribe(this.context, topic, handler)
}
