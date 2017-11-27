package main

import (
	"fmt"
	"time"
)

func CheckIfIFrame(buffer []byte, offset int, length int) bool {
	var startOffset int = offset

	for i := startOffset; i < (length -7); i++ {
		var j int = offset + i
		if buffer[j] == 0 {
			if buffer[j+1] == 0 {
				if buffer[j+2] == 1 {
					var val short= (short) (0x001f & buffer[j + 3])
					if val == 7 {
						return true // Start GOP was detected
					}
				}
			}
		}
	}
	return false
}


func consumer(videoFrames frameQueue) {
	h264Buffer := make([]byte, 60*1024*1024)
	var length int = 0
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
		var frame *Frame  = videoFrames.Poll();
		if !iframeDetected {
			if CheckIfIFrame(frame.GetData(),0, frame.Size()) {
				iframeDetected = true
			} else {
				videoFrames.Recylce(frame)
				continue
			}
		}
		if CheckIfIFrame(frame.GetData(),0, frame.Size()) {
			numIframes++
		}
		if numIframes >= 2 {
			videoFrames.Recylce(frame)
			continue
		}

		ArrayCopy(frame.GetData(),0, h264Buffer,length,frame.Size())
		length += frame.Size()
		videoFrames.Recylce(frame)
	}

	///todo - this is testMP4 rest of function. for now just print what we got.
	fmt.Println(h264Buffer)
}




