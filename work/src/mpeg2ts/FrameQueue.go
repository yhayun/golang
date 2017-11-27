package main

//   todo - refernce:
//   https://stackoverflow.com/questions/25657207/golang-how-to-know-a-buffered-channel-is-full

type frameQueue struct {
	queue chan frame
}

//CONSTRUCTOR:
func NewFrameQueue(capacity int, frameCapacity int) *frameQueue {
	var res *frameQueue =  &frameQueue{
		queue: make(chan frame, capacity),
	}
	for i:=0 ; i < capacity; i++ {
		res.queue <- *NewFrame(frameCapacity)
	}

	return res
}


/**
 * @return A new pre-allocated instance of type E.
 */
func (q frameQueue) NewElement() frame {
	return <- q.queue
}

/**
 * @param e
 *            Recycles a given instance of type E back to the pool.
 */
func (q frameQueue) Recylce(f *frame) {
	if f != nil {
		f.Clear()
		q.queue <- *f
	}
}

//This behaves as a non-blocking put.
func (q frameQueue) Offer(f *frame) bool {
	select {
	case q.queue <- *f:
		return true
	default:
		return false
	}
}


