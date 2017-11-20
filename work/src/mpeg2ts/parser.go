//Collects Mpeg2TSPackets and builds a single frame (access unit) out of them.
//Frames are added to a PooledQueue given upon creation.
package main

import (
	"fmt"
	"time"
)

// todo - remove. this is just for syntax errors ::
func (non Mpeg2TSPacket) getContinuityCounter () int {return 0}
func (non Mpeg2TSPacket) isPayloadExist () bool {return true}
func (non Mpeg2TSPacket) isStartOfPES () bool {return true}
func (non Mpeg2TSPacket) getPTS () int64 {return 0}

func (non Mpeg2TSPacket) getData () []byte {return nil}
func (non Mpeg2TSPacket) getPayloadOffset () int {return 0}
func (non Mpeg2TSPacket) getPayloadLength () int {return 0}
// todo - end of list.

const MAX_COUNTER = 16

type Mpeg2TSParser struct {
	counter int
    output frameQueue
	currentFrame * frame
	iFrameFound bool
	lastFrameTime int64
	endFlag bool
}

 func (ps *Mpeg2TSParser) NewMpeg2TSParser(output frameQueue) *Mpeg2TSParser{
 	ps.counter = 0
 	ps.lastFrameTime = 0
 	ps.iFrameFound = false
 	ps.endFlag = false
	 return &Mpeg2TSParser{
		 counter: -1,
		 lastFrameTime: 0,
		 iFrameFound: false,
		 endFlag: false,
		 output: output,
	 }
 	// todo - complete this function.
 }

 func (ps Mpeg2TSParser) getPreviousCounter(counter int) int{
	if counter == 0 {
		return 15
	}
	return counter - 1
 }

 func (ps Mpeg2TSParser) close() {
	ps.endFlag = true;
 }


/**
* Deletes the content of the current frame and recycles it.
*/
 func (ps Mpeg2TSParser) flush() {
	ps.output.recylce(ps.currentFrame)
	ps.currentFrame = nil
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
	 	var temp frame = ps.output.newElement()
 		ps.currentFrame = &temp
 		ps.currentFrame.clear()
	 }

	 if ps.currentFrame == nil { //todo - Makes no sense.
		 fmt.Println("currentFrame is null")
		 time.Sleep(5 * time.Millisecond)
		 return
	 }

	 if packet.isStartOfPES() {
		 if !ps.currentFrame.isEmpty() {
			 if packetCounter == ps.counter {
			 	for !ps.output.offer(ps.currentFrame) {
			 		if (ps.endFlag){
			 			break
					}
					time.Sleep(5 * time.Millisecond)
				}
				 var temp frame = ps.output.newElement()
				 ps.currentFrame = &temp
				 ps.currentFrame.clear();

				 //todo check if frame is null again - makes no sense in GO
			 } else {
				 fmt.Println("currentFrame.clear()");
				 ps.currentFrame.clear();
			 }
		 }
		 ps.currentFrame.setPTS(packet.getPTS()/90)

	 } else if (ps.currentFrame.isEmpty() || packetCounter != ps.counter) {
		 fmt.Println("Continuity = " ,packetCounter, " counter = ", ps.counter);
		 ps.counter = (packetCounter + 1) % MAX_COUNTER;
		 return;
	 }


	 if (packet.isPayloadExist()) {
		 if(ps.currentFrame.isEmpty()) {
			 ps.currentFrame.setCurrentSize(12);
		 }
		 ps.currentFrame.append(packet.getData(), packet.getPayloadOffset(),
			 packet.getPayloadLength());
		 ps.counter = (packetCounter + 1) % MAX_COUNTER;
	 }

 }