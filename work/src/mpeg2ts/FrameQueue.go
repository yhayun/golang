package main

import "fmt"

//   todo - refernce:
//   https://stackoverflow.com/questions/25657207/golang-how-to-know-a-buffered-channel-is-full

type frameQueue struct {
	queue chan Frame
}

//CONSTRUCTOR:
func NewFrameQueue(capacity int, frameCapacity int) *frameQueue {
	var res *frameQueue =  &frameQueue{
		queue: make(chan Frame, capacity),
	}
	for i:=0 ; i < capacity; i++ {
		res.queue <- *NewFrame(frameCapacity)
	}

	return res
}


/**
 * @return A new pre-allocated instance of type E.
 */
func (q frameQueue) NewElement() *Frame {
	select {
	case tmp := <- q.queue:
		fmt.Println("yes")
		return &tmp
	default:
		fmt.Println("no")
		return nil
	}
}

/**
 *  Wrapper around NewElement, does the same thing.
 */
func (q frameQueue) Poll() *Frame {
	return q.NewElement()
}

/**
 * @param e
 *            Recycles a given instance of type E back to the pool.
 */
func (q frameQueue) Recylce(f *Frame) {
	if f != nil {
		f.Clear()
		q.queue <- *f
	}
}

//This behaves as a non-blocking put.
func (q frameQueue) Offer(f *Frame) bool {
	select {
	case q.queue <- *f:
		return true
	default:
		return false
	}
}


func ( q frameQueue) IsEmpty() bool {
	return 0 == len(q.queue)
}