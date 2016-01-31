package pubsub

import (
	"github.com/v2ray/v2ray-core/app"
)

const (
	APP_ID = app.ID(3)
)

type PubsubMessage []byte
type TopicHandler func(PubsubMessage)

type Pubsub interface {
	Publish(topic string, message PubsubMessage)
	Subscribe(topic string, handler TopicHandler)
}

type pubsubWithContext interface {
	Publish(context app.Context, topic string, message PubsubMessage)
	Subscribe(context app.Context, topic string, handler TopicHandler)
}

type contextedPubsub struct {
	context app.Context
	pubsub  pubsubWithContext
}

func (this *contextedPubsub) Publish(topic string, message PubsubMessage) {
	this.pubsub.Publish(this.context, topic, message)
}

func (this *contextedPubsub) Subscribe(topic string, handler TopicHandler) {
	this.pubsub.Subscribe(this.context, topic, handler)
}

func init() {
	app.RegisterApp(APP_ID, func(context app.Context, obj interface{}) interface{} {
		pubsub := obj.(pubsubWithContext)
		return &contextedPubsub{
			context: context,
			pubsub:  pubsub,
		}
	})
}
