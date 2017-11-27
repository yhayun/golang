package main

import (
	"fmt"
	"time"
)

func consumer(videoFrames frameQueue) {
	var iframeDetected bool = false
	var numIframes int = 0
	for {
		if videoFrames.IsEmpty() {
			time.Sleep(100 * time.Millisecond)
			if videoFrames.IsEmpty() {
				time.Sleep(8000 * time.Millisecond)
				if videoFrames.IsEmpty() {
					break
				}
			}
		}
	}

	var _frame frame  = videoFrames.poll();
	if !iframeDetected {
		if
	}
}

boolean iframeDetected = false;
int numIframes = 0;
while(true) {
	if(videoFrames.isEmpty()) {
		Sleeper.sleep(100);
		if(videoFrames.isEmpty()) {
			Sleeper.sleep(8000);
			if(videoFrames.isEmpty()) {
				break;
			}
		}
	}
	Frame frame = videoFrames.poll();
	if(!iframeDetected) {
		if(checkIfIFrame(frame.getData(), 0, frame.size())) {
			iframeDetected = true;
		}
		else {
			videoFrames.recycle(frame);
		continue;
		}
	}
	if(checkIfIFrame(frame.getData(), 0, frame.size())) {
		numIframes++;
	}
	if(numIframes >= 2) {
		videoFrames.recycle(frame);
	continue;
	}

	//System.out.println(frame.size());
	System.arraycopy(frame.getData(), 0, h264Buffer, length, frame.size());
	length += frame.size();
	videoFrames.recycle(frame);
}