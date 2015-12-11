package pubsub

import (
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/internal"
)

type TopicHandlerList struct {
	sync.RWMutex
	handlers []app.TopicHandler
}

func NewTopicHandlerList(handlers ...app.TopicHandler) *TopicHandlerList {
	return &TopicHandlerList{
		handlers: handlers,
	}
}

func (this *TopicHandlerList) Add(handler app.TopicHandler) {
	this.Lock()
	this.handlers = append(this.handlers, handler)
	this.Unlock()
}

func (this *TopicHandlerList) Dispatch(message app.PubsubMessage) {
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

func New() internal.PubsubWithContext {
	return &Pubsub{
		topics: make(map[string]*TopicHandlerList),
	}
}

func (this *Pubsub) Publish(context app.Context, topic string, message app.PubsubMessage) {
	this.RLock()
	list, found := this.topics[topic]
	this.RUnlock()

	if found {
		list.Dispatch(message)
	}
}

func (this *Pubsub) Subscribe(context app.Context, topic string, handler app.TopicHandler) {
	this.Lock()
	defer this.Unlock()
	if list, found := this.topics[topic]; found {
		list.Add(handler)
	} else {
		this.topics[topic] = NewTopicHandlerList(handler)
	}
}
