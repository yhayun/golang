//Collects Mpeg2TSPackets and builds a single frame (access unit) out of them.
//Frames are added to a PooledQueue given upon creation.
package main

import (
	"fmt"
	"time"
)

const MAX_COUNTER = 16

type Mpeg2TSParser struct {
	counter int
    output frameQueue
	currentFrame * Frame
	iFrameFound bool
	lastFrameTime long
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
 }

 func (ps Mpeg2TSParser) GetPreviousCounter(counter int) int{
	if counter == 0 {
		return 15
	}
	return counter - 1
 }

 func (ps Mpeg2TSParser) Close() {
	ps.endFlag = true;
 }


/**
* Deletes the content of the current frame and recycles it.
*/
 func (ps Mpeg2TSParser) Flush() {
	ps.output.Recylce(ps.currentFrame)
	ps.currentFrame = nil
 }


 func (ps Mpeg2TSParser) Write(packet Mpeg2TSPacket) {
 	var packetCounter int =  packet.GetContinuityCounter()

	 if !packet.IsPayloadExist() {
		 fmt.Println("Payload doens't exist")
		 return
	 }

	 if packetCounter == ps.GetPreviousCounter(ps.counter) {
		 fmt.Println("Duplicate packet found")
		 return
	 }

	 if ps.currentFrame == nil {
 		ps.currentFrame = ps.output.NewElement()
 		ps.currentFrame.Clear()
	 }

	 if ps.currentFrame == nil {
		 fmt.Println("currentFrame is null")
		 time.Sleep(5 * time.Millisecond)
		 return
	 }

	 if packet.IsStartOfPES() {
		 if !ps.currentFrame.IsEmpty() {
			 if packetCounter == ps.counter {
			 	for !ps.output.Offer(ps.currentFrame) {
			 		if (ps.endFlag){
			 			break
					}
					time.Sleep(5 * time.Millisecond)
				}
				 ps.currentFrame = ps.output.NewElement()
				 ps.currentFrame.Clear();

				 if ps.currentFrame == nil {
					 fmt.Println("currentFrame is null")
					 return
				 }
			 } else {
				 fmt.Println("currentFrame.clear()");
				 ps.currentFrame.Clear();
			 }
		 }
		 ps.currentFrame.SetPTS(packet.GetPTS()/90)

	 } else if (ps.currentFrame.IsEmpty() || packetCounter != ps.counter) {
		 fmt.Println("Continuity = " ,packetCounter, " counter = ", ps.counter);
		 ps.counter = (packetCounter + 1) % MAX_COUNTER;
		 return;
	 }


	 if (packet.IsPayloadExist()) {
		 if(ps.currentFrame.IsEmpty()) {
			 ps.currentFrame.SetCurrentSize(12);
		 }
		 ps.currentFrame.Append(packet.GetData(), packet.GetPayloadOffset(),
			 packet.GetPayloadLength());
		 ps.counter = (packetCounter + 1) % MAX_COUNTER;
	 }

 }