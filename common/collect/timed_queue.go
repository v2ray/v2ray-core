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
	queue           timedQueueImpl
	access          sync.Mutex
	removedCallback func(interface{})
}

func NewTimedQueue(updateInterval int, removedCallback func(interface{})) *TimedQueue {
	queue := &TimedQueue{
		queue:           make([]*timedQueueEntry, 0, 256),
		removedCallback: removedCallback,
		access:          sync.Mutex{},
	}
	go queue.cleanup(time.Tick(time.Duration(updateInterval) * time.Second))
	return queue
}

func (queue *TimedQueue) Add(value interface{}, time2Remove int64) {
	newEntry := &timedQueueEntry{
		timeSec: time2Remove,
		value:   value,
	}
	var removedEntry *timedQueueEntry
	queue.access.Lock()
	nowSec := time.Now().Unix()
	if queue.queue.Len() > 0 && queue.queue[0].timeSec < nowSec {
		removedEntry = queue.queue[0]
		queue.queue[0] = newEntry
		heap.Fix(&queue.queue, 0)
	} else {
		heap.Push(&queue.queue, newEntry)
	}
	queue.access.Unlock()
	if removedEntry != nil {
		queue.removedCallback(removedEntry.value)
	}
}

func (queue *TimedQueue) cleanup(tick <-chan time.Time) {
	for now := range tick {
		nowSec := now.Unix()
		removedEntries := make([]*timedQueueEntry, 0, 128)
		queue.access.Lock()
		changed := false
		for i := 0; i < queue.queue.Len(); i++ {
			entry := queue.queue[i]
			if entry.timeSec < nowSec {
				removedEntries = append(removedEntries, entry)
				queue.queue.Swap(i, queue.queue.Len()-1)
				queue.queue.Pop()
				changed = true
			}
		}
		if changed {
			heap.Init(&queue.queue)
		}
		queue.access.Unlock()
		for _, entry := range removedEntries {
			queue.removedCallback(entry.value)
		}
	}
}
