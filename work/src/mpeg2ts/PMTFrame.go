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

func (pmt PMTFrame) GetStreamType(streamTypeByte byte) int {
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

func ArrayCopy (src []byte, srcPos int, dst []byte, dstPos int , length int) {
	for i := 0; i < length; i++ {
		dst[dstPos + i] = src[srcPos + i]
	}
}

func (pmt PMTFrame) SetProgramPID(programID int) {
	pmt.programPID = programID
}

func (pmt PMTFrame) GetProgramPID() int {
	return pmt.programPID
}

func (pmt PMTFrame) SetExpectedSize(expectedSize int) {
	pmt.expectedSize = expectedSize
}

func (pmt PMTFrame) GetExpectedSize() int {
	return pmt.expectedSize
}


func (pmt PMTFrame) AddPacket(p Mpeg2TSPacket) bool {
	if p.getType() != Mpeg2TSPacketType_PMT {
		return false
	}

	if pmt.continuityCounter == -1 && !p.IsStart() {
		return false
	}

	if p.IsStart() {
		pmt.expectedSize = p.GetPMTLength()
	}

	if pmt.continuityCounter != -1 && p.GetContinuityCounter() != (pmt.continuityCounter + 1)% MAX_CONTINUITY_COUNTER {
		pmt.size = 0
		pmt.continuityCounter = -1
		return false
	} else {
		ArrayCopy(p.getData(), p.GetDataOffset(), pmt.data, pmt.size, p.GetDataLength())
		pmt.size += p.GetDataLength()
		pmt.continuityCounter = p.GetContinuityCounter()
		if pmt.size >= pmt.expectedSize {
			return true
		}
	}

	return false
}

func (pmt PMTFrame) GetPcrPID() int {
	var pcrPID int = ((int)((pmt.data[9] & 0x1F) << 8) & 0x0000ffff) + (int)(pmt.data[10] & 0x00ff)
	return pcrPID
}

func (pmt PMTFrame) GetStreamInfos() map[int]StreamInfo {
	var result = make(map[int]StreamInfo)
	var offset  = 2
	var programMapSectionLength int = ((int)(pmt.data[offset]&0x0f) <<8) + (int)(pmt.data[offset+1])

	//Skip everything till start of streams array
	offset += 9
	var programInfoLength int = (int) (((pmt.data[offset]&0xf) <<8) + pmt.data[offset+1]);
	offset+= programInfoLength + 2

	for offset < programMapSectionLength -1 {
		var streamInfo StreamInfo
		streamInfo.streamType = pmt.GetStreamType(pmt.data[offset])
		offset +=1
		streamInfo.streamPID = (int)((((pmt.data[offset] & 0x1F) << 8)& 0x0000ffff) + (pmt.data[offset+1]& 0x00ff))
		offset += 2
		var esInfoLength int = (int)(((pmt.data[offset]&0xf)<<8)+ pmt.data[offset+1])
		offset += esInfoLength +2;
		result[streamInfo.streamPID] = streamInfo;
	}

	return result;
}
