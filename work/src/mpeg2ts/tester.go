package main

import (
	"fmt"
	"time"
	"io/ioutil"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteFile(buffer []byte, file string){
	err := ioutil.WriteFile(file, buffer, 0644)
	check(err)
}

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
	fmt.Println("entered consmer")
	h264Buffer := make([]byte, 60*1024*1024)
	var length int = 0
	var iframeDetected bool = false
	var numIframes int = 0
	for {
		//fmt.Println("entered consmer loop")
		if videoFrames.IsEmpty() {
			fmt.Println("consume sleep 100ms")
			time.Sleep(100 * time.Millisecond)
			if videoFrames.IsEmpty() {
				fmt.Println("consume sleep 8000ms")
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
		fmt.Println(h264Buffer)
	}
	fmt.Println("left consumer loop.")
	///todo - this is testMP4 rest of function. for now just print what we got.
	//fmt.Println(h264Buffer)
	first_100 := h264Buffer[:100]
	fmt.Println(first_100)
	Done <- true
	//WriteFile(h264Buffer,"tempfile.txt");
}




