package collect

import (
	"container/heap"
	"sync"
	"time"
)

type timedQueueEntry struct {
	timeSec int64
	value   interface{}
}

type timedQueue []*timedQueueEntry

func (queue timedQueue) Len() int {
	return len(queue)
}

func (queue timedQueue) Less(i, j int) bool {
	return queue[i].timeSec < queue[j].timeSec
}

func (queue timedQueue) Swap(i, j int) {
	tmp := queue[i]
	queue[i] = queue[j]
	queue[j] = tmp
}

func (queue *timedQueue) Push(value interface{}) {
	entry := value.(*timedQueueEntry)
	*queue = append(*queue, entry)
}

func (queue *timedQueue) Pop() interface{} {
	old := *queue
	n := len(old)
	v := old[n-1]
	*queue = old[:n-1]
	return v
}

type TimedStringMap struct {
	timedQueue
	queueMutex sync.Mutex
	dataMutext sync.RWMutex
	data       map[string]interface{}
	interval   int
}

func NewTimedStringMap(updateInterval int) *TimedStringMap {
	m := &TimedStringMap{
		timedQueue: make([]*timedQueueEntry, 0, 1024),
		queueMutex: sync.Mutex{},
		dataMutext: sync.RWMutex{},
		data:       make(map[string]interface{}, 1024),
		interval:   updateInterval,
	}
	m.initialize()
	return m
}

func (m *TimedStringMap) initialize() {
	go m.cleanup(time.Tick(time.Duration(m.interval) * time.Second))
}

func (m *TimedStringMap) cleanup(tick <-chan time.Time) {
	for {
		now := <-tick
		nowSec := now.UTC().Unix()
		if m.timedQueue.Len() == 0 {
			continue
		}
		for m.timedQueue.Len() > 0 {
			entry := m.timedQueue[0]
			if entry.timeSec > nowSec {
				break
			}
			m.queueMutex.Lock()
			entry = heap.Pop(&m.timedQueue).(*timedQueueEntry)
			m.queueMutex.Unlock()
			m.Remove(entry.value.(string))
		}
	}
}

func (m *TimedStringMap) Get(key string) (interface{}, bool) {
	m.dataMutext.RLock()
	value, ok := m.data[key]
	m.dataMutext.RUnlock()
	return value, ok
}

func (m *TimedStringMap) Set(key string, value interface{}, time2Delete int64) {
	m.dataMutext.Lock()
	m.data[key] = value
	m.dataMutext.Unlock()

	m.queueMutex.Lock()
	heap.Push(&m.timedQueue, &timedQueueEntry{
		timeSec: time2Delete,
		value:   key,
	})
	m.queueMutex.Unlock()
}

func (m *TimedStringMap) Remove(key string) {
	m.dataMutext.Lock()
	delete(m.data, key)
	m.dataMutext.Unlock()
}
