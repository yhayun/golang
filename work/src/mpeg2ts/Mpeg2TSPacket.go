package main

import (
	"container/list"
)

type long uint64
const TRANSPORT_PACKET_SIZE int = 188
const MAX_CONTINUITY_COUNTER int = 16

type Mpeg2TSPacket struct {
	data []byte
	offset int
	valid bool
	Mpeg2TSPacketType int // Mpeg2TSPacketType PAT = 0, PMT =1, PES = 2 for now.. TODO define the enum correctl
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

func validMpegTsPacket(data []byte, offset int) bool{
	return data[offset] == 0x47
}

func (tsP* Mpeg2TSPacket) fromBytes(data []byte, offset int, programPID int) {
	tsP.reset()

	if !validMpegTsPacket(data, offset) {
		tsP.valid = false;
		return;
	}

	tsP.data = data
	tsP.offset = offset

	tsP.valid = true

	tsP.pid = getPID(data, offset);

	tsP.dataOffset = tsDataOffset(data, offset);
	tsP.dataLength = TRANSPORT_PACKET_SIZE + offset - tsP.dataOffset;
	tsP.continuityCounter = getContinuityCounter(data, offset);
	tsP.adaptationExists = isAdaptationFieldExist(data, offset);
	//TODO PCR ignored

	if tsP.pid == 0 {
		tsP.Mpeg2TSPacketType = 0
		tsP.programPID = getProgramPID(data, offset);
	} else if (tsP.pid == programPID) {
		tsP.Mpeg2TSPacketType = 1
		var start bool = isStart(data, offset);
		if (start) {
			tsP.pmtLength = getPMTLength(data, offset);
		}
	} else {
		tsP.Mpeg2TSPacketType = 2
		tsP.startOfPES = isStartOfPES(data, offset);
		tsP.payloadExist = payloadExists(data, offset);
		if (tsP.payloadExist) {
			tsP.payloadOffset = getPayloadOffset(data, offset);
			tsP.payloadLength = TRANSPORT_PACKET_SIZE + offset - tsP.payloadOffset;
		}
		if (tsP.startOfPES) {
			tsP.pts = getPTS(data, offset);
		}
	}
}


func (tsP* Mpeg2TSPacket) reset() {
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

func isAdaptationFieldExist(buffer [] byte, offset int) bool{
	return ((buffer[3 + offset] & 0x20) != 0);
}

func adaptationFieldLength(buffer []byte , offset int) int {
	var length int = (int)(buffer[4 + offset] & 0x00ff) ;
	return length;
}

func tsDataOffset(buffer []byte , offset int) int {
	var internalOffset int = 4 + offset
	if isAdaptationFieldExist(buffer, offset) {
		internalOffset += 1 + adaptationFieldLength(buffer, offset)
	}
	return internalOffset;
}

func getPID(buffer []byte , offset int) int {
	var pid int = (int)(((buffer[1 + offset] & 0x1f) << 8) & 0x0000ffff)+ (int)(buffer[2 + offset] & 0x00ff)
	return pid;
}

func getContinuityCounter(buffer []byte , offset int) int {
	var counter int = int(buffer[3 + offset] & 0x0f)
	return counter
}


func getProgramPID(buffer []byte , offset int) int {
	var dataOffset int = tsDataOffset(buffer, offset);
	return (int)((buffer[dataOffset + 11] & 0x1F) << 8) + (int)(0x0ff & buffer[dataOffset + 12]);
}

func isStart(tsPacket []byte, offset int) bool {
	return (tsPacket[1 + offset] & 0x40) != 0;
}

func getPMTLength(buffer []byte, offset int) int{
	var tsDataOffset int = tsDataOffset(buffer, offset) + 2
	var pmtSectionLength int = (int)((buffer[tsDataOffset] & 0xf) << 8) + (int)(buffer[1+tsDataOffset])
	return pmtSectionLength + 3;
}
func isStartOfPES(buffer []byte, offset int) bool {
	if (buffer[1+offset] & 0x40) == 0 {
		return false
	}
	var tsDataOffset int = tsDataOffset(buffer, offset)
	var b bool = buffer[tsDataOffset] == 0 && buffer[tsDataOffset + 1] == 0 && buffer[tsDataOffset + 2] == 0x01
	return b
}
func payloadExists(buffer []byte, offset int) bool{
	return (buffer[3 + offset] & 0x10) != 0
}

func getPayloadOffset(buffer []byte, offset int) int{
	var internalOffset int = tsDataOffset(buffer, offset)

	if (isStartOfPES(buffer, offset)) {
		var pes_header_data_length byte = buffer[internalOffset + 8]
		internalOffset += 9 + (int)(pes_header_data_length)
	}
	return internalOffset;
}


func getPTS(buffer []byte, offset int) long {

var dataOffset int = tsDataOffset(buffer, offset);

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