package collect

type SizedQueue struct {
	elements []interface{}
	nextPos  int
}

func NewSizedQueue(size int) *SizedQueue {
	return &SizedQueue{
		elements: make([]interface{}, size),
		nextPos:  0,
	}
}

// Put puts a new element into the queue and pop out the first element if queue is full.
func (this *SizedQueue) Put(element interface{}) interface{} {
	res := this.elements[this.nextPos]
	this.elements[this.nextPos] = element
	this.nextPos++
	if this.nextPos == len(this.elements) {
		this.nextPos = 0
	}
	return res
}
