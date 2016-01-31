package internal

import (
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/pubsub"
)

type TopicHandlerList struct {
	sync.RWMutex
	handlers []pubsub.TopicHandler
}

func NewTopicHandlerList(handlers ...pubsub.TopicHandler) *TopicHandlerList {
	return &TopicHandlerList{
		handlers: handlers,
	}
}

func (this *TopicHandlerList) Add(handler pubsub.TopicHandler) {
	this.Lock()
	this.handlers = append(this.handlers, handler)
	this.Unlock()
}

func (this *TopicHandlerList) Dispatch(message pubsub.PubsubMessage) {
	this.RLock()
	for _, handler := range this.handlers {
		go handler(message)
	}
	this.RUnlock()
}

type Pubsub struct {
	topics map[string]*TopicHandlerList
	sync.RWMutex
}

func New() *Pubsub {
	return &Pubsub{
		topics: make(map[string]*TopicHandlerList),
	}
}

func (this *Pubsub) Publish(context app.Context, topic string, message pubsub.PubsubMessage) {
	this.RLock()
	list, found := this.topics[topic]
	this.RUnlock()

	if found {
		list.Dispatch(message)
	}
}

func (this *Pubsub) Subscribe(context app.Context, topic string, handler pubsub.TopicHandler) {
	this.Lock()
	defer this.Unlock()
	if list, found := this.topics[topic]; found {
		list.Add(handler)
	} else {
		this.topics[topic] = NewTopicHandlerList(handler)
	}
}
