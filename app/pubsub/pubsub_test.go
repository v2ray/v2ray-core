package pubsub_test

import (
	"testing"
	"time"

	"github.com/v2ray/v2ray-core/app"
	. "github.com/v2ray/v2ray-core/app/pubsub"
	apptesting "github.com/v2ray/v2ray-core/app/testing"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestPubsub(t *testing.T) {
	v2testing.Current(t)

	messages := make(map[string]app.PubsubMessage)

	pubsub := New()
	pubsub.Subscribe(&apptesting.Context{}, "t1", func(message app.PubsubMessage) {
		messages["t1"] = message
	})

	pubsub.Subscribe(&apptesting.Context{}, "t2", func(message app.PubsubMessage) {
		messages["t2"] = message
	})

	message := app.PubsubMessage([]byte("This is a pubsub message."))
	pubsub.Publish(&apptesting.Context{}, "t2", message)
	<-time.Tick(time.Second)

	_, found := messages["t1"]
	assert.Bool(found).IsFalse()

	actualMessage, found := messages["t2"]
	assert.Bool(found).IsTrue()
	assert.StringLiteral(string(actualMessage)).Equals(string(message))
}
