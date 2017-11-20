package main

type PMTFrame struct {
	data              []byte //= new byte[1280];
	size              int    // = 0;
	expectedSize      int    //= -1;
	programPID        int    //= -1;
	continuityCounter int    // = -1;
}

//CONSTRUCTOR:
func NewPMTFrame() *PMTFrame {
	return &PMTFrame{
		size: 0,
		data: make([]byte,1280),
		expectedSize: -1,
		programPID: -1,
		continuityCounter: -1,
	}
}

func (pmt PMTFrame) getStreamType(streamTypeByte byte) int {
	switch streamTypeByte {
	case 0x02:
		return TS_STREAM_TYPE_MPEG2V;
	case 0x1B:
		return TS_STREAM_TYPE_H264;
	case 0x03:
		return TS_STREAM_TYPE_MP3;
	case 0xF:
		return TS_STREAM_TYPE_AAC;
	default:
		return TS_STREAM_TYPE_UNKNOWN;
	}
}

func arrayCopy (src []byte, srcPos int, dst []byte, dstPos int , length int) {
	for i := 0; i < length; i++ {
		dst[dstPos + i] = src[srcPos + i]
	}
}

func (pmt PMTFrame) setProgramPID(programID int) {
	pmt.programPID = programID
}

func (pmt PMTFrame) getProgramPID() int {
	return pmt.programPID
}

func (pmt PMTFrame) setExpectedSize(expectedSize int) {
	pmt.expectedSize = expectedSize
}

func (pmt PMTFrame) getExpectedSize() int {
	return pmt.expectedSize
}


func (pmt PMTFrame) addPacket(p Mpeg2TSPacket) bool {
	if p.getType() != Mpeg2TSPacketType_PMT {
		return false
	}

	if pmt.continuityCounter == -1 && !p.isStart() {
		return false
	}

	if p.isStart() {
		pmt.expectedSize = p.getPMTLength()
	}

	if pmt.continuityCounter != -1 && p.getContinuityCounter() != (pmt.continuityCounter + 1)% MAX_CONTINUITY_COUNTER {
		pmt.size = 0
		pmt.continuityCounter = -1
		return false
	} else {
		arrayCopy(p.getData(), p.getDataOffset(), pmt.data, pmt.size, p.getDataLength())
		pmt.size += p.getDataLength()
		pmt.continuityCounter = p.getContinuityCounter()
		if pmt.size >= pmt.expectedSize {
			return true
		}
	}

	return false
}

func (pmt PMTFrame) getPcrPID() int {
	var pcrPID int = ((int)((pmt.data[9] & 0x1F) << 8) & 0x0000ffff) + (int)(pmt.data[10] & 0x00ff)
	return pcrPID
}

func (pmt PMTFrame) getStreamInfos() map[int]StreamInfo { //todo - make streaminfo class
	var m  = make(map[int]StreamInfo)
	var offset  = 2
	var programMapSectionLength int = ((int)(pmt.data[offset]&0x0f) <<8) + (int)(pmt.data[offset+1])

}

//public Map<Integer, StreamInfo> getStreamInfos(){
//HashMap<Integer, StreamInfo> result = new HashMap<Integer, StreamInfo>();
//
//int offset = 2;
//
//int programMapSectionLength= ((data[offset]&0x0f) <<8) + data[offset+1];
//
////Skip everything till start of streams array
//offset+=9;
//int programInfoLength = ((data[offset]&0xf) <<8) + data[offset+1];
//offset+=programInfoLength+2;
//
//while(offset < programMapSectionLength-1){
//StreamInfo streamInfo = new StreamInfo();
//
//streamInfo.streamType = getStreamType(data[offset]);
//offset +=1;
//streamInfo.streamPID = (((data[offset] & 0x1F) << 8)& 0x0000ffff) + (data[offset+1]& 0x00ff);
//offset += 2;
//int esInfoLength = ((data[offset]&0xf)<<8)+data[offset+1];
//offset += esInfoLength +2;
//result.put(streamInfo.streamPID, streamInfo);
//}
//
//return result;
//
//}