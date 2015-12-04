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

type timedQueueImpl []*timedQueueEntry

func (queue timedQueueImpl) Len() int {
	return len(queue)
}

func (queue timedQueueImpl) Less(i, j int) bool {
	return queue[i].timeSec < queue[j].timeSec
}

func (queue timedQueueImpl) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
}

func (queue *timedQueueImpl) Push(value interface{}) {
	entry := value.(*timedQueueEntry)
	*queue = append(*queue, entry)
}

func (queue *timedQueueImpl) Pop() interface{} {
	old := *queue
	n := len(old)
	v := old[n-1]
	old[n-1] = nil
	*queue = old[:n-1]
	return v
}

// TimedQueue is a priority queue that entries with oldest timestamp get removed first.
type TimedQueue struct {
	queue   timedQueueImpl
	access  sync.RWMutex
	removed chan interface{}
}

func NewTimedQueue(updateInterval int) *TimedQueue {
	queue := &TimedQueue{
		queue:   make([]*timedQueueEntry, 0, 256),
		removed: make(chan interface{}, 16),
		access:  sync.RWMutex{},
	}
	go queue.cleanup(time.Tick(time.Duration(updateInterval) * time.Second))
	return queue
}

func (queue *TimedQueue) Add(value interface{}, time2Remove int64) {
	queue.access.Lock()
	heap.Push(&queue.queue, &timedQueueEntry{
		timeSec: time2Remove,
		value:   value,
	})
	queue.access.Unlock()
}

func (queue *TimedQueue) RemovedEntries() <-chan interface{} {
	return queue.removed
}

func (queue *TimedQueue) cleanup(tick <-chan time.Time) {
	for now := range tick {
		nowSec := now.Unix()
		for {
			queue.access.RLock()
			queueLen := queue.queue.Len()
			queue.access.RUnlock()
			if queueLen == 0 {
				break
			}
			queue.access.RLock()
			entry := queue.queue[0]
			queue.access.RUnlock()
			if entry.timeSec > nowSec {
				break
			}
			queue.access.Lock()
			heap.Pop(&queue.queue)
			queue.access.Unlock()

			queue.removed <- entry.value
		}
	}
}
