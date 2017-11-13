//Collects Mpeg2TSPackets and builds a single frame (access unit) out of them.
//Frames are added to a PooledQueue given upon creation.
package main

import (
	"container/list"
	"fmt"
	"log"
)

// todo - remove. this is just for syntax errors ::
type Mpeg2TSPacket struct {}
func (non Mpeg2TSPacket) getContinuityCounter () int {return 0}
func (non Mpeg2TSPacket) isPayloadExist () bool {return true}
func (non Mpeg2TSPacket) isStartOfPES () bool {return true}
// todo - end of list.


const MAX_COUNTER = 16

type Frame struct{
	array int		// todo - fill this struct to something meaningful.
}


type Mpeg2TSParser struct {
	counter int
    output list.List
	currentFrame * Frame
	iFrameFound bool
	lastFrameTime uint64
	endFlag bool
}

 func (ps *Mpeg2TSParser) New(){
 	ps.counter = 0
 	ps.lastFrameTime = 0
 	ps.iFrameFound = false
 	ps.endFlag = false
 	// todo - complete this function.
 }

 func (ps *Mpeg2TSParser) getPreviousCounter(counter int) int{
	if counter == 0 {
		return 15
	}
	return counter - 1
 }


 func (ps *Mpeg2TSParser) write(packet Mpeg2TSPacket) {
 	var packetCounter int =  packet.getContinuityCounter()

	 if !packet.isPayloadExist() {
		 fmt.Println("Payload doens't exist")
		 return
	 }

	 if packetCounter == ps.getPreviousCounter(ps.counter) {
		 fmt.Println("Duplicate packet found")
		 return
	 }

	 if ps.currentFrame == nil {
	 	 frame := new(Frame)
  		 ps.output.PushBack(frame)
		 ps.currentFrame  = frame
		 //ps.currentFrame.clear(); // todo - make a matching method.
	 }

	 if ps.currentFrame == nil {
		 fmt.Println("currentFrame is null")
		 Sleeper.sleep(5)			// todo - make SLEEPER!
		 return
	 }

	 if packet.isStartOfPES() {
		 if !ps.currentFrame.isEmpty() { //todo - make isEMpty()
			 if packetCounter == ps.counter {
				 for true {
					 ps.output.PushBack(ps.currentFrame)
					 if ps.endFlag {
						 break;
					 }
				 }
			 }
		 }
	 } 
 }