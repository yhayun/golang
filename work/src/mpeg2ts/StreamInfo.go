package main

import "fmt"

const TS_STREAM_TYPE_MPEG2V = 0
const TS_STREAM_TYPE_H264 = 1
const TS_STREAM_TYPE_AAC = 2
const TS_STREAM_TYPE_MP3 = 3
const TS_STREAM_TYPE_AC3 = 4
const TS_STREAM_TYPE_G711 = 5
const TS_STREAM_TYPE_UNKNOWN = 6


type StreamInfo struct {
	streamType int
	streamPID int
}

//CONSTRUCTOR:
func NewStreamInfo(capacity int) *StreamInfo {
	return &StreamInfo{}
}

func ( s StreamInfo) ToString() string{
	 res := fmt.Sprintf( "Type: %d , PID: %d",s.streamType, s.streamPID)

	 return res
}
