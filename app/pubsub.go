package app

type PubsubMessage []byte
type TopicHandler func(PubsubMessage)

type Pubsub interface {
	Publish(topic string, message PubsubMessage)
	Subscribe(topic string, handler TopicHandler)
}
