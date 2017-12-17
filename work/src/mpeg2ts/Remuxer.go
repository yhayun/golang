package main

// todo - https://github.com/video-dev/hls.js/blob/ce19344d1bfb025e4d55c5b1521f0c8f02c06b5b/src/demux/tsdemuxer.js

type VideoTrack struct {
	id uint
	pid uint
	inputTimeScale uint
	sequenceNumber uint
	samples []byte
	len uint
	dropped uint
	isAAC bool
	duration bool

	//For the first frame also:
	width uint
	height uint
	sps uint//(GOP header)
	pps uint //(I frame header)

}

func NewVideoTrack() *VideoTrack {
	return &VideoTrack{
		id: 0,
		pid : -1,
		inputTimeScale : 90000,
		sequenceNumber: 0,
		//samples : [],
		len : 0,
		dropped: 0, // type === 'video' ? 0 : undefined,
		isAAC: false, //type === 'audio' ? true : undefined,
		duration: false, //type === 'audio' ? duration : undefined
		
		//First Frame only:
		width: -1,
		height: -1,
		sps: -1,
		pps: -1,
	}
}


type MP4Remuxer struct {

}

func NewMP4Remuxer() *MP4Remuxer {
	return &MP4Remuxer{

		}
}

