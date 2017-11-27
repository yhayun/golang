package main


type long int64
type short int16
const TRANSPORT_PACKET_SIZE int = 188
const MAX_CONTINUITY_COUNTER int = 16
const Mpeg2TSPacketType_PAT int  = 0
const Mpeg2TSPacketType_PMT int  = 1
const Mpeg2TSPacketType_PES int  = 2
const VideoFrameType_I_START = 0
const VideoFrameType_P_START = 1
const VideoFrameType_OTHER = 2
type Mpeg2TSPacket struct {
	data []byte
	offset int
	valid bool
	Mpeg2TSPacketType int // Mpeg2TSPacketType PAT = 0, PMT =1, PES = 2 for now.. TODO define the enum correctl
	VideoFrameType int // I_START = 0, P_START = 1, OTHER = 2
	pid int
	continuityCounter int
	dataOffset int
	dataLength int
	adaptationExists bool
	pcrExists bool
	pcr long
	//// in case of PES
	startOfPES bool
	pts long
	payloadExist bool
	payloadOffset int
	payloadLength int
	//// in case of PAT
	programPID int
	//// in case of PMT
	pmtLength int
	start bool
	receiveTime long
}

func (tsP* Mpeg2TSPacket) IsPayloadExist() bool{
	return tsP.payloadExist
}





func (tsP* Mpeg2TSPacket) IsStartOfPES() bool{
	return tsP.startOfPES;
}

func ValidMpegTsPacket(data []byte, offset int) bool{
	return data[offset] == 0x47
}

func (tsP* Mpeg2TSPacket) FromBytes(data []byte, offset int, programPID int) {
	tsP.Reset()

	if !ValidMpegTsPacket(data, offset) {
		tsP.valid = false;
		return;
	}

	tsP.data = data
	tsP.offset = offset

	tsP.valid = true

	tsP.pid = GetPID(data, offset);

	tsP.dataOffset = TsDataOffset(data, offset);
	tsP.dataLength = TRANSPORT_PACKET_SIZE + offset - tsP.dataOffset;
	tsP.continuityCounter = GetContinuityCounter(data, offset);
	tsP.adaptationExists = IsAdaptationFieldExist(data, offset);
	//TODO PCR ignored

	if tsP.pid == 0 {
		tsP.Mpeg2TSPacketType = 0
		tsP.programPID = GetProgramPID(data, offset);
	} else if (tsP.pid == programPID) {
		tsP.Mpeg2TSPacketType = 1
		var start bool = IsStart(data, offset);
		if (start) {
			tsP.pmtLength = GetPMTLength(data, offset);
		}
	} else {
		tsP.Mpeg2TSPacketType = 2
		tsP.startOfPES = IsStartOfPES(data, offset);
		tsP.payloadExist = PayloadExists(data, offset);
		if (tsP.payloadExist) {
			tsP.payloadOffset = GetPayloadOffset(data, offset);
			tsP.payloadLength = TRANSPORT_PACKET_SIZE + offset - tsP.payloadOffset;
		}
		if (tsP.startOfPES) {
			tsP.pts = GetPTS(data, offset);
		}
	}
}


func (tsP* Mpeg2TSPacket) Reset() {
	tsP.valid = false
	tsP.pid = 0;
	tsP.Mpeg2TSPacketType = -1;
	tsP.continuityCounter = 0;
	tsP.dataOffset = 0;
	tsP.dataLength = 0;
	tsP.adaptationExists = false;
	tsP.pcrExists = false;
	tsP.pcr = 0;
	tsP.startOfPES = false;
	tsP.pts = 0;
	tsP.payloadExist = false;
	tsP.payloadOffset = 0;
	tsP.payloadLength = 0;
	tsP.programPID = 0;
	tsP.pmtLength = 0;
	tsP.start = false;
	tsP.receiveTime = 0;
}

func IsAdaptationFieldExist(buffer [] byte, offset int) bool{
	return ((buffer[3 + offset] & 0x20) != 0);
}

func AdaptationFieldLength(buffer []byte , offset int) int {
	var length int = (int)(buffer[4 + offset] & 0x00ff) ;
	return length;
}

func TsDataOffset(buffer []byte , offset int) int {
	var internalOffset int = 4 + offset
	if IsAdaptationFieldExist(buffer, offset) {
		internalOffset += 1 + AdaptationFieldLength(buffer, offset)
	}
	return internalOffset;
}

func GetPID(buffer []byte , offset int) int {
	var pid int = (int)(((buffer[1 + offset] & 0x1f) << 8) & 0x0000ffff)+ (int)(buffer[2 + offset] & 0x00ff)
	return pid;
}

func GetContinuityCounter(buffer []byte , offset int) int {
	var counter int = int(buffer[3 + offset] & 0x0f)
	return counter
}

func (tsP Mpeg2TSPacket) GetContinuityCounter() int {
	return tsP.continuityCounter
}


func GetProgramPID(buffer []byte , offset int) int {
	var dataOffset int = TsDataOffset(buffer, offset);
	return (int)((buffer[dataOffset + 11] & 0x1F) << 8) + (int)(0x0ff & buffer[dataOffset + 12]);
}

func IsStart(tsPacket []byte, offset int) bool {
	return (tsPacket[1 + offset] & 0x40) != 0;
}

func (tsP Mpeg2TSPacket) IsStart() bool{
	return tsP.start
}

func (tsP* Mpeg2TSPacket) GetType() int {
	return tsP.Mpeg2TSPacketType;
}

func (tsP* Mpeg2TSPacket) GetData() []byte {
	return tsP.data;
}

func GetPMTLength(buffer []byte, offset int) int{
	var tsDataOffset int = TsDataOffset(buffer, offset) + 2
	var pmtSectionLength int = (int)((buffer[tsDataOffset] & 0xf) << 8) + (int)(buffer[1+tsDataOffset])
	return pmtSectionLength + 3;
}

func (tsP Mpeg2TSPacket) GetPMTLength() int{
	return tsP.pmtLength
}

func (tsP Mpeg2TSPacket) GetDataOffset() int{
	return tsP.dataOffset
}

func (tsP Mpeg2TSPacket) GetDataLength() int{
	return tsP.dataLength
}

func IsStartOfPES(buffer []byte, offset int) bool {
	if (buffer[1+offset] & 0x40) == 0 {
		return false
	}
	var tsDataOffset int = TsDataOffset(buffer, offset)
	var b bool = buffer[tsDataOffset] == 0 && buffer[tsDataOffset + 1] == 0 && buffer[tsDataOffset + 2] == 0x01
	return b
}
func PayloadExists(buffer []byte, offset int) bool{
	return (buffer[3 + offset] & 0x10) != 0
}

func GetPayloadOffset(buffer []byte, offset int) int{
	var internalOffset int = TsDataOffset(buffer, offset)

	if (IsStartOfPES(buffer, offset)) {
		var pes_header_data_length byte = buffer[internalOffset + 8]
		internalOffset += 9 + (int)(pes_header_data_length)
	}
	return internalOffset;
}

func (tsP Mpeg2TSPacket) GetPayloadOffset() int {
	return tsP.payloadOffset
}

func (tsP Mpeg2TSPacket) GetPayloadLength() int {
	return tsP.payloadLength
}

func GetPTS(buffer []byte, offset int) long {

	var dataOffset int = TsDataOffset(buffer, offset);

	if ((buffer[7 + dataOffset] & 0x80) == 0) {
	return -1;
	}

	var ptsOffset int= 9 + dataOffset;

	var pts long;

	pts = ((long) ((buffer[ptsOffset] & 0x0e) >> 1)) << 30;
	pts += ((long) (buffer[1 + ptsOffset] & 0xff) << 22);
	pts += ((long) ((buffer[2 + ptsOffset] & 0xfe) >> 1)) << 15;
	pts += ((long) (buffer[3 + ptsOffset] & 0xff) << 7);
	pts += (long)((buffer[4 + ptsOffset] & 0xfe) >> 1);
	return pts;
}

func (tsP Mpeg2TSPacket) GetPTS() long {
	return tsP.pts
}

func (tsP* Mpeg2TSPacket) GetH264type() int16 {

	var buffer []byte = tsP.data;

	var startOffset int = TsDataOffset(buffer, tsP.offset) - tsP.offset;

	// If start of PES, look for Sequence Parameter Set (SPS)
	// If it was found then it is I frame / GOP start otherwise P frame
	if (tsP.IsStartOfPES()) {
		for i := startOffset; i < (TRANSPORT_PACKET_SIZE - 5); i++ {
			var j int= tsP.offset + i;
			if (buffer[j] == 0) {
				if (buffer[j + 1] == 0) {
					if (buffer[j + 2] == 1) {
						var val short= (short) (0x001f & buffer[j + 3]);
						if (val == 7) {
							return 0; // Start GOP was detected
						}
					}
				}
			}
		}
		return 1; // P frame was detected
	}
	return -1;
}

//todo increaseContinuityCounter()
func (tsP* Mpeg2TSPacket) GetH264VideoFrameType() int{
	var _type int16 = tsP.GetH264type();
	if _type == 0 {
		return VideoFrameType_I_START;
	} else if (_type == 1) || (_type == 6) {
		return VideoFrameType_P_START;
	}
return VideoFrameType_OTHER;

}


